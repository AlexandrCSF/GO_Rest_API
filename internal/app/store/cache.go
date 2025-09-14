package store

import (
	"log"
	"sync"
	"time"
	"wb_cource/internal/app/config"
	"wb_cource/internal/app/model"
)

type CacheElement struct {
	orderUID  string
	response  model.Order
	ttl       time.Duration
	createdAt time.Time
}

type Cache struct {
	orders     map[string]*model.Order
	cacheElems []*CacheElement
	mutex      sync.RWMutex
	config     config.Config
	stopChan   chan struct{}
}

func NewCache(cfg config.Config) *Cache {
	cache := &Cache{
		orders:     make(map[string]*model.Order),
		cacheElems: make([]*CacheElement, 0),
		config:     cfg,
		stopChan:   make(chan struct{}),
	}

	go cache.startCleanup()
	return cache
}

func (c *Cache) Set(order *model.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.orders[order.OrderUID] = order
	c.cacheElems = append(c.cacheElems, &CacheElement{
		orderUID:  order.OrderUID,
		response:  *order,
		ttl:       time.Duration(c.config.CacheTTL) * time.Second,
		createdAt: time.Now(),
	})
}

func (c *Cache) Get(orderUID string) (*model.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	order, exists := c.orders[orderUID]
	if exists {
		log.Printf("Значение для orderUID %s взято из кеша", orderUID)
	}
	return order, exists
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopChan:
			return
		}
	}
}

func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	var validElems []*CacheElement

	for _, elem := range c.cacheElems {
		if now.Sub(elem.createdAt) < elem.ttl {
			validElems = append(validElems, elem)
		} else {
			delete(c.orders, elem.orderUID)
		}
	}

	c.cacheElems = validElems
	log.Printf("Очистка кеша: удалено %d устаревших элементов", len(c.cacheElems)-len(validElems))
}

func (c *Cache) Stop() {
	close(c.stopChan)
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
	c.cacheElems = make([]*CacheElement, 0)

	now := time.Now()
	for _, order := range orders {
		c.orders[order.OrderUID] = order
		c.cacheElems = append(c.cacheElems, &CacheElement{
			orderUID:  order.OrderUID,
			response:  *order,
			ttl:       time.Duration(c.config.CacheTTL) * time.Second,
			createdAt: now,
		})
	}

	return nil
}
