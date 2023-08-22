package interfaceadapters

import "receipt-processor-challenge/internal/domain/receipt"

type Service struct {
	Repository receipt.Repository
}
