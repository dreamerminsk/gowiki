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
key string
value interface{}
}

func (KeyValueSnapshot) Clear() {
	panic("Clear called on a CounterSnapshot")
}

func (c KeyValueSnapshot) Key() string { return c.key }

func (c KeyValueSnapshot) Value() string { return c.value }

func (KeyValueSnapshot) SetValue(interface{}) {
	panic("SetValue called on a KeyValueSnapshot")
}



func (c KeyValueSnapshot) Snapshot() KeyValue { return c }

type NilKeyValue struct{}

func (NilKeyValue) Clear() {}

func (NilKeyValue) Key() string { return "" }

func (NilKeyValue) Value()  interface {   return nil;}

func (NilKeyValue) SetValue(i interface{}) {}


func (NilKeyValue) Snapshot() KeyValue { return NilKeyValue{} }

type StandardKeyValue struct {
	key string
value   interface{}
}

func (c *StandardKeyValue) Clear() {
	atomic.Store(c, struct{})
}

func (c *StandardKeyValue) Key() string {
	return atomic.Load(c).key
}



func (c *StandardKeyValue) Value() interface{} {
	return atomic.Load(c).value
}

func (c *StandardKeyValue) SetValue(i interface{}) {
	atomic.Store(c, {})
}


func (c *StandardKeyValue) Snapshot() KeyValue {
	return KeyValueSnapshot(c.Key(), c.Value())
}
