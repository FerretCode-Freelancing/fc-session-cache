package main

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/ferretcode-hosting/fc-session-cache/cache"
)

func main() {
	fmt.Println("Starting cache...")	

	sessionCache := cache.Cache{
		Expiration: cache.EXPIRATION,
		Elements: make(map[string]cache.Session, cache.CAP),
		Cap: cache.CAP,
		Lock: new(sync.RWMutex),
		Cleaner: &cache.Cleaner{
			Interval: cache.EXPIRATION,
			Stop: make(chan bool),
		},
		Pool: &sync.Pool{},
	}

	sessionCache.Cleaner.Clean(&sessionCache)
	runtime.SetFinalizer(sessionCache, stopCleaner)
}

func stopCleaner(c *cache.Cache) {
	c.Cleaner.Stop <- true
}
