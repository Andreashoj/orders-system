package domain

type OrderProduct struct {
	OrderID   int `json:"order_id,omitempty"`
	ProductID int `json:"product_id,omitempty"`
	Quantity  int `json:"quantity,omitempty"`
}
