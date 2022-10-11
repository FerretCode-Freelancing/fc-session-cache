package main

import (
	"sync"
	"time"
)

type FcCache struct {
	expiration time.Duration
	elements   map[string]string
	cap        int64
	size       int64
	lock       *sync.RWMutex
	pool       *sync.Pool
	cleaner    *Cleaner
}

type Cleaner struct {
	interval time.Duration
	stop     chan bool
}

func (c *FcCache) Get(k string) (v string, err error) {
	return "", nil
}
