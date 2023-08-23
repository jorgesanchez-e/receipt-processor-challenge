package memory_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"receipt-processor-challenge/internal/domain/receipt"
	. "receipt-processor-challenge/internal/interfaceadapters/storage/memory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var mStorage *Engine

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())

	mStorage = New(ctx)
	code := m.Run()

	cancel()
	os.Exit(code)
}

func Test_Save(t *testing.T) {
	cases := []struct {
		name           string
		contextBuilder func() (context.Context, context.CancelFunc)
		input          []receipt.Points
		expectedResult []error
	}{
		{
			name: "normal-case",
			contextBuilder: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			input: []receipt.Points{
				{Points: 10},
				{Points: 20},
				{Points: 30},
				{Points: 40},
				{Points: 50},
			},
			expectedResult: []error{
				nil,
				nil,
				nil,
				nil,
				nil,
			},
		}, {
			name: "time-out-case",
			contextBuilder: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))

				return ctx, cancel
			},
			input: []receipt.Points{
				{Points: 10},
			},
			expectedResult: []error{context.DeadlineExceeded},
		},
	}

	for _, c := range cases {
		ctx, cancel := c.contextBuilder()
		input := c.input
		expectedResult := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			defer cancel()

			for i := 0; i < len(input); i++ {
				uuid, err := mStorage.Save(ctx, input[i])

				assert.NotNil(t, uuid)
				assert.Equal(t, expectedResult[i], err)
			}
		})
	}
}

func Test_Get(t *testing.T) {
	ctx := context.Background()

	point, err := mStorage.Get(ctx, uuid.New())
	assert.Equal(t, errors.New("points not found"), err)
	assert.Nil(t, point)

	newPoint := receipt.Points{Points: 199}
	newID, err := mStorage.Save(ctx, newPoint)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, newID)

	lastPoint, err := mStorage.Get(ctx, newID)
	assert.NoError(t, err)
	assert.Equal(t, newPoint.Points, lastPoint.Points)
}
