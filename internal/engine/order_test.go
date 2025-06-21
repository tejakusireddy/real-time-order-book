package engine

import (
	"testing"
	"time"
)

func TestOrderValidation(t *testing.T) {
	tests := []struct {
		name    string
		order   Order
		wantErr bool
	}{
		{
			name: "valid limit buy order",
			order: Order{
				ID:        "123",
				Timestamp: time.Now(),
				Price:     100.0,
				Quantity:  1.5,
				Type:      LimitOrder,
				Side:      Buy,
			},
			wantErr: false,
		},
		{
			name: "invalid order type",
			order: Order{
				ID:        "124",
				Timestamp: time.Now(),
				Price:     100.0,
				Quantity:  1.5,
				Type:      "invalid",
				Side:      Buy,
			},
			wantErr: true,
		},
		{
			name: "zero quantity",
			order: Order{
				ID:        "125",
				Timestamp: time.Now(),
				Price:     100.0,
				Quantity:  0,
				Type:      LimitOrder,
				Side:      Sell,
			},
			wantErr: true,
		},
		{
			name: "market order no price (valid)",
			order: Order{
				ID:        "126",
				Timestamp: time.Now(),
				Quantity:  2.0,
				Type:      MarketOrder,
				Side:      Buy,
			},
			wantErr: false,
		},
		{
			name: "limit order no price (invalid)",
			order: Order{
				ID:        "127",
				Timestamp: time.Now(),
				Quantity:  2.0,
				Type:      LimitOrder,
				Side:      Sell,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
