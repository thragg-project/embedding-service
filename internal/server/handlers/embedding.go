package handlers

import (
	"context"
	"fmt"

	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/server/gen"
	"github.com/fr33dman/go-template/internal/services"
)

type EmbeddingService interface {
	Models(ctx context.Context) ([]core.Model, error)
	Generate(
		ctx context.Context,
		Model string,
		InputType core.InputType,
		Input []core.EmbeddingInput,
	) (services.GeneratedEmbeddings, error)
}

type EmbeddingHandler struct {
	service EmbeddingService
}

func NewEmbeddingHandler(service EmbeddingService) *EmbeddingHandler {
	return &EmbeddingHandler{service: service}
}

func (h *EmbeddingHandler) GetModels(
	ctx context.Context,
	_ gen.GetModelsRequestObject,
) (gen.GetModelsResponseObject, error) {
	models, err := h.service.Models(ctx)
	if err != nil {
		return nil, err
	}

	return gen.GetModels200JSONResponse{
		Models: mapModels(models),
	}, nil
}

func (h *EmbeddingHandler) GenerateVectors(
	ctx context.Context,
	request gen.GenerateVectorsRequestObject,
) (gen.GenerateVectorsResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("%w: request body is required", apperrors.ErrInvalidArgument)
	}

	input := make([]core.EmbeddingInput, len(request.Body.Input))
	for i, item := range request.Body.Input {
		input[i] = core.EmbeddingInput{
			ID:   item.Id,
			Text: item.Text,
		}
	}
	result, err := h.service.Generate(
		ctx,
		request.Body.Model,
		core.InputType(request.Body.InputType),
		input,
	)
	if err != nil {
		return nil, err
	}

	return gen.GenerateVectors200JSONResponse(mapGeneratedEmbeddings(result)), nil
}

func mapModels(source []core.Model) []gen.Model {
	models := make([]gen.Model, len(source))
	for i, model := range source {
		models[i] = gen.Model{
			Name:           model.Name,
			Dimensions:     model.Dimensions,
			MaxInputTokens: model.MaxInputTokens,
			Multilingual:   model.Multilingual,
		}
	}

	return models
}

func mapGeneratedEmbeddings(source services.GeneratedEmbeddings) gen.GeneratedEmbeddings {
	embeddings := make([]gen.EmbeddingData, len(source.Embeddings))
	for i, item := range source.Embeddings {
		embeddings[i] = gen.EmbeddingData{
			Id:     item.ID,
			Vector: item.Vector,
		}
	}

	return gen.GeneratedEmbeddings{
		Model:      source.Model,
		Dimensions: source.Dimensions,
		Embeddings: embeddings,
	}
}
