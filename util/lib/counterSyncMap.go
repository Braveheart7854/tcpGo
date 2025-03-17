/**
 * @author  tongh
 * @date  2025/3/17 18:52
 */
package lib

import (
	"sync"
	"sync/atomic"
)

type CounterSyncMap struct {
	m     sync.Map
	count atomic.Int64
}

func (c *CounterSyncMap) Store(key, value interface{}) {
	c.m.Store(key, value)
	c.count.Add(1)
}

func (c *CounterSyncMap) Delete(key interface{}) {
	_, load := c.m.LoadAndDelete(key)
	if load {
		c.count.Add(-1)
	}
}

func (c *CounterSyncMap) Load(key interface{}) (value interface{}, ok bool) {
	return c.m.Load(key)
}

func (c *CounterSyncMap) Len() int64 {
	return c.count.Load()
}
