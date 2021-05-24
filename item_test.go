package cache

import (
	"log"
	"testing"
	"time"
)

func TestItem(t *testing.T) {
	i := &item{
		value:     10,
		expiredAt: time.Now(),
	}

	log.Println(i.expiredAt)

	i.addTime(2 * time.Second)
	log.Println(i.expiredAt)
}
