package memoryorder

import (
	"container/list"
	"sync"

	"github.com/AndrejDubinin/wbtech-l0/internal/domain"
)

type LRUCache struct {
	capacity int64
	mx       sync.Mutex
	data     map[string]*list.Element
	list     *list.List
}

func New(capacity int64) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		data:     make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(orderUID string) *domain.Order {
	c.mx.Lock()
	defer c.mx.Unlock()

	if elem, inMap := c.data[orderUID]; inMap {
		c.list.MoveToFront(elem)
		return elem.Value.(*domain.Order)
	}

	return nil
}

func (c *LRUCache) Put(order *domain.Order) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if elem, inMap := c.data[order.OrderUID]; inMap {
		elem.Value = order
		c.list.MoveToFront(elem)
		return
	}

	elem := c.list.PushFront(order)
	c.data[order.OrderUID] = elem

	if int64(c.list.Len()) > c.capacity {
		last := c.list.Back()
		if last != nil {
			c.list.Remove(last)
			delete(c.data, last.Value.(*domain.Order).OrderUID)
		}
	}
}
