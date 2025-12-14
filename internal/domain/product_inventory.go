package domain

import "time"

type ProductInventory struct {
	ID             int       `json:"id,omitempty"`
	ProductID      int       `json:"product_id,omitempty"`
	Quantity       int       `json:"quantity,omitempty"`
	NextShipmentAt time.Time `json:"next_shipment_at"`
}
