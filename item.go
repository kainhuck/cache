package cache

import "time"

// 这个类型是cache存储的基本类型
type item struct {
	// 值
	value interface{}
	// 过期时间
	expiredAt time.Time
}

func (i *item) isExpired() bool {
	return time.Now().After(i.expiredAt)
}

func (i *item) addTime(d time.Duration) {
	i.expiredAt = i.expiredAt.Add(d)
}
