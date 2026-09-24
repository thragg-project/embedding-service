package embeddings

import (
	"context"
	"fmt"
	"slices"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
)

type (
	Gateway struct {
		providers map[string]Provider
		models    map[string]core.Model
	}
)

func NewGateway(cfg Config, providers map[string]Provider) (*Gateway, error) {
	models := make(map[string]core.Model, len(cfg.Models))
	for _, model := range cfg.Models {
		if !model.Enabled {
			continue
		}
		if _, exists := providers[model.Provider]; !exists {
			return nil, fmt.Errorf(
				"provider %q for model %q is not configured",
				model.Provider,
				model.Name,
			)
		}
		models[model.Name] = core.Model{
			Name:           model.Name,
			Provider:       model.Provider,
			UpstreamModel:  model.UpstreamModel,
			Dimensions:     model.Dimensions,
			MaxInputTokens: model.MaxInputTokens,
			Multilingual:   model.Multilingual,
			Enabled:        model.Enabled,
		}
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("embeddings config has no enabled models")
	}

	return &Gateway{
		providers: providers,
		models:    models,
	}, nil
}

func (g *Gateway) Models(_ context.Context) ([]core.Model, error) {
	models := make([]core.Model, 0, len(g.models))
	for _, model := range g.models {
		models = append(models, model)
	}
	slices.SortFunc(models, func(a, b core.Model) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	return models, nil
}

func (g *Gateway) Generate(
	ctx context.Context,
	req GenerateRequest,
) (Generated, error) {
	model, exists := g.models[req.Model]
	if !exists {
		return Generated{}, fmt.Errorf(
			"%w: embedding model %q",
			apperrors.ErrNotFound,
			req.Model,
		)
	}

	provider, exists := g.providers[model.Provider]
	if !exists {
		return Generated{}, fmt.Errorf(
			"provider %q for model %q is not configured",
			model.Provider,
			model.Name,
		)
	}

	result, err := provider.Generate(ctx, model, GenerateRequest{
		Model:     req.Model,
		InputType: req.InputType,
		Input:     req.Input,
	})
	if err != nil {
		return Generated{}, err
	}

	return Generated{
		Model:      result.Model,
		Dimensions: result.Dimensions,
		Embeddings: result.Embeddings,
	}, nil
}
