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
		order := elem.Value.(domain.Order)
		return &order
	}

	return nil
}

func (c *LRUCache) Put(order *domain.Order) {
	c.mx.Lock()
	defer c.mx.Unlock()

	orderCopy := *order

	if elem, inMap := c.data[orderCopy.OrderUID]; inMap {
		elem.Value = orderCopy
		c.list.MoveToFront(elem)
		return
	}

	elem := c.list.PushFront(orderCopy)
	c.data[orderCopy.OrderUID] = elem

	if int64(c.list.Len()) > c.capacity {
		last := c.list.Back()
		if last != nil {
			c.list.Remove(last)
			delete(c.data, last.Value.(domain.Order).OrderUID)
		}
	}
}
