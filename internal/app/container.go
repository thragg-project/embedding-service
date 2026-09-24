package app

import (
	"context"
	"fmt"

	"github.com/fr33dman/go-template/internal/clients/embeddings"
	"github.com/fr33dman/go-template/internal/clients/embeddings/openai"
	"github.com/fr33dman/go-template/internal/heartbeat"
	"github.com/fr33dman/go-template/internal/services"
	"github.com/fr33dman/go-template/pkg/probes"
)

type (
	Container struct {
		// Probes
		Probes *probes.Probes

		// Services
		EmbeddingService *services.EmbeddingService
	}
)

func NewContainer(ctx context.Context, cfg Config) (Container, error) {
	_ = ctx

	providers, err := newEmbeddingProviders(cfg.Embeddings)
	if err != nil {
		return Container{}, err
	}
	embeddingGateway, err := embeddings.NewGateway(cfg.Embeddings, providers)
	if err != nil {
		return Container{}, err
	}

	// Setup Probes
	appProbes := heartbeat.NewProbes()

	return Container{
		Probes:           appProbes,
		EmbeddingService: services.NewEmbeddingService(embeddingGateway),
	}, nil
}

func (c Container) Close() {}

func newEmbeddingProviders(cfg embeddings.Config) (map[string]embeddings.Provider, error) {
	providers := make(map[string]embeddings.Provider, len(cfg.Providers))
	for _, providerCfg := range cfg.Providers {
		switch providerCfg.Type {
		case embeddings.ProviderTypeOpenAICompatible:
			provider, err := openai.NewClient(providerCfg)
			if err != nil {
				return nil, err
			}
			providers[providerCfg.Name] = provider
		default:
			return nil, fmt.Errorf("unsupported embeddings provider type %q", providerCfg.Type)
		}
	}

	return providers, nil
}
