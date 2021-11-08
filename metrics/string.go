package metrics

import "sync/atomic"

type String interface {
	Snapshot() String
	Update(int64)
	Value() int64
}

func GetOrRegisterString(name string, r Registry) Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGauge).(Gauge)
}

func NewString() String {
	if UseNilMetrics {
		return NilString{}
	}
	return &StandardString{0}
}

func NewRegisteredGString(name string, r Registry) String {
	c := NewString()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

func NewFunctionalString(f func() int64) String {
	if UseNilMetrics {
		return NilString{}
	}
	return &FunctionalString{value: f}
}

func NewRegisteredFunctionalString(name string, r Registry, f func() int64) String {
	c := NewFunctionalString(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

type StringSnapshot int64

func (g StringSnapshot) Snapshot() String { return g }

func (StringSnapshot) Update(int64) {
	panic("Update called on a GaugeSnapshot")
}

func (g StringSnapshot) Value() int64 { return int64(g) }

type NilString struct{}

func (NilString) Snapshot() String { return NilString{} }

func (NilString) Update(v int64) {}

func (NilString) Value() int64 { return 0 }

type StandardString struct {
	value int64
}

func (g *StandardString) Snapshot() String {
	return StringSnapshot(g.Value())
}

func (g *StandardString) Update(v int64) {
	atomic.StoreInt64(&g.value, v)
}

func (g *StandardString) Value() int64 {
	return atomic.LoadInt64(&g.value)
}

type FunctionalString struct {
	value func() int64
}

func (g FunctionalString) Value() int64 {
	return g.value()
}

func (g FunctionalString) Snapshot() String { return StringSnapshot(g.Value()) }

func (FunctionalString) Update(int64) {
	panic("Update called on a FunctionalString")
}
