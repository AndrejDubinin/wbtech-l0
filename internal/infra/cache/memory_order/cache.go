package memoryorder

import (
	"sync"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type Cache struct {
	mx   sync.RWMutex
	data map[string]*domain.Order
}

func New() *Cache {
	return &Cache{
		data: make(map[string]*domain.Order),
	}
}

func (c *Cache) Get(orderUID string) *domain.Order {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if order, inMap := c.data[orderUID]; inMap {
		return order
	}

	return nil
}

func (c *Cache) Add(order *domain.Order) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.data[order.OrderUID] = order
}
