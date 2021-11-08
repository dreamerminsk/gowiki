package metrics

import "github.com/dreamerminsk/gowiki/atomic"

type String interface {
	Snapshot() String
	Update(string)
	Value() string
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
	return &StandardString{atomic.NewString("")}
}

func NewRegisteredGString(name string, r Registry) String {
	c := NewString()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

func NewFunctionalString(f func() string) String {
	if UseNilMetrics {
		return NilString{}
	}
	return &FunctionalString{value: f}
}

func NewRegisteredFunctionalString(name string, r Registry, f func() string) String {
	c := NewFunctionalString(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

type StringSnapshot string

func (g StringSnapshot) Snapshot() String { return g }

func (StringSnapshot) Update(string) {
	panic("Update called on a GaugeSnapshot")
}

func (g StringSnapshot) Value() string { return string(g) }

type NilString struct{}

func (NilString) Snapshot() String { return NilString{} }

func (NilString) Update(v string) {}

func (NilString) Value() string { return "" }

type StandardString struct {
	value atomic.String
}

func (g *StandardString) Snapshot() String {
	return StringSnapshot(g.Value())
}

func (g *StandardString) Update(v string) {
	g.value.Store(v)
}

func (g *StandardString) Value() string {
	return g.value.Load()
}

type FunctionalString struct {
	value func() string
}

func (g FunctionalString) Value() string {
	return g.value()
}

func (g FunctionalString) Snapshot() String { return StringSnapshot(g.Value()) }

func (FunctionalString) Update(string) {
	panic("Update called on a FunctionalString")
}
