package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/fr33dman/go-template/internal/heartbeat"
	"github.com/fr33dman/go-template/internal/server/handlers"
	"github.com/fr33dman/go-template/pkg/probes"
)

const shutdownTimeout = 10 * time.Second

func RunServer(
	ctx context.Context,
	host string,
	apiPort, healthcheckPort int,
	appProbes *probes.Probes,
	embeddingService handlers.EmbeddingService,
	log *slog.Logger,
) error {
	apiServer := NewAPIServer(host, apiPort, log)
	if err := apiServer.RegisterHandlers(
		handlers.NewEmbeddingHandler(embeddingService),
	); err != nil {
		return err
	}
	probesServer := heartbeat.NewProbesServer(host, healthcheckPort, appProbes)
	appProbes.Heartbeat(ctx)

	errCh := make(chan error, 2)

	go func() {
		log.InfoContext(
			ctx,
			"api server started",
			slog.String("host", host),
			slog.Int("port", apiPort),
		)
		if err := apiServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("api server: %w", err)
		}
	}()

	go func() {
		log.InfoContext(
			ctx,
			"probes server started",
			slog.String("host", host),
			slog.Int("port", healthcheckPort),
		)
		if err := probesServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("probes server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.InfoContext(ctx, "shutdown requested")
		return shutdown(apiServer, probesServer)
	case err := <-errCh:
		return errors.Join(err, shutdown(apiServer, probesServer))
	}
}

func shutdown(apiServer APIServer, probesServer heartbeat.ProbesServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return errors.Join(
		apiServer.Shutdown(ctx),
		probesServer.Shutdown(ctx),
	)
}
