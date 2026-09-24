package embeddings

import (
	"context"
	"errors"
	"testing"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
)

type fakeProvider struct{}

func (p fakeProvider) Generate(
	ctx context.Context,
	model core.Model,
	req GenerateRequest,
) (Generated, error) {
	return Generated{
		Model:      model.Name,
		Dimensions: model.Dimensions,
		Embeddings: []core.Embedding{
			{ID: req.Input[0].ID, Vector: []float32{0.1, 0.2}},
		},
	}, nil
}

func TestGatewayModelsReturnsOnlyEnabledAllowlist(t *testing.T) {
	gateway, err := NewGateway(Config{
		Models: []ModelConfig{
			{
				Name:           "disabled",
				Provider:       "provider",
				UpstreamModel:  "disabled",
				Dimensions:     2,
				MaxInputTokens: 128,
				Enabled:        false,
			},
			{
				Name:           "enabled",
				Provider:       "provider",
				UpstreamModel:  "enabled",
				Dimensions:     2,
				MaxInputTokens: 128,
				Enabled:        true,
			},
		},
	}, map[string]Provider{"provider": fakeProvider{}})
	if err != nil {
		t.Fatalf("new gateway: %v", err)
	}

	models, err := gateway.Models(context.Background())
	if err != nil {
		t.Fatalf("models: %v", err)
	}
	if len(models) != 1 || models[0].Name != "enabled" {
		t.Fatalf("expected only enabled model, got %+v", models)
	}
}

func TestGatewayGenerateRejectsUnknownOrDisabledModel(t *testing.T) {
	gateway, err := NewGateway(Config{
		Models: []ModelConfig{
			{
				Name:           "disabled",
				Provider:       "provider",
				UpstreamModel:  "disabled",
				Dimensions:     2,
				MaxInputTokens: 128,
				Enabled:        false,
			},
			{
				Name:           "enabled",
				Provider:       "provider",
				UpstreamModel:  "enabled",
				Dimensions:     2,
				MaxInputTokens: 128,
				Enabled:        true,
			},
		},
	}, map[string]Provider{"provider": fakeProvider{}})
	if err != nil {
		t.Fatalf("new gateway: %v", err)
	}

	_, err = gateway.Generate(context.Background(), GenerateRequest{
		Model:     "disabled",
		InputType: core.InputTypeQuery,
		Input:     []core.EmbeddingInput{{ID: "chunk#1", Text: "hello"}},
	})
	if !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
