package commands

import (
	"context"

	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

type Calculator interface {
	Points(ctx context.Context, r receipt.Receipt) (*receipt.Points, error)
}

type AddReceiptPointsHandler struct {
	repo receipt.Repository
	calc Calculator
}

// NewAddReceiptPointsHandler Initializes an addReceiptPointsHandler
func NewAddReceiptPointsHandler(repo receipt.Repository, calc Calculator) AddReceiptPointsHandler {
	return AddReceiptPointsHandler{
		repo: repo,
		calc: calc,
	}
}

func (h AddReceiptPointsHandler) Handle(ctx context.Context, r receipt.Receipt) (uuid.UUID, error) {
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
