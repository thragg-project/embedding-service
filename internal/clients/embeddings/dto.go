package embeddings

import "github.com/thragg-project/embedding-service/internal/core"

type (
	GenerateRequest struct {
		Model     string
		InputType core.InputType
		Input     []core.EmbeddingInput
	}
	Generated struct {
		Model      string
		Dimensions int
		Embeddings []core.Embedding
	}
)
