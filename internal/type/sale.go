package _type

type SaleItem struct {
	Item string `json:"Item"`
	Price float64 `json:"Price"`
	Amount float64 `json:"Amount"`
	City string `json:"City"`
	Count int `json:"Count"`
	Receipt string `json:"Receipt"`
	Seller string `json:"Seller"`
	SoldAt string `json:"SoldAt"`
}
