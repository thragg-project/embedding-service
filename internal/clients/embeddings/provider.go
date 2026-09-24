package embeddings

import (
	"context"

	"github.com/fr33dman/go-template/internal/core"
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
