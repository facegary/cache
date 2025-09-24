package cache

import (
	"time"
)

type iObject interface {
	Value() any
	IsExpired() bool
}

type object struct {
	value     any
	expiredAt *time.Time
}

func newObject(value any, ttl time.Duration) iObject {
	object := &object{value: value}
	if ttl > 0 {
		expiredAt := time.Now().Add(ttl)
		object.expiredAt = &expiredAt
	}
	return object
}

func (o *object) Value() any {
	return o.value
}

func (o *object) IsExpired() bool {
	if o.expiredAt != nil {
		return o.expiredAt.Before(time.Now())
	}
	return false
}
