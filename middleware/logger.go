package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jramsgz/articpad/config"
	"github.com/rs/zerolog"
)

type logFields struct {
	ID         string
	RemoteIP   string
	Host       string
	Method     string
	Path       string
	Protocol   string
	StatusCode int
	Latency    float64
	Error      error
	Stack      []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("id", lf.ID).
		Str("remote_ip", lf.RemoteIP).
		Str("host", lf.Host).
		Str("method", lf.Method).
		Str("path", lf.Path).
		Str("protocol", lf.Protocol).
		Int("status_code", lf.StatusCode).
		Float64("latency", lf.Latency).
		Str("tag", "request")

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}

// Logger requestid + logger + recover for request traceability
func Logger(log zerolog.Logger, filter func(*fiber.Ctx) bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if filter != nil && filter(c) {
			c.Next()
			return nil
		}

		start := time.Now()

		// Only generate request id for non-static files
		rid := "static"
		// Check if request path starts with /api
		if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
			rid = uuid.New().String()
			c.Set(fiber.HeaderXRequestID, rid)
		}

		fields := &logFields{
			ID:       rid,
			RemoteIP: c.IP(),
			Method:   c.Method(),
			Host:     c.Hostname(),
			Path:     c.Path(),
			Protocol: c.Protocol(),
		}

		var err error
		defer func() {
			// There are two possible ways to get an error here:
			// 1. c.Next() returns an error
			// 2. panic() is called

			// We only send the error to the client if DEBUG is not set
			var isProduction bool = config.Config("DEBUG", "false") == "false"

			// Check if c.Next() returned an error
			if err != nil {
				// Status code defaults to 500
				code := fiber.StatusInternalServerError
				// Message defaults to "Internal Server Error"
				message := "Internal Server Error"

				// Check if it's a fiber.Error
				if e, ok := err.(*fiber.Error); ok {
					// Override status code if fiber.Error type
					code = e.Code
					if !isProduction {
						message = e.Message
					}
				}

				fields.Error = err
				fields.Stack = debug.Stack()

				// Send custom error page
				c.Status(code).JSON(fiber.Map{
					"success":   false,
					"error":     message,
					"requestId": rid,
				})
			}

			// Check if panic() was called
			rvr := recover()
			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				fields.Error = err
				fields.Stack = debug.Stack()

				message := "Internal Server Error"
				if !isProduction {
					message = err.Error()
				}

				c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"success":   false,
					"error":     message,
					"requestId": rid,
				})
			}

			fields.StatusCode = c.Response().StatusCode()
			fields.Latency = time.Since(start).Seconds()

			switch {
			case rvr != nil:
				log.Error().EmbedObject(fields).Msg("panic recover")
			case fields.StatusCode >= 500:
				log.Error().EmbedObject(fields).Msg("server error")
			case fields.StatusCode >= 400:
				log.Warn().EmbedObject(fields).Msg("client error")
			case fields.StatusCode >= 300:
				log.Warn().EmbedObject(fields).Msg("redirect")
			case fields.StatusCode >= 200:
				log.Info().EmbedObject(fields).Msg("success")
			case fields.StatusCode >= 100:
				log.Info().EmbedObject(fields).Msg("informative")
			default:
				log.Warn().EmbedObject(fields).Msg("unknown status")
			}
		}()

		err = c.Next()
		return nil
	}
}
