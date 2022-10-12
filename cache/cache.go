package cache

import (
	"fmt"
	"sync"
	"time"
)

const (
	EXPIRATION = 24 * 7 * time.Hour
	CAP = 1024
)

type Cache struct {
	Expiration time.Duration
	Elements map[string]Session
	Cap int64
	Size int64
	Lock *sync.RWMutex
	Pool *sync.Pool
	Cleaner *Cleaner
}

type Session struct {
	C string // session cookie
	S interface{} //session object
	Expiration int64
	LastAccess int64
}

type Cleaner struct {
	Interval time.Duration
	Stop chan bool
}

func (c *Cache) Put(co string, s interface{}) error {
	expiration := time.Now().Add(EXPIRATION).UnixNano()
	lastAccess := time.Now().UnixNano()

	if c.Size + 1 > c.Cap {
		if _, err := c.Remove(c.LRU().C); err != nil {
			return err
		}
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()

	if session, ok := c.Elements[co]; ok {
		session.S = s
		session.Expiration = expiration
		session.LastAccess = lastAccess

		return nil
	}

	session := Session{
		S: s,
		Expiration: expiration,
		LastAccess: lastAccess,
	}

	c.Pool.Put(session)
	c.Elements[co] = session
	c.Size += 1

	return nil
}

func (c *Cache) Get(co string) (s interface{}, err error) {
	session := c.Pool.Get()

	if sess, ok := session.(Session); ok {
		if sess.C == co {
			return sess.S, nil
		}
	}

	expiration := time.Now().Add(EXPIRATION).UnixNano()
	lastAccess := time.Now().UnixNano()

	c.Lock.RLock()
	defer c.Lock.RUnlock()

	if sess, ok := c.Elements[co]; ok {
		sess.Expiration = expiration
		sess.LastAccess = lastAccess

		return sess.S, nil
	}

	return nil, nil
}

func (c *Cache) Remove(co string) (isFound bool, err error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	s := c.Pool.Get()

	if s != nil && s.(Session).C != co {
		c.Pool.Put(s)
	}

	for sess := range c.Elements {
		if sess == co {
			delete(c.Elements, sess)
			return true, nil
		}
	}

	return false, nil
}

func (c *Cache) Flush() error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	for k := range c.Elements {
		delete(c.Elements, k)
	}

	c.Pool = &sync.Pool{}

	return nil
}

func (c *Cache) CleanExpired() {
	now := time.Now().UnixNano()
	
	c.Lock.Lock()
	defer c.Lock.Unlock()

	for co, s := range c.Elements {
		if s.Expiration > 0 && now > s.Expiration {
			if _, err := c.Remove(co); err != nil {
				fmt.Println(fmt.Sprintf("Error cleaning cached session: %s", err))
			}	
		}
	}
}

func (cl *Cleaner) Clean(c *Cache) {
	ticker := time.NewTicker(cl.Interval)

	for {
		select {
		case <- ticker.C:
			c.CleanExpired()
		case <- cl.Stop:
			ticker.Stop()
			return
		}
	}
}
