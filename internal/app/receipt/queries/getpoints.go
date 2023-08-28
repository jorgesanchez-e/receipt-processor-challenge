package queries

import (
	"context"
	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

type PointsGetter struct {
	repo receipt.Repository
}

// NewReceiptPointsRequestHandler Handler Constructor.
func NewGetterReceiptPoints(repo receipt.Repository) PointsGetter {
	return PointsGetter{repo: repo}
}

// Handle Handlers the GetReceiptPoints request.
func (pg PointsGetter) GetPoints(ctx context.Context, id uuid.UUID) (*receipt.Points, error) {
	return pg.repo.Get(ctx, id)
}
