package memory

import (
	"context"
	"receipt-processor-challenge/internal/domain/receipt"
	"testing"

	"github.com/google/uuid"
)

var id uuid.UUID

func BenchmarkSave(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	store := New(ctx)

	var newID uuid.UUID
	for i := 0; i < b.N; i++ {
		newID, _ = store.Save(ctx, receipt.Points{Points: 10})
	}

	id = newID
	cancel()
}
