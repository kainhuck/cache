package cache

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := New(2*time.Second, 3*time.Second)
	cache.SetByDefaultExpiration("name", "Kainhuck")
	//cache.Set("age", "Kainhuck", 5 *time.Second)
	for i := 0; i < 10; i++ {
		name, ok := cache.Get("name")
		if ok {
			log.Println(name)
			cache.Grow("name", 2*time.Second)
		} else {
			log.Println("not ok")
		}
		time.Sleep(1 * time.Second)
	}
}

func TestTime(t *testing.T) {
	a := time.Now()
	fmt.Println(a)
	b := a.Add(2 * time.Second)
	fmt.Println(b)
}

func TestGrow(t *testing.T) {
	cache := New(2*time.Second, 3*time.Second)
	cache.Set("age", 19, 2*time.Second)
	for i := 0; i < 10; i++ {
		cache.Grow("age", 2*time.Second)
		time.Sleep(1 * time.Second)
	}
}
