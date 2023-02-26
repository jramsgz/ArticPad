package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/middleware"
	"github.com/jramsgz/articpad/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	// Version of the application
	Version = "dev"
	// BuildTime of the application
	BuildTime = time.Now().Format(time.RFC3339)
	// Commit of the application
	Commit = "dev build"

	// Build information. Populated at build-time.
	// go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=2020-01-01T00:00:00Z -X main.Commit=abcdef"
)

func main() {
	var isProduction bool = config.Config("DEBUG", "false") == "false"

	app := fiber.New(fiber.Config{
		Prefork:               isProduction,
		ServerHeader:          "ArticPad Server " + Version,
		AppName:               "ArticPad",
		DisableStartupMessage: isProduction,
	})

	// Logging
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		os.MkdirAll("./logs", 0755)
	}
	logFile, err := os.OpenFile("./logs/articpad.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${locals:requestid} ${status} - ${latency} ${ip} ${method} ${path}\n",
		Output: mw,
		// Logs are disabled for requests to the static files (not prefixed with /api)
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/api"
		},
	}))
	log.SetOutput(mw)

	// Middleware registration
	middleware.RegisterMiddlewares(app)

	//database.ConnectDB()

	router.SetupRoutes(app)

	if !fiber.IsChild() {
		log.Printf("Starting ArticPad %s with isProduction: %t", Version, isProduction)
		log.Printf("BuildTime: %s | Commit: %s", BuildTime, Commit)
		log.Printf("Listening on %s", config.Config("APP_ADDR", ":3000"))
	}
	log.Fatal(app.Listen(config.Config("APP_ADDR", ":3000")))
}
