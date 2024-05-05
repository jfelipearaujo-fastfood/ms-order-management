package order_entity

type Item struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	UnitPrice float64 `json:"unit_price"`
	Quantity  int     `json:"quantity"`
}

func NewItem(id string, name string, unitPrice float64, quantity int) Item {
	return Item{
		Id:        id,
		Name:      name,
		UnitPrice: unitPrice,
		Quantity:  quantity,
	}
}
