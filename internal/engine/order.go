package engine

import (
	"errors"
	"strings"
	"time"
)

type OrderType string
type Side string

const (
	LimitOrder  OrderType = "limit"
	MarketOrder OrderType = "market"

	Buy  Side = "buy"
	Sell Side = "sell"
)

type Order struct {
	ID        string
	Timestamp time.Time
	Price     float64
	Quantity  float64
	Type      OrderType
	Side      Side
}

func (s Side) IsValid() bool {
	switch strings.ToLower(string(s)) {
	case "buy", "sell":
		return true
	default:
		return false
	}
}

func (o OrderType) IsValid() bool {
	switch strings.ToLower(string(o)) {
	case "limit", "market":
		return true
	default:
		return false
	}
}

func (o *Order) Validate() error {
	if o.ID == "" {
		return errors.New("missing order ID")
	}
	if !o.Type.IsValid() {
		return errors.New("invalid order type")
	}
	if !o.Side.IsValid() {
		return errors.New("invalid side")
	}
	if o.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if o.Type == LimitOrder && o.Price <= 0 {
		return errors.New("limit order must have positive price")
	}
	return nil
}
