package infrastructure

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/internal/auth"
	"github.com/jramsgz/articpad/internal/health"
	"github.com/jramsgz/articpad/internal/logging"
	"github.com/jramsgz/articpad/internal/misc"
	"github.com/jramsgz/articpad/internal/user"
)

// Run ArticPad API & Static Server
func Run() {
	// Load configuration from .env file.
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	// Start logger.
	logger, _, logFile := StartLogger(&LoggerConfig{
		Level: config.GetString("LOG_LEVEL", "debug"),
		Dir:   config.GetString("LOG_DIR", "./logs"),
	})

	// Try to connect to the specified database.
	db, err := ConnectToDB(&DatabaseConfig{
		Driver:   config.GetString("DB_DRIVER", "sqlite"),
		Host:     config.GetString("DB_HOST", "localhost"),
		Username: config.GetString("DB_USERNAME", "root"),
		Password: config.GetString("DB_PASSWORD", ""),
		Port:     config.GetInt("DB_PORT", 3306),
		Database: config.GetString("DB_DATABASE", "config/articpad.db"),
	})
	if err != nil || db == nil {
		log.Fatal("Database connection error: $s", err)
	}

	// Set trusted proxies
	var trustedProxies []string
	if config.GetString("TRUSTED_PROXIES", "") != "" {
		trustedProxies = strings.Split(config.GetString("TRUSTED_PROXIES", ""), ",")
	}
	var enableProxy bool = len(trustedProxies) > 0

	// Creates a new Fiber instance.
	var isProduction bool = config.GetString("DEBUG", "false") == "false"
	app := fiber.New(fiber.Config{
		Prefork:                 isProduction,
		ServerHeader:            "ArticPad Server " + config.Version,
		AppName:                 "ArticPad",
		DisableStartupMessage:   isProduction,
		EnableTrustedProxyCheck: enableProxy,
		TrustedProxies:          trustedProxies,
	})

	// Setup graceful shutdown.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	var serverShutdown sync.WaitGroup

	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(60 * time.Second)
	}()

	if !fiber.IsChild() {
		// Auto-migrate database models
		err := db.AutoMigrate(&user.User{})
		if err != nil {
			log.Fatal("failed to automigrate models:", err.Error())
			return
		}
	}

	// Use global middlewares.
	app.Use(logging.Logger(logger, func(c *fiber.Ctx) bool {
		return c.Path() == "/health" // skip logging for health check
	}))
	app.Use(cors.New(cors.Config{
		MaxAge:       1800,
		AllowOrigins: config.GetString("APP_URL", "http://localhost:8080"),
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(etag.New())
	app.Use(limiter.New(limiter.Config{
		Max: 100,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(&fiber.Map{
				"status":  "fail",
				"message": "You have requested too many in a single time-frame! Please wait another minute!",
			})
		},
	}))

	// Create repositories.
	//cityRepository := city.NewCityRepository(db)
	userRepository := user.NewUserRepository(db)

	// Create all of our services.
	//cityService := city.NewCityService(cityRepository)
	userService := user.NewUserService(userRepository)

	// Prepare our endpoints for the API.
	api := app.Group("/api")
	apiv1 := api.Group("/v1")

	misc.NewMiscHandler(apiv1)
	health.NewHealthHandler(app.Group("/health"))
	auth.NewAuthHandler(apiv1.Group("/auth"))
	//city.NewCityHandler(apiv1.Group("/cities"), cityService)
	user.NewUserHandler(apiv1.Group("/users"), userService)

	// Prepare an endpoint for 'Not Found'.
	api.All("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Not Found",
		})
	})

	// Serve Single Page application on "/"
	// assume static file at static folder
	app.Static("/", config.GetString("STATIC_DIR", "static"), fiber.Static{
		Compress: true,
		MaxAge:   3600,
	})

	// Panic test route, this brings up an error
	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + config.GetString("STATIC_DIR", "static") + "/index.html")
	})

	if !fiber.IsChild() {
		log.Printf("Starting ArticPad %s with isProduction: %t", config.Version, isProduction)
		log.Printf("BuildTime: %s | Commit: %s", config.BuildTime, config.Commit)
		log.Printf("Listening on %s", config.GetString("APP_ADDR", ":8080"))
	}
	if err := app.Listen(config.GetString("APP_ADDR", ":8080")); err != nil {
		logger.Fatal().Err(err).Msg("Error starting server")
	}

	if !fiber.IsChild() {
		serverShutdown.Wait()

		fmt.Println("Running cleanup tasks...")
		_ = logFile.Close()
	}
}
