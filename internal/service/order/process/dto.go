package process

type ProcessMessageDto struct {
	OrderId string `json:"order_id"`

	PaymentResponse *PaymentResponse `json:"payment"`
	OrderResponse   *OrderResponse   `json:"order"`
}

type PaymentResponse struct {
	PaymentId string `json:"id"`
	State     string `json:"state"`
}

type OrderResponse struct {
	State string `json:"state"`
}
