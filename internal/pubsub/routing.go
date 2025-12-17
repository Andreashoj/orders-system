package pubsub

const ExchangeOrderDirect string = "order_direct"

const (
	QueueTransaction      string = "order.transaction"
	QueueTransactionReply string = "order.transaction.reply"
	QueueShipping         string = "order.shipping"
	QueueShippingReply    string = "order.shipping.reply"
	QueueInventory        string = "order.inventory"
	QueueInventoryReply   string = "order.inventory.Reply"
)

const (
	TransactionKey      string = "transaction"
	TransactionKeyReply string = "transaction.reply"
	ShippingKey         string = "shipping"
	ShippingKeyReply    string = "shipping.reply"
	InventoryKey        string = "inventory"
	InventoryKeyReply   string = "inventory.reply"
)
