package app

import (
	"receipt-processor-challenge/internal/app/receipt/commands"
	"receipt-processor-challenge/internal/app/receipt/queries"
	"receipt-processor-challenge/internal/domain/receipt"
)

/*********************************************************************
* WE NEED TO MOVE THIS INTERFACES TO OTHER PACKAGE
type CreateAddReceiptPointsHandler interface {
	Handle(ctx context.Context, r receipt.Receipt) (uuid.UUID, error)
}

type GetReceiptPointsHandler interface {
	Handle(ctx context.Context, id uuid.UUID) (*receipt.Points, error)
}
**********************************************************************/

// Services contains all exposed services of the application layer
type Service struct {
	commands.AddReceiptPointsHandler
	queries.GetReceiptPointsHandler
}

// NewServices Bootstraps Application Layer dependencies
func NewServices(repo receipt.Repository, calc commands.Calculator) Service {
	return Service{
		commands.NewAddReceiptPointsHandler(repo, calc),
		queries.NewGetCragRequestHandler(repo),
	}
}
