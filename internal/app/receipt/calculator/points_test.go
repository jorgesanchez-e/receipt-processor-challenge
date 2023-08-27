package calculator

import (
	"context"
	"receipt-processor-challenge/internal/domain/receipt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RetailerNamePoints(t *testing.T) {
	cases := []struct {
		name           string
		retailName     string
		expectedResult int
	}{
		{
			name:           "case-Target",
			retailName:     "Target",
			expectedResult: 6,
		},
		{
			name:           "case-M&M Corner Market",
			retailName:     "M&M Corner Market",
			expectedResult: 14,
		},
	}

	for _, c := range cases {
		rname := c.retailName
		expectedPoints := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			points := retailerNamePoints(rname)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_roundDollarPoints(t *testing.T) {
	cases := []struct {
		name           string
		total          float64
		expectedResult int
	}{
		{
			name:           "total-rounded",
			total:          9.00,
			expectedResult: 50,
		},
		{
			name:           "total-no-rounded",
			total:          9.25,
			expectedResult: 0,
		},
	}

	for _, c := range cases {
		total := c.total
		expectedPoints := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			points := roundDollarPoints(total)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_multipleOf25CentsPoints(t *testing.T) {
	cases := []struct {
		name           string
		total          float64
		expectedResult int
	}{
		{
			name:           "is-multiple",
			total:          9.0,
			expectedResult: 25,
		},
		{
			name:           "is-not-multiple",
			total:          9.1,
			expectedResult: 0,
		},
	}

	for _, c := range cases {
		total := c.total
		expectedPoints := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			points := multipleOf25CentsPoints(total)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_itemsPoints(t *testing.T) {
	cases := []struct {
		name           string
		items          []receipt.Item
		expectedResult int
	}{
		{
			name:           "test-0-tems",
			items:          []receipt.Item{},
			expectedResult: 0,
		},
		{
			name: "test-5-tems",
			items: []receipt.Item{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            6.49,
				}, {
					ShortDescription: "Emils Cheese Pizza",
					Price:            12.25,
				}, {
					ShortDescription: "Knorr Creamy Chicken",
					Price:            1.26,
				}, {
					ShortDescription: "Doritos Nacho Cheese",
					Price:            3.35,
				}, {
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            12.00,
				},
			},
			expectedResult: 10,
		},
	}

	for _, c := range cases {
		items := c.items
		expectedPoints := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			points := itemsPoints(items)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_trimmedDescriptionPoints(t *testing.T) {
	cases := []struct {
		name           string
		items          []receipt.Item
		expectedResult int
	}{
		{
			name: "Klarbrunn-case",
			items: []receipt.Item{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            6.49,
				}, {
					ShortDescription: "Emils Cheese Pizza",
					Price:            12.25,
				}, {
					ShortDescription: "Knorr Creamy Chicken",
					Price:            1.26,
				}, {
					ShortDescription: "Doritos Nacho Cheese",
					Price:            3.35,
				}, {
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            12.00,
				},
			},
			expectedResult: 6,
		},
	}

	for _, c := range cases {
		items := c.items
		expectedPoints := c.expectedResult

		t.Run(c.name, func(t *testing.T) {
			points := trimmedDescriptionPoints(items)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_oddPurchaseDayPoints(t *testing.T) {
	cases := []struct {
		name           string
		date           string
		expectedResult int
	}{
		{
			name:           "odd-day-case",
			date:           "2022-01-01",
			expectedResult: 6,
		},
		{
			name:           "even-day-case",
			date:           "2022-01-20",
			expectedResult: 0,
		},
	}

	for _, c := range cases {
		expectedPoints := c.expectedResult
		date, err := time.Parse(receipt.DatePurchaseFormat, c.date)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(c.name, func(t *testing.T) {
			points := oddPurchaseDayPoints(date)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_timePurchasePoints(t *testing.T) {
	cases := []struct {
		name           string
		tdate          string
		expectedResult int
	}{
		{
			name:           "hour-in-range",
			tdate:          "14:00",
			expectedResult: 10,
		},
		{
			name:           "hour-out-range",
			tdate:          "13:59",
			expectedResult: 0,
		},
	}

	for _, c := range cases {
		expectedPoints := c.expectedResult
		ctime, err := time.Parse(receipt.TimePurchaseFormat, c.tdate)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(c.name, func(t *testing.T) {
			points := timePurchasePoints(ctime)

			assert.Equal(t, expectedPoints, points)
		})
	}
}

func Test_Points(t *testing.T) {
	cases := []struct {
		name           string
		receipt        receipt.Receipt
		expectedResult *receipt.Points
		expectedError  error
	}{
		{
			name: "Target-case",
			receipt: receipt.Receipt{
				Retailer:     "Target",
				PurchaseDate: purchaseDate(t, "2022-01-01"),
				PurchaseTime: purchaseTime(t, "13:01"),
				Items: []receipt.Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            6.49,
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            12.25,
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            1.26,
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            3.35,
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            12.00,
					},
				},
				Total: 35.35,
			},
			expectedResult: &receipt.Points{Points: 28},
			expectedError:  nil,
		},
		{
			name: "M&M Corner Market-case",
			receipt: receipt.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: purchaseDate(t, "2022-03-20"),
				PurchaseTime: purchaseTime(t, "14:33"),
				Items: []receipt.Item{
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            12.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            12.25,
					},
				},
				Total: 9.00,
			},
			expectedResult: &receipt.Points{Points: 109},
			expectedError:  nil,
		},
	}

	for _, c := range cases {
		ctx := context.Background()
		r := c.receipt
		expectedPoints := c.expectedResult
		expectedError := c.expectedError
		cal := New(ctx)

		t.Run(c.name, func(t *testing.T) {
			points, err := cal.Points(ctx, r)

			assert.Equal(t, expectedPoints, points)
			assert.Equal(t, expectedError, err)
		})
	}
}

func purchaseDate(t *testing.T, date string) time.Time {
	tdate, err := time.Parse(receipt.DatePurchaseFormat, date)
	if err != nil {
		t.Fatal(err)
	}

	return tdate
}

func purchaseTime(t *testing.T, date string) time.Time {
	tdate, err := time.Parse(receipt.TimePurchaseFormat, date)
	if err != nil {
		t.Fatal(err)
	}

	return tdate
}
