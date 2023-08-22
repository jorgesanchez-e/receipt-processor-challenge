package memory_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"receipt-processor-challenge/internal/domain/receipt"
	. "receipt-processor-challenge/internal/interfaceadapters/storage/memory"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	mStorage *Engine

	wg sync.WaitGroup
)

type result struct {
	id     uuid.UUID
	points receipt.Points
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())

	mStorage = New(ctx, logrus.DebugLevel)
	code := m.Run()

	cancel()
	os.Exit(code)
}

func Test_Operations(t *testing.T) {
	ctx := context.Background()

	operations := 10

	results := newExectedResults(t, operations)
	saveFunctions := make([]func(*result), operations)

	trigger := make(chan struct{})
	for i := 0; i < operations; i++ {
		saveFunctions[i] = newSaveOperation(t, ctx, trigger, &results[i])
		go saveFunctions[i](&results[i])
	}

	close(trigger)
	wg.Wait()

	trigger = make(chan struct{})
	getFunctions := make([]func(result), operations)

	for i := 0; i < operations; i++ {
		getFunctions[i] = newGetOperations(t, ctx, trigger, results[i])
		go getFunctions[i](results[i])
	}

	close(trigger)
	wg.Wait()

	_, err := mStorage.Get(ctx, uuid.New())
	assert.Equal(t, err, ErrNotFound)
}

func newExectedResults(t *testing.T, n int) []result {
	t.Helper()

	results := make([]result, n)

	for i := 0; i < n; i++ {
		results[i].points = receipt.Points{Points: i}
	}

	return results
}

func newSaveOperation(t *testing.T, ctx context.Context, trigger chan struct{}, data *result) func(*result) {
	t.Helper()

	return func(r *result) {
		wg.Add(1)
		<-trigger

		uuid, err := mStorage.Save(ctx, r.points)
		assert.NoError(t, err)
		assert.NotNil(t, uuid)

		r.id = uuid
		wg.Done()
	}
}

func newGetOperations(t *testing.T, ctx context.Context, trigger chan struct{}, data result) func(result) {
	t.Helper()
	return func(r result) {
		wg.Add(1)
		<-trigger

		points, err := mStorage.Get(ctx, r.id)
		assert.NoError(t, err)
		assert.Equal(t, *points, r.points)

		wg.Done()
	}
}
