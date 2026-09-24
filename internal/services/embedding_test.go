package services_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/fr33dman/go-template/internal/clients/embeddings"
	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
	"github.com/fr33dman/go-template/internal/services"
	"github.com/fr33dman/go-template/internal/services/mocks"
)

func TestEmbeddingServiceGenerateDelegatesToGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	gateway := mocks.NewMockEmbeddingGateway(ctrl)
	service := services.NewEmbeddingService(gateway)
	input := []core.EmbeddingInput{
		{ID: "chunk#1", Text: "hello"},
	}
	req := embeddings.GenerateRequest{
		Model:     "bge-m3",
		InputType: core.InputTypeQuery,
		Input:     input,
	}
	generated := embeddings.Generated{
		Model:      req.Model,
		Dimensions: 2,
		Embeddings: []core.Embedding{
			{Vector: []float32{0.1, 0.2}},
		},
	}
	expected := services.GeneratedEmbeddings(generated)

	gateway.EXPECT().Generate(ctx, req).Return(generated, nil)

	got, err := service.Generate(ctx, req.Model, req.InputType, req.Input)
	if err != nil {
		t.Fatalf("generate embeddings: %v", err)
	}
	if got.Model != expected.Model || got.Dimensions != expected.Dimensions ||
		len(got.Embeddings) != len(expected.Embeddings) {
		t.Fatalf("expected %+v, got %+v", expected, got)
	}
}

func TestEmbeddingServiceGenerateRejectsEmptyInputText(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := services.NewEmbeddingService(mocks.NewMockEmbeddingGateway(ctrl))

	_, err := service.Generate(
		context.Background(),
		"bge-m3",
		core.InputTypeQuery,
		[]core.EmbeddingInput{{ID: "chunk#1", Text: "  "}},
	)
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}
