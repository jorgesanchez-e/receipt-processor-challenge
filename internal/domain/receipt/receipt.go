package receipt

import "time"

const (
	DatePurchaseFormat string = "2006-01-02"
	TimePurchaseFormat string = "15:04"
)

type Receipt struct {
	Retailer     string
	PurchaseDate time.Time
	PurchaseTime time.Time
	Items        []Item
	Total        float64
}

type Item struct {
	ShortDescription string
	Price            float64
}

type Points struct {
	Points int
}
