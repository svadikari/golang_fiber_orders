package middleware

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Logger(c *fiber.Ctx) error {
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, ok := a.Value.Any().(*slog.Source)
				if ok {
					//Extract module name
					function := source.Function[strings.LastIndex(source.Function, "/")+1:]
					if strings.Contains(function, ".") {
						function = strings.Split(function, ".")[0]
					}
					//Extract file name
					fileName := source.File[strings.LastIndex(source.File, "/")+1:]
					return slog.String(slog.SourceKey, function+":"+fileName+":"+slog.AnyValue(source.Line).String())
				}
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))

	// Generate a request ID
	requestID := uuid.New().String()

	// Create a new context with the request ID
	cxt := context.WithValue(c.Context(), "requestID", requestID)

	c.SetUserContext(cxt)

	c.Locals("logger", logger.With(
		slog.String("requestID", requestID),
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
	))

	log := c.Locals("logger").(*slog.Logger)

	log.Info("Request",
		slog.String("method", c.Method()),
		slog.String("url", c.OriginalURL()),
		slog.String("ip", c.IP()),
	)

	// Proceed to the next middleware or handler
	return c.Next()
}
