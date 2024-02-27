package infrastructure

import (
	"strings"

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
	"github.com/jramsgz/articpad/pkg/i18n"
	"github.com/jramsgz/articpad/pkg/mail"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// startFiberServer starts the Fiber server.
func startFiberServer(logger zerolog.Logger, db *gorm.DB, mailClient *mail.Mailer, i18n *i18n.I18n) *fiber.App {
	var trustedProxies []string
	if config.GetString("TRUSTED_PROXIES") != "" {
		trustedProxies = strings.Split(config.GetString("TRUSTED_PROXIES"), ",")
	}
	var enableProxy bool = len(trustedProxies) > 0

	var isProduction bool = config.GetString("DEBUG") == "false"
	app := fiber.New(fiber.Config{
		Prefork:                 isProduction,
		ServerHeader:            "ArticPad Server " + config.Version,
		AppName:                 "ArticPad",
		DisableStartupMessage:   isProduction,
		EnableTrustedProxyCheck: enableProxy,
		TrustedProxies:          trustedProxies,
	})

	app.Use(logging.Logger(logger, func(c *fiber.Ctx) bool {
		return c.Path() == "/health" // skip logging for health check
	}))
	app.Use(cors.New(cors.Config{
		MaxAge:       1800,
		AllowOrigins: config.GetString("APP_URL"),
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(etag.New())
	app.Use(limiter.New(limiter.Config{
		// TODO: Make this configurable and enable it for auth routes. Also fix that every fork has its own counter by using redis or something.
		Max: 100,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "You have exceeded the maximum number of requests. Please try again later.")
		},
	}))

	// Create repositories.
	userRepository := user.NewUserRepository(db)

	// Create all of our services.
	userService := user.NewUserService(userRepository)

	// Prepare our endpoints for the API.
	api := app.Group("/api")
	apiv1 := api.Group("/v1")

	misc.NewMiscHandler(apiv1)
	health.NewHealthHandler(app.Group("/health"))
	auth.NewAuthHandler(apiv1.Group("/auth"), userService, mailClient, i18n)
	//user.NewUserHandler(apiv1.Group("/users"), userService)

	// Prepare an endpoint for 'Not Found'.
	api.All("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Not Found",
		})
	})

	// Serve Single Page application on "/"
	// assume static file at static folder
	app.Static("/", config.GetString("STATIC_DIR"), fiber.Static{
		Compress: true,
		MaxAge:   3600,
	})

	// Panic test route, this brings up an error
	// TODO: Remove this route in production
	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + config.GetString("STATIC_DIR") + "/index.html")
	})

	return app
}
