package websocket

import (
	"encoding/json"
	"log"

	"github.com/tejakusireddy/real-time-order-book/internal/engine"
)

func (h *Hub) OnMatch(event engine.MatchEvent) {
	msg, err := json.Marshal(struct {
		Type  string            `json:"type"`
		Event engine.MatchEvent `json:"event"`
	}{
		Type:  "match",
		Event: event,
	})
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	h.broadcast <- msg
}

func (h *Hub) OnAdd(order engine.Order) {
	msg, err := json.Marshal(struct {
		Type  string       `json:"type"`
		Order engine.Order `json:"order"`
	}{
		Type:  "order_added",
		Order: order,
	})
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	h.broadcast <- msg
}
