package main

import (
	"context"
	"github.com/kainhuck/cache"
	"log"
	"time"
)

type IntCache struct {
	cache *cache.Cache
	cancel context.CancelFunc
}

func NewIntCache() *IntCache{
	ctx, cancel := context.WithCancel(context.Background())
	ca := cache.New(ctx, 2 * time.Second, 3 * time.Second)
	ca.RegisterBeforeDelete(func(key string, value interface{}) {
		v := value.(int)

		log.Printf("delete `%v`", v)
	})
	return &IntCache{
		cache: ca,
		cancel: cancel,
	}
}

func (c *IntCache)Shutdown() {
	if c.cancel != nil{
		defer c.cancel()
	}

	c.cache.For(func(key string, value interface{}) {
		v := value.(int)

		log.Printf("before stop `%v`", v)
	})
}

func (c *IntCache)Store(key string, value int){
	//c.cache.SetByDefaultExpiration(key, value) // 2 * time.Second
	c.cache.Set(key, value, 3 * time.Second)
}

func (c *IntCache)Load(key string) (int, bool){
	valueI, ok := c.cache.Get(key)
	if !ok {
		return 0, ok
	}
	value := valueI.(int)

	return value, true
}

func main() {
	intCache := NewIntCache()

	// store
	intCache.Store("bar", 18)
	intCache.Store("foo", 14)

	// load
	bar, ok := intCache.Load("bar")
	if ok {
		log.Println(bar)
	}

	time.Sleep(10 * time.Second)

	intCache.Shutdown()
}