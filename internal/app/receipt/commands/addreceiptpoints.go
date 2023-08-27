package commands

import (
	"context"

	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

type Calculator interface {
	Points(ctx context.Context, r receipt.Receipt) (*receipt.Points, error)
}

type PointsSaver struct {
	repo receipt.Repository
	calc Calculator
}

// NewAddReceiptPointsHandler Initializes an addReceiptPointsHandler
func NewSaverReceiptPoint(repo receipt.Repository, calc Calculator) PointsSaver {
	return PointsSaver{
		repo: repo,
		calc: calc,
	}
}

func (ps PointsSaver) SavePoints(ctx context.Context, r receipt.Receipt) (uuid.UUID, error) {
	points, err := ps.calc.Points(ctx, r)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := ps.repo.Save(ctx, *points)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
