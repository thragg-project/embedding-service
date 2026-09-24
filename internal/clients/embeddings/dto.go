package embeddings

import "github.com/fr33dman/go-template/internal/core"

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
