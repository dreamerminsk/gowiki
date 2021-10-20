package metrics

import "sync/atomic"

type KeyValue interface {
	Clear()
           Value() interface{}
	Key() string
	SetValue(interface{})
	Snapshot() KeyValue
}

func GetOrRegisterKeyValue(name string, r Registry) KeyValue {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewKeyValue).(KeyValue)
}

func NewKeyValue() KeyValue {
	if UseNilMetrics {
		return NilKeyValue{}
	}
	return &StandardKeyValue{0}
}

func NewRegisteredKeyValue(name string, r Registry) KeyValue {
	c := NewKeyValue()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

type KeyValueSnapshot {
Key string
Value interface{}
}

func (KeyValueSnapshot) Clear() {
	panic("Clear called on a CounterSnapshot")
}

func (c KeyValueSnapshot) Count() int64 { return int64(c) }

func (KeyValueSnapshot) Dec(int64) {
	panic("Dec called on a CounterSnapshot")
}

func (KeyValueSnapshot) Inc(int64) {
	panic("Inc called on a CounterSnapshot")
}

func (c KeyValueSnapshot) Snapshot() KeyValue { return c }

type NilKeyValue struct{}

func (NilKeyValue) Clear() {}

func (NilKeyValue) Count() int64 { return 0 }

func (NilKeyValue) Dec(i int64) {}

func (NilKeyValue) Inc(i int64) {}

func (NilKeyValue) Snapshot() KeyValue { return NilKeyValue{} }

type StandardKeyValue struct {
	key string
value   interface{}
}

func (c *StandardKeyValue) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

func (c *StandardKeyValue) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

func (c *StandardKeyValue) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

func (c *StandardKeyValue) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}

func (c *StandardKeyValue) Snapshot() KeyValue {
	return KeyValueSnapshot(c.Count())
}
