package store

import (
	"sync"
	"wb_cource/internal/app/model"
)

type Cache struct {
	orders map[string]*model.Order
	mutex  sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]*model.Order),
	}
}

func (c *Cache) Set(order *model.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) Get(orderUID string) (*model.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	order, exists := c.orders[orderUID]
	return order, exists
}

func (c *Cache) GetAll() []*model.Order {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	orders := make([]*model.Order, 0, len(c.orders))
	for _, order := range c.orders {
		orders = append(orders, order)
	}
	return orders
}

func (c *Cache) LoadFromStore(store Store) error {
	orders, err := store.Order().GetAll()
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.orders = make(map[string]*model.Order)
	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}

	return nil
}
