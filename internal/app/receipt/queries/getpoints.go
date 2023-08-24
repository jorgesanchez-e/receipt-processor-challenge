package queries

import (
	"context"

	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

type GetReceiptPointsHandler struct {
	repo receipt.Repository
}

// NewGetCragRequestHandler Handler Constructor
func NewGetCragRequestHandler(repo receipt.Repository) GetReceiptPointsHandler {
	return GetReceiptPointsHandler{repo: repo}
}

// Handle Handlers the GetCragRequest query
func (h GetReceiptPointsHandler) Handle(ctx context.Context, id uuid.UUID) (*receipt.Points, error) {
	return h.repo.Get(ctx, id)
}
