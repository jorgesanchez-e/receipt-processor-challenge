package storage

import (
	"context"

	"receipt-processor-challenge/internal/domain/receipt"
	"receipt-processor-challenge/internal/interfaceadapters/storage/memory"
)

type Service struct {
	Repository receipt.Repository
}

func New(ctx context.Context) Service {
	return Service{
		Repository: memory.New(ctx),
	}
}
