package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/server/handlers"
	"github.com/fr33dman/go-template/internal/services"
)

type fakeEmbeddingService struct{}

var _ handlers.EmbeddingService = (*fakeEmbeddingService)(nil)

func (s *fakeEmbeddingService) Models(ctx context.Context) ([]core.Model, error) {
	return []core.Model{
		{
			Name:           "bge-m3",
			Dimensions:     1024,
			MaxInputTokens: 8192,
			Multilingual:   true,
		},
	}, nil
}

func (s *fakeEmbeddingService) Generate(
	ctx context.Context,
	model string,
	inputType core.InputType,
	input []core.EmbeddingInput,
) (services.GeneratedEmbeddings, error) {
	if model == "missing" {
		return services.GeneratedEmbeddings{}, apperrors.ErrNotFound
	}
	return services.GeneratedEmbeddings{
		Model:      model,
		Dimensions: 2,
		Embeddings: []core.Embedding{
			{ID: input[0].ID, Vector: []float32{0.1, 0.2}},
		},
	}, nil
}

func TestAPIServerValidatesRequests(t *testing.T) {
	apiServer := NewAPIServer("127.0.0.1", 0, testLogger())
	err := apiServer.RegisterHandlers(handlers.NewEmbeddingHandler(&fakeEmbeddingService{}))
	if err != nil {
		t.Fatalf("register handlers: %v", err)
	}

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/v1/embeddings/generate",
		strings.NewReader(
			`{"model":"bge-m3","input_type":"query","input":[{"id":"chunk#1","text":"hello","extra":"rejected"}]}`,
		),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	apiServer.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestAPIServerMapsNotFoundError(t *testing.T) {
	apiServer := NewAPIServer("127.0.0.1", 0, testLogger())
	err := apiServer.RegisterHandlers(handlers.NewEmbeddingHandler(&fakeEmbeddingService{}))
	if err != nil {
		t.Fatalf("register handlers: %v", err)
	}

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/v1/embeddings/generate",
		strings.NewReader(
			`{"model":"missing","input_type":"query","input":[{"id":"chunk#1","text":"hello"}]}`,
		),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	apiServer.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
