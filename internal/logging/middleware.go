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

		rid := generateRequestID(c)
		fields := createLogFields(c, rid)

		var err error
		defer func() {
			handleError(log, c, fields, start, err)
			rvr := recover()
			handlePanic(log, c, fields, start, rvr)
			logEvent(log, fields, rvr)
		}()

		err = c.Next()
		return nil
	}
}

func generateRequestID(c *fiber.Ctx) string {
	if len(c.Path()) >= 4 && c.Path()[:4] == "/api" {
		return uuid.New().String()
	}
	return "static"
}

func createLogFields(c *fiber.Ctx, rid string) *logFields {
	return &logFields{
		ID:       rid,
		RemoteIP: c.IP(),
		Method:   c.Method(),
		Host:     c.Hostname(),
		Path:     c.Path(),
		Protocol: c.Protocol(),
	}
}

func handleError(log zerolog.Logger, c *fiber.Ctx, fields *logFields, start time.Time, err error) {
	if err != nil {
		status, message, code := getErrorMessage(err)
		fields.StatusCode = status
		fields.ErrorCode = code
		fields.Error = err

		c.Status(status).JSON(fiber.Map{
			"success":    false,
			"error_code": code,
			"error":      message,
			"requestId":  fields.ID,
		})
	}

	fields.StatusCode = c.Response().StatusCode()
	fields.Latency = time.Since(start).Seconds()
}

func handlePanic(log zerolog.Logger, c *fiber.Ctx, fields *logFields, start time.Time, rvr any) {
	if rvr != nil {
		err, ok := rvr.(error)
		if !ok {
			err = fmt.Errorf("%v", rvr)
		}

		status, message, code := getErrorMessage(err)
		fields.StatusCode = status
		fields.ErrorCode = code
		fields.Error = err
		fields.Stack = debug.Stack()

		c.Status(status).JSON(fiber.Map{
			"success":    false,
			"error_code": code,
			"error":      message,
			"requestId":  fields.ID,
		})
	}
}

func getErrorMessage(err error) (int, string, string) {
	isProduction := config.GetString(config.Debug) == "false"
	status := fiber.StatusInternalServerError
	message := "Internal Server Error"
	code := "unknown_error"

	if e, ok := err.(*fiber.Error); ok {
		status = e.Code
		if !isProduction || (status >= 400 && status < 500) {
			message = e.Message
		}
	} else if e, ok := err.(*apierror.Error); ok {
		status = e.Status
		code = e.Code
		if !isProduction || (status >= 400 && status < 500) || e.Show {
			message = e.Message
		}
	}

	return status, message, code
}

func logEvent(log zerolog.Logger, fields *logFields, rvr any) {
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
}
