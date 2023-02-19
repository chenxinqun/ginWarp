package configs

import (
	"sync"
	"time"

	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
)

var _ MapHandel = (*mapHandel)(nil)

type MapHandel interface {
	Kvs() map[string]interface{}
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	GetString(key string) (s string)
	GetBool(key string) (b bool)
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetFloat64(key string) (f64 float64)
	GetTime(key string) (t time.Time)
	GetDuration(key string) (d time.Duration)
	GetIntSlice(key string) (si []int)
	GetStringSlice(key string) (ss []string)
	GetStringMap(key string) (sm map[string]interface{})
	GetStringMapString(key string) (sms map[string]string)
	GetStringMapStringSlice(key string) (smss map[string][]string)
}

type mapHandel struct {
	mu  sync.RWMutex
	kvs map[string]interface{}
}

func NewMapHandel(data map[string]interface{}) MapHandel {
	return &mapHandel{kvs: data}
}

func (c *mapHandel) Kvs() map[string]interface{} {
	return c.kvs
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.kvs if it was not used previously.
func (c *mapHandel) Set(key string, value interface{}) {
	c.mu.Lock()
	c.kvs[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *mapHandel) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.kvs[key]
	c.mu.RUnlock()
	return
}

// GetString returns the value associated with the key as a string.
func (c *mapHandel) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *mapHandel) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *mapHandel) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *mapHandel) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *mapHandel) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *mapHandel) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *mapHandel) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *mapHandel) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *mapHandel) GetIntSlice(key string) (ss []int) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]int)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *mapHandel) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, ok = val.(etcdx.ConfigStringMap)
		if !ok {
			sm, _ = val.(map[string]interface{})
		}
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *mapHandel) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *mapHandel) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
