package commands

import (
	"context"

	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

type Calculator interface {
	Points(ctx context.Context, r receipt.Receipt) (*receipt.Points, error)
}

type SaveReceiptPointsHandler struct {
	repo receipt.Repository
	calc Calculator
}

// NewAddReceiptPointsHandler Initializes an addReceiptPointsHandler
func NewSaveReceiptPointsHandler(repo receipt.Repository, calc Calculator) SaveReceiptPointsHandler {
	return SaveReceiptPointsHandler{
		repo: repo,
		calc: calc,
	}
}

func (h SaveReceiptPointsHandler) Handle(ctx context.Context, r receipt.Receipt) (uuid.UUID, error) {
	points, err := h.calc.Points(ctx, r)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := h.repo.Save(ctx, *points)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
