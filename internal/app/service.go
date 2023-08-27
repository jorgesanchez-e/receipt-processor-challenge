package app

import (
	"receipt-processor-challenge/internal/app/receipt/commands"
	"receipt-processor-challenge/internal/app/receipt/queries"
	"receipt-processor-challenge/internal/domain/receipt"
)

// Services contains all exposed services of the application layer
type Service struct {
	commands.PointsSaver
	queries.PointsGetter
}

// NewServices Bootstraps Application Layer dependencies
func NewServices(repo receipt.Repository, calc commands.Calculator) Service {
	return Service{
		commands.NewSaverReceiptPoint(repo, calc),
		queries.NewGetterReceiptPoints(repo),
	}
}
