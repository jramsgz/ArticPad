package infrastructure

import (
	"strings"
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

// startFiberServer starts the Fiber server.
func (a *App) startFiberServer() *fiber.App {
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

	app.Use(logging.Logger(a.logger, func(c *fiber.Ctx) bool {
		return c.Path() == "/health"
	}))
	app.Use(cors.New(cors.Config{
		MaxAge:       1800,
		AllowOrigins: config.GetString("APP_URL"),
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(etag.New())
	if config.GetString("RATE_LIMIT_AUTH") == "true" {
		app.Use(limiter.New(limiter.Config{
			Max:        40,
			Expiration: 1 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return fiber.NewError(fiber.StatusTooManyRequests, "You have exceeded the maximum number of requests. Please try again later.")
			},
			Storage: func() fiber.Storage {
				if a.redis != nil {
					return a.redis
				}
				return limiter.ConfigDefault.Storage
			}(),
			Next: func(c *fiber.Ctx) bool {
				return !strings.HasPrefix(c.Path(), "/api/v1/auth")
			},
		}))
	}

	userRepository := user.NewUserRepository(a.db)

	userService := user.NewUserService(userRepository)

	api := app.Group("/api")
	apiv1 := api.Group("/v1")

	misc.NewMiscHandler(apiv1)
	health.NewHealthHandler(app.Group("/health"))
	auth.NewAuthHandler(apiv1.Group("/auth"), userService, a.mail, a.i18n)
	//user.NewUserHandler(apiv1.Group("/users"), userService)

	api.All("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "Not Found",
		})
	})

	app.Static("/", config.GetString("STATIC_DIR"), fiber.Static{
		Compress: true,
		MaxAge:   3600,
	})

	// TODO: Remove this route in production
	app.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./" + config.GetString("STATIC_DIR") + "/index.html")
	})

	return app
}
