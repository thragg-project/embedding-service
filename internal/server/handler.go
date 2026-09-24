package server

import (
	"context"
	"fmt"

	"github.com/fr33dman/go-template/internal/server/gen"
)

type EmbeddingHandler interface {
	GenerateVectors(
		ctx context.Context,
		request gen.GenerateVectorsRequestObject,
	) (gen.GenerateVectorsResponseObject, error)
	GetModels(
		ctx context.Context,
		request gen.GetModelsRequestObject,
	) (gen.GetModelsResponseObject, error)
}

type APIHandler struct {
	embeddingHandler EmbeddingHandler
}

func NewAPIHandler(
	embeddingHandler EmbeddingHandler,
) *APIHandler {
	return &APIHandler{
		embeddingHandler: embeddingHandler,
	}
}

var _ gen.StrictServerInterface = (*APIHandler)(nil)

func (h *APIHandler) GenerateVectors(
	ctx context.Context,
	request gen.GenerateVectorsRequestObject,
) (gen.GenerateVectorsResponseObject, error) {
	if h.embeddingHandler == nil {
		return nil, notImplemented("embedding handler not implemented")
	}
	return h.embeddingHandler.GenerateVectors(ctx, request)
}

func (h *APIHandler) GetModels(
	ctx context.Context,
	request gen.GetModelsRequestObject,
) (gen.GetModelsResponseObject, error) {
	if h.embeddingHandler == nil {
		return nil, notImplemented("embedding handler not implemented")
	}
	return h.embeddingHandler.GetModels(ctx, request)
}

func notImplemented(name string) error {
	return fmt.Errorf("%s is not configured", name)
}
