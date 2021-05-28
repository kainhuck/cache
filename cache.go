package cache

import (
	"context"
	"log"
	"sync"
	"time"
)

type CallBackFunc func(key string, value interface{})

type Cache struct {
	container         map[string]*item // 用于存放数据的容器
	mutex             sync.Mutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	ctx               context.Context

	beforeDelete []CallBackFunc
}

func New(ctx context.Context, defaultExpiration, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		container:         make(map[string]*item),
		mutex:             sync.Mutex{},
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		ctx:               ctx,

		beforeDelete: make([]CallBackFunc, 0),
	}
	go c.run()

	return c
}

// 注册回调
func (c *Cache) RegisterBeforeDelete(f ...CallBackFunc) {
	c.beforeDelete = append(c.beforeDelete, f...)
}

// 存值并指定过期时间
func (c *Cache) Set(key string, value interface{}, d time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// 将数据组合称item存进容器中
	c.container[key] = &item{
		value:     value,
		expiredAt: time.Now().Add(d),
	}
	log.Println("SET SUCCESS", key)
}

// 使用默认的过期时间
func (c *Cache) SetByDefaultExpiration(key string, value interface{}) {
	c.Set(key, value, c.defaultExpiration)
}

// 从cache中获取数据, 如果数据过期或者不存在就返回false
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.container[key]
	if !ok || item.isExpired() {
		return nil, false
	}

	return item.value, true
}

// 将指定的项延长寿命
func (c *Cache) Grow(key string, d time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ok := c.container[key]
	if ok {
		item.addTime(d)
	}
}

// for循环
func (c *Cache) For(f CallBackFunc)  {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.container{
		f(k, v.value)
	}
}

// 清除过期的item
func (c *Cache) cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, item := range c.container {
		if item.isExpired() {
			if len(c.beforeDelete) > 0 {
				for _, f := range c.beforeDelete{
					f(key, item.value)
				}
			}
			delete(c.container, key)
		}
	}
}

func (c *Cache) run() {
	ticker := time.NewTicker(c.cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <- c.ctx.Done():
			return
		}
	}
}
