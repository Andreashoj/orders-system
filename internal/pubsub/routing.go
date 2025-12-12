package pubsub

type ExchangeKey string

const (
	ExchangeOrderDirect ExchangeKey = "order_direct"
)

type QueueName string

const (
	QueueTransaction QueueName = "order.transaction"
	QueueShipping    QueueName = "order.shipping"
	QueueInventory   QueueName = "order.inventory"
)
