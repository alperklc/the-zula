package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type Cache[T any] struct {
	store *ttlcache.Cache[string, T]
}

type CacheInterface[T any] interface {
	Reset(id string)
	Write(id string, obj T)
	Read(id string) *T
}

// Ensure Cache implements CacheInterface
var _ CacheInterface[any] = (*Cache[any])(nil)

func NewCache[T any](ttl time.Duration) (*Cache[T], error) {
	store := ttlcache.New[string, T](
		ttlcache.WithTTL[string, T](ttl),
	)

	go store.Start()

	return &Cache[T]{
		store: store,
	}, nil
}

func (c *Cache[T]) Reset(id string) {
	c.store.Delete(id)
}

func (c *Cache[T]) Write(id string, obj T) {
	c.store.Set(id, obj, ttlcache.DefaultTTL)
}

func (c *Cache[T]) Read(id string) *T {
	item := c.store.Get(id)
	if item == nil {
		return nil
	}

	response := item.Value()
	return &response
}
