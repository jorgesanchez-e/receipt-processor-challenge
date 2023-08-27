package calculator

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"receipt-processor-challenge/internal/domain/receipt"
)

/*
*
RULES
One point for every alphanumeric character in the retailer name.
50 points if the total is a round dollar amount with no cents.
25 points if the total is a multiple of 0.25.
5 points for every two items on the receipt.
If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
6 points if the day in the purchase date is odd.
10 points if the time of purchase is after 2:00pm and before 4:00pm.
**
*/
type Calculator struct{}

func New(ctx context.Context) Calculator {
	return Calculator{}
}

func (c Calculator) Points(ctx context.Context, r receipt.Receipt) (*receipt.Points, error) {
	points := retailerNamePoints(r.Retailer)
	points += roundDollarPoints(r.Total)
	points += multipleOf25CentsPoints(r.Total)
	points += itemsPoints(r.Items)
	points += trimmedDescriptionPoints(r.Items)
	points += oddPurchaseDayPoints(r.PurchaseDate)
	points += timePurchasePoints(r.PurchaseTime)

	return &receipt.Points{Points: points}, nil
}

func retailerNamePoints(name string) int {
	points := 0
	for _, c := range name {
		if unicode.IsLetter(c) || unicode.IsNumber(c) {
			points++
		}
	}

	return points
}

func roundDollarPoints(total float64) int {
	if total == math.Trunc(total) {
		return 50
	}

	return 0
}

func multipleOf25CentsPoints(total float64) int {
	const frac float64 = 0.25
	if math.Mod(float64(total), frac) == 0 {
		return 25
	}

	return 0
}

func itemsPoints(items []receipt.Item) int {
	return (len(items) / 2) * 5
}

func trimmedDescriptionPoints(items []receipt.Item) int {
	total := 0

	for _, item := range items {
		x := strings.TrimSpace(item.ShortDescription)
		l := len(x)
		if l%3 == 0 {
			test := float64(item.Price * 0.2)

			if test == math.Trunc(test) {
				total += int(test)
				continue
			}

			fl := fmt.Sprintf("%f", test)
			text := strings.Split(fl, ".")
			itext, err := strconv.Atoi(text[0])
			if err == nil {
				total += itext + 1
			}
		}

	}

	return total
}

func oddPurchaseDayPoints(date time.Time) int {
	if date.Day()%2 == 0 {
		return 0
	}
	return 6
}

func timePurchasePoints(t time.Time) int {
	if t.Hour() >= 14 && t.Hour() <= 16 {
		return 10
	}

	return 0
}
