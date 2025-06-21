package engine

import (
	"container/list"
	"sort"
	"sync"
	"time"
)

type OrderBook struct {
	mu    sync.Mutex
	buys  map[float64]*list.List // price → orders (FIFO)
	sells map[float64]*list.List

	buyPrices  []float64
	sellPrices []float64
	notifier   Notifiee
}

func NewOrderBook(n Notifiee) *OrderBook {
	return &OrderBook{
		buys:       make(map[float64]*list.List),
		sells:      make(map[float64]*list.List),
		notifier:   n,
		buyPrices:  []float64{},
		sellPrices: []float64{},
	}
}

func (ob *OrderBook) AddOrder(order Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if err := order.Validate(); err != nil {
		return // ignore invalid orders silently (for now)
	}

	switch order.Side {
	case Buy:
		ob.matchBuy(order)
	case Sell:
		ob.matchSell(order)
	}
	if ob.notifier != nil && order.Quantity > 0 && order.Type == LimitOrder {
		ob.notifier.OnAdd(order)
	}
}

func (ob *OrderBook) matchBuy(order Order) {
	for len(ob.sellPrices) > 0 && order.Quantity > 0 {
		bestPrice := ob.sellPrices[0]
		if order.Type == LimitOrder && order.Price < bestPrice {
			break
		}

		queue := ob.sells[bestPrice]
		for queue.Len() > 0 && order.Quantity > 0 {
			head := queue.Front()
			existing := head.Value.(Order)

			if order.Quantity >= existing.Quantity {
				matchedQty := existing.Quantity
				order.Quantity -= matchedQty
				queue.Remove(head)

				if ob.notifier != nil {
					ob.notifier.OnMatch(MatchEvent{
						BuyOrderID:   order.ID,
						SellOrderID:  existing.ID,
						Price:        bestPrice,
						Quantity:     matchedQty,
						TimestampUTC: time.Now().UTC().Format(time.RFC3339),
					})
				}
			} else {
				matchedQty := order.Quantity
				existing.Quantity -= matchedQty
				order.Quantity = 0
				head.Value = existing

				if ob.notifier != nil {
					ob.notifier.OnMatch(MatchEvent{
						BuyOrderID:   order.ID,
						SellOrderID:  existing.ID,
						Price:        bestPrice,
						Quantity:     matchedQty,
						TimestampUTC: time.Now().UTC().Format(time.RFC3339),
					})
				}
			}
		}

		if queue.Len() == 0 {
			delete(ob.sells, bestPrice)
			ob.sellPrices = ob.sellPrices[1:]
		}
	}

	if order.Quantity > 0 && order.Type == LimitOrder {
		ob.enqueue(ob.buys, &ob.buyPrices, order, true)
	}
}

func (ob *OrderBook) matchSell(order Order) {
	for len(ob.buyPrices) > 0 && order.Quantity > 0 {
		bestPrice := ob.buyPrices[0]
		if order.Type == LimitOrder && order.Price > bestPrice {
			break
		}

		queue := ob.buys[bestPrice]
		for queue.Len() > 0 && order.Quantity > 0 {
			head := queue.Front()
			existing := head.Value.(Order)

			if order.Quantity >= existing.Quantity {
				matchedQty := existing.Quantity
				order.Quantity -= matchedQty
				queue.Remove(head)

				if ob.notifier != nil {
					ob.notifier.OnMatch(MatchEvent{
						BuyOrderID:   existing.ID,
						SellOrderID:  order.ID,
						Price:        bestPrice,
						Quantity:     matchedQty,
						TimestampUTC: time.Now().UTC().Format(time.RFC3339),
					})
				}
			} else {
				matchedQty := order.Quantity
				existing.Quantity -= matchedQty
				order.Quantity = 0
				head.Value = existing

				if ob.notifier != nil {
					ob.notifier.OnMatch(MatchEvent{
						BuyOrderID:   existing.ID,
						SellOrderID:  order.ID,
						Price:        bestPrice,
						Quantity:     matchedQty,
						TimestampUTC: time.Now().UTC().Format(time.RFC3339),
					})
				}
			}
		}

		if queue.Len() == 0 {
			delete(ob.buys, bestPrice)
			ob.buyPrices = ob.buyPrices[1:]
		}
	}

	if order.Quantity > 0 && order.Type == LimitOrder {
		ob.enqueue(ob.sells, &ob.sellPrices, order, false)
	}
}

func (ob *OrderBook) enqueue(book map[float64]*list.List, prices *[]float64, order Order, isBuy bool) {
	if _, ok := book[order.Price]; !ok {
		book[order.Price] = list.New()
		*prices = append(*prices, order.Price)

		if isBuy {
			sort.Sort(sort.Reverse(sort.Float64Slice(*prices)))
		} else {
			sort.Float64s(*prices)
		}
	}
	book[order.Price].PushBack(order)
}
