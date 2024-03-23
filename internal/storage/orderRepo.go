package storage

import (
	"errors"

	"github.com/federlizer/alokin-zota-integration/internal"
)

// OrderRepo is a simple in-memory storage for created Orders
type OrderRepo struct {
	// Orders holds all orders that have been created.
	// The key is the ID of the order and the value is the order itself.
	orders map[string]*internal.Order
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{
		orders: make(map[string]*internal.Order),
	}
}

func (r *OrderRepo) AddOrder(order *internal.Order) error {
	_, exists := r.orders[order.Id.String()]
	if exists {
		return errors.New("Another order with the same ID already exists")
	}

	r.orders[order.Id.String()] = order
	return nil
}

func (r *OrderRepo) GetOrder(id string) *internal.Order {
	order, exists := r.orders[id]
	if !exists {
		return nil
	}

	return order
}

func (r *OrderRepo) GetAll() []*internal.Order {
	orderArray := make([]*internal.Order, 0)

	for _, order := range r.orders {
		orderArray = append(orderArray, order)
	}

	return orderArray
}
