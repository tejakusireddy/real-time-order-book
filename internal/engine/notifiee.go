package engine

type Notifiee interface {
	OnMatch(match MatchEvent)
	OnAdd(order Order)
}

type MatchEvent struct {
	BuyOrderID   string  `json:"buy_order_id"`
	SellOrderID  string  `json:"sell_order_id"`
	Price        float64 `json:"price"`
	Quantity     float64 `json:"quantity"`
	TimestampUTC string  `json:"timestamp"`
}
