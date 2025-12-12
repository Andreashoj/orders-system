package models

import "time"

type ShipmentStatus string

const (
	Processing     ShipmentStatus = "processing"
	OutForDelivery ShipmentStatus = "out_for_delivery"
	Delivered      ShipmentStatus = "delivered"
)

type Shipment struct {
	ID        int            `json:"id,omitempty"`
	OrderID   int            `json:"order_id,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	Status    ShipmentStatus `json:"status,omitempty"`
}
