package receipt

type Receipt struct {
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Items        []Item
	Total        float32
}

type Item struct {
	ShortDescription string
	Price            float32
}

type Points struct {
	Points int
}
