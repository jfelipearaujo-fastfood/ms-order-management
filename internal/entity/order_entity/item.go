package order_entity

type Item struct {
	Id        string  `json:"id"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  int     `json:"quantity"`
}

func NewItem(id string, unitPrice float64, quantity int) Item {
	return Item{
		Id:        id,
		UnitPrice: unitPrice,
		Quantity:  quantity,
	}
}
