package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/fr33dman/go-template/internal/app"
	"github.com/fr33dman/go-template/internal/build"
	"github.com/fr33dman/go-template/internal/logger"
	"github.com/fr33dman/go-template/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appLogger, _ := logger.New("", build.Version)
	slog.SetDefault(appLogger)

	cfg, err := app.NewConfig()
	if err != nil {
		logger.Fatal(ctx, appLogger, "load config", err)
	}

	appLogger, err = logger.New(cfg.Logger.Level, build.Version)
	if err != nil {
		logger.Fatal(ctx, appLogger, "create logger", err)
	}
	slog.SetDefault(appLogger)

	container, err := app.NewContainer(ctx, cfg)
	if err != nil {
		logger.Fatal(ctx, appLogger, "create container", err)
	}
	defer container.Close()

	if err := server.RunServer(
		ctx,
		cfg.Server.Host,
		cfg.Server.APIPort,
		cfg.Server.HealthcheckPort,
		container.Probes,
		container.EmbeddingService,
		appLogger,
	); err != nil {
		logger.Fatal(ctx, appLogger, "run server", err)
	}
}
