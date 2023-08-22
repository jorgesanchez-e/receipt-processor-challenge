package receipt

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, points Points) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*Points, error)
}
