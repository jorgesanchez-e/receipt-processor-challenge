package calculator

import (
	"fmt"
	"math"
	"receipt-processor-challenge/internal/domain/receipt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	oddDayPoints        int     = 6
	tPurchasePoints     int     = 10
	roundTotalPoints    int     = 50
	multipleOf25Points  int     = 25
	trimmedFactorPoints float64 = 0.2
	itemsFactorPoints   int     = 5
	numItems            int     = 2
)

/*
RULES
  - One point for every alphanumeric character in the retailer name.
  - 50 points if the total is a round dollar amount with no cents.
  - 25 points if the total is a multiple of 0.25.
  - 5 points for every two items on the receipt.
  - If the trimmed length of the item description is a multiple of 3,
    multiply the price by 0.2 and round up to the nearest integer. The result
    is the number of points earned.
  - 6 points if the day in the purchase date is odd.
  - 10 points if the time of purchase is after 2:00pm and before 4:00pm.
*/
type Calculator struct{}

func New() Calculator {
	return Calculator{}
}

func (c Calculator) Points(rcpt receipt.Receipt) (*receipt.Points, error) {
	points := retailerNamePoints(rcpt.Retailer)
	points += roundDollarPoints(rcpt.Total)
	points += multipleOf25CentsPoints(rcpt.Total)
	points += itemsPoints(rcpt.Items)
	points += trimmedDescriptionPoints(rcpt.Items)
	points += oddPurchaseDayPoints(rcpt.PurchaseDate)
	points += timePurchasePoints(rcpt.PurchaseTime)

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
		return roundTotalPoints
	}

	return 0
}

func multipleOf25CentsPoints(total float64) int {
	const frac float64 = 0.25
	if math.Mod(total, frac) == 0 {
		return multipleOf25Points
	}

	return 0
}

func itemsPoints(items []receipt.Item) int {
	return (len(items) / numItems) * itemsFactorPoints
}

func trimmedDescriptionPoints(items []receipt.Item) int {
	total := 0

	for _, item := range items {
		x := strings.TrimSpace(item.ShortDescription)
		l := len(x)

		if l%3 == 0 {
			test := item.Price * trimmedFactorPoints

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

	return oddDayPoints
}

func timePurchasePoints(t time.Time) int {
	if t.Hour() >= 14 && t.Hour() <= 16 {
		return tPurchasePoints
	}

	return 0
}
