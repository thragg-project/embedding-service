package embeddings

import (
	"context"

	"github.com/thragg-project/embedding-service/internal/core"
)

type (
	Provider interface {
		Generate(
			ctx context.Context,
			model core.Model,
			req GenerateRequest,
		) (Generated, error)
	}
)
