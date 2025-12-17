package pubsub

const ExchangeOrderDirect string = "order_direct"

const (
	QueueTransaction string = "order.transaction"
	QueueShipping    string = "order.shipping"
	QueueInventory   string = "order.inventory"
)

const (
	TransactionKey string = "transaction"
	ShippingKey    string = "shipping"
	InventoryKey   string = "inventory"
)
