package db

type Product struct {
	ID       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Price    int    `json:"price,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}
