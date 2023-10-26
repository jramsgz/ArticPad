package logging

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jramsgz/articpad/config"
	"github.com/jramsgz/articpad/pkg/apierror"
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
	ErrorCode  string
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

	if lf.ErrorCode != "" {
		e.Str("error_code", lf.ErrorCode)
	}
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
			var isProduction bool = config.GetString("DEBUG") == "false"

			// Status code defaults to 500
			status := fiber.StatusInternalServerError
			// Message defaults to "Internal Server Error"
			message := "Internal Server Error"
			// Code defaults to "unknown_error"
			code := "unknown_error"

			// Check if c.Next() returned an error
			if err != nil {

				// Check if it's a fiber.Error
				if e, ok := err.(*fiber.Error); ok {
					// Override status code if fiber.Error type
					status = e.Code
					// If the error is not a server error, send the error message to the client
					if !isProduction || (status >= 400 && status < 500) {
						message = e.Message
					}
				} else if e, ok := err.(*apierror.Error); ok {
					// Override status code and error code if apierror.Error type
					status = e.Status
					code = e.Code
					// If the error is not a server error, send the error message to the client
					if !isProduction || (status >= 400 && status < 500) || e.Show {
						message = e.Message
					}
				}

				fields.ErrorCode = code
				fields.Error = err

				// Send custom error page
				c.Status(status).JSON(fiber.Map{
					"success":    false,
					"error_code": code,
					"error":      message,
					"requestId":  rid,
				})
			}

			// Check if panic() was called
			rvr := recover()
			if rvr != nil {
				err, ok := rvr.(error)
				if !ok {
					err = fmt.Errorf("%v", rvr)
				}

				fields.ErrorCode = code
				fields.Error = err
				fields.Stack = debug.Stack()

				if !isProduction {
					message = err.Error()
				}

				c.Status(status).JSON(fiber.Map{
					"success":    false,
					"error_code": code,
					"error":      message,
					"requestId":  rid,
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
