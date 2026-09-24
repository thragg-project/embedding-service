package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/fr33dman/go-template/internal/clients/embeddings"
	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
)

type (
	EmbeddingGateway interface {
		Models(ctx context.Context) ([]core.Model, error)
		Generate(ctx context.Context, req embeddings.GenerateRequest) (embeddings.Generated, error)
	}
	EmbeddingService struct {
		gateway EmbeddingGateway
	}
	GeneratedEmbeddings struct {
		Model      string
		Dimensions int
		Embeddings []core.Embedding
	}
)

func NewEmbeddingService(gateway EmbeddingGateway) *EmbeddingService {
	return &EmbeddingService{gateway: gateway}
}

func (s *EmbeddingService) Models(ctx context.Context) ([]core.Model, error) {
	return s.gateway.Models(ctx)
}

func (s *EmbeddingService) Generate(
	ctx context.Context,
	Model string,
	InputType core.InputType,
	Input []core.EmbeddingInput,
) (GeneratedEmbeddings, error) {
	if strings.TrimSpace(Model) == "" {
		return GeneratedEmbeddings{}, fmt.Errorf(
			"%w: model is required",
			apperrors.ErrInvalidArgument,
		)
	}
	if InputType != core.InputTypeQuery && InputType != core.InputTypeDocument {
		return GeneratedEmbeddings{}, fmt.Errorf(
			"%w: unsupported input_type %q",
			apperrors.ErrInvalidArgument,
			InputType,
		)
	}
	if len(Input) == 0 {
		return GeneratedEmbeddings{}, fmt.Errorf(
			"%w: input is required",
			apperrors.ErrInvalidArgument,
		)
	}

	for i, input := range Input {
		if strings.TrimSpace(input.ID) == "" {
			return GeneratedEmbeddings{}, fmt.Errorf(
				"%w: input[%d].id is required",
				apperrors.ErrInvalidArgument,
				i,
			)
		}
		if strings.TrimSpace(input.Text) == "" {
			return GeneratedEmbeddings{}, fmt.Errorf(
				"%w: input[%d].text is required",
				apperrors.ErrInvalidArgument,
				i,
			)
		}
	}

	generated, err := s.gateway.Generate(ctx, embeddings.GenerateRequest{
		Model:     Model,
		InputType: InputType,
		Input:     Input,
	})
	if err != nil {
		return GeneratedEmbeddings{}, err
	}
	return GeneratedEmbeddings{
		Model:      generated.Model,
		Dimensions: generated.Dimensions,
		Embeddings: generated.Embeddings,
	}, nil
}
