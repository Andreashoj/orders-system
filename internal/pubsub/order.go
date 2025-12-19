package pubsub

type TransactionReplyMessage struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	Success       bool   `json:"success,omitempty"`
}

type OrderProcessRequirements struct {
	TransactionComplete bool
	ShipmentComplete    bool
	InventoryComplete   bool
}

type PubTransaction struct {
	OrderId string `json:"order_id"`
}
