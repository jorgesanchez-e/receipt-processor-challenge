package queries

import (
	"context"

	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

type GetReceiptPointsHandler struct {
	repo receipt.Repository
}

// NewReceiptPointsRequestHandler Handler Constructor
func NewReceiptPointsRequestHandler(repo receipt.Repository) GetReceiptPointsHandler {
	return GetReceiptPointsHandler{repo: repo}
}

// Handle Handlers the GetReceiptPoints request.
func (h GetReceiptPointsHandler) Handle(ctx context.Context, id uuid.UUID) (*receipt.Points, error) {
	return h.repo.Get(ctx, id)
}
