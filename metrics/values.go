package metrics

import (
	"sync"
	"sync/atomic"
)

type Values interface {
	Clear()
	Get(key string) interface{}
	Add(key string, value interface{})
	Snapshot() Values
}

func GetOrRegisterValues(name string, r Registry) Values {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewValues).(Values)
}

func NewValues() Values {
	if UseNilMetrics {
		return NilValues{}
	}
	return &StandardValues{}
}

func NewRegisteredKeyValue(name string, r Registry) Values {
	c := NewValues()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

type ValuesSnapshot map[string]interface{}

func (ValuesSnapshot) Clear() {
	panic("Clear called on a ValuesSnapshot")
}

func (c ValuesSnapshot) Get(key string) interface{} { return c[key] }

func (ValuesSnapshot) Add(key string, value interface{}) {
	panic("Add called on a ValuesSnapshot")
}

func (c ValuesSnapshot) Snapshot() Values { return c }

type NilValues map[string]interface{}

func (NilValues) Clear() {}

func (NilValues) Get(key string) interface{} { return nil }

func (NilValues) Add(key string, value interface{}) {}

func (NilValues) Snapshot() Values { return NilValues{} }

type StandardValues struct {
	values atomic.Value
	m      sync.RWMutex
}

func (c *StandardValues) Clear() {
	c.values.Store(make(map[string]interface{}))
}

func (c *StandardValues) Get(key string) interface{} {
	return c.values.Load().(map[string]interface{})[key]
}

func (c *StandardValues) Add(key string, value interface{}) {
	c.m.Lock()
	defer c.m.Unlock()
	m := c.values.Load().(map[string]interface{})
	m[key] = value
	c.values.Store(StandardValues{})
}

func (c *StandardValues) Snapshot() Values {
	return c.values.Load().(ValuesSnapshot)
}
