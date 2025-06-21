package engine

import (
	"testing"
	"time"
)

func TestBuyMatchesSell(t *testing.T) {
	ob := NewOrderBook()

	sell := Order{
		ID:        "sell-1",
		Timestamp: time.Now(),
		Price:     100.0,
		Quantity:  1.0,
		Type:      LimitOrder,
		Side:      Sell,
	}
	ob.AddOrder(sell)

	buy := Order{
		ID:        "buy-1",
		Timestamp: time.Now(),
		Price:     105.0, // Matches sell
		Quantity:  1.0,
		Type:      LimitOrder,
		Side:      Buy,
	}
	ob.AddOrder(buy)

	if len(ob.sellPrices) != 0 {
		t.Errorf("Expected sell order to be fully matched and removed")
	}
}

func TestBuyTooLowGoesToBook(t *testing.T) {
	ob := NewOrderBook()

	sell := Order{
		ID:        "sell-2",
		Timestamp: time.Now(),
		Price:     105.0,
		Quantity:  1.0,
		Type:      LimitOrder,
		Side:      Sell,
	}
	ob.AddOrder(sell)

	buy := Order{
		ID:        "buy-2",
		Timestamp: time.Now(),
		Price:     100.0, // Too low to match
		Quantity:  1.0,
		Type:      LimitOrder,
		Side:      Buy,
	}
	ob.AddOrder(buy)

	if len(ob.buyPrices) != 1 || len(ob.sellPrices) != 1 {
		t.Errorf("Expected buy and sell orders to remain on book")
	}
}

func TestPartialMatch(t *testing.T) {
	ob := NewOrderBook()

	sell := Order{
		ID:        "sell-3",
		Timestamp: time.Now(),
		Price:     100.0,
		Quantity:  2.0,
		Type:      LimitOrder,
		Side:      Sell,
	}
	ob.AddOrder(sell)

	buy := Order{
		ID:        "buy-3",
		Timestamp: time.Now(),
		Price:     105.0,
		Quantity:  1.0,
		Type:      LimitOrder,
		Side:      Buy,
	}
	ob.AddOrder(buy)

	if len(ob.sellPrices) != 1 {
		t.Fatalf("Sell should still be on book after partial match")
	}
	remaining := ob.sells[100.0].Front().Value.(Order).Quantity
	if remaining != 1.0 {
		t.Errorf("Expected 1.0 quantity remaining, got %v", remaining)
	}
}
