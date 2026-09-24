package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echomiddleware "github.com/oapi-codegen/echo-middleware"

	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/server/gen"
	"github.com/fr33dman/go-template/internal/server/handlers"
)

type APIServer struct {
	Host string
	Port int
	log  *slog.Logger
	*echo.Echo
}

func NewAPIServer(
	host string,
	port int,
	log *slog.Logger,
) APIServer {
	if log == nil {
		log = slog.Default()
	}
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return APIServer{
		Host: host,
		Port: port,
		log:  log,
		Echo: e,
	}
}

func (s *APIServer) RegisterHandlers(
	embeddings *handlers.EmbeddingHandler,
) error {
	s.HTTPErrorHandler = errorHandler(s.log)
	s.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))
	s.Use(requestLogger(s.log))
	s.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))

	swagger, err := gen.GetSwagger()
	if err != nil {
		return fmt.Errorf("load OpenAPI spec: %w", err)
	}
	s.Use(echomiddleware.OapiRequestValidator(swagger))

	gen.RegisterHandlersWithBaseURL(
		s,
		gen.NewStrictHandler(NewAPIHandler(embeddings), nil),
		"/api/v1/embeddings",
	)
	return nil
}

func (s *APIServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	return s.Echo.Start(address)
}

func requestLogger(log *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRemoteIP: true,
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("remote_ip", values.RemoteIP),
				slog.String("method", values.Method),
				slog.String("uri", values.URI),
				slog.Int("status", values.Status),
				slog.Duration("latency", values.Latency),
			}
			if values.Error != nil {
				attrs = append(attrs, slog.Any("error", values.Error))
			}

			log.LogAttrs(c.Request().Context(), slog.LevelInfo, "http request", attrs...)
			return nil
		},
	})
}

func errorHandler(log *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		status := http.StatusInternalServerError
		message := http.StatusText(status)

		var httpErr *echo.HTTPError
		switch {
		case errors.Is(err, apperrors.ErrInvalidArgument):
			status = http.StatusBadRequest
			message = err.Error()
		case errors.Is(err, apperrors.ErrNotFound):
			status = http.StatusNotFound
			message = "not found"
		case errors.Is(err, apperrors.ErrUnavailable):
			status = http.StatusServiceUnavailable
			message = "service unavailable"
		case errors.As(err, &httpErr):
			status = httpErr.Code
			message = fmt.Sprint(httpErr.Message)
		}

		if status >= http.StatusInternalServerError {
			log.ErrorContext(
				c.Request().Context(),
				"http error",
				slog.Int("status", status),
				slog.Any("error", err),
			)
		}

		if c.Response().Committed {
			return
		}

		if err := c.JSON(status, gen.Error{Message: message}); err != nil {
			log.ErrorContext(c.Request().Context(), "write error response", slog.Any("error", err))
		}
	}
}
