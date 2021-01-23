package cache

import (
	"log"
	"time"

	"github.com/allegro/bigcache"
)

type store struct {
	cache *bigcache.BigCache
}

func (s *store) Get(key string) ([]byte, error) {
	res, err := s.cache.Get(key)
	if err != nil {
		log.Printf("got an error from bigcache: %v", err)
		return nil, err
	}

	return res, nil
}

func (s *store) Set(key string, bb []byte) error {
	err := s.cache.Set(key, bb)
	if err != nil {
		log.Printf("got an error from bigcache: %v", err)
		return err
	}

	return nil
}

func New(eviction int64, shards int) (*store, error) {
	config := bigcache.Config{
		Shards:             shards,
		LifeWindow:         time.Duration(eviction) * time.Second,
		CleanWindow:        0,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		Verbose:            true,
		HardMaxCacheSize:   0,
	}
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	return &store{cache: cache}, nil
}
