package exp

import (
	"expvar"
	"fmt"
	"net/http"
	"sync"

	"github.com/dreamerminsk/gowiki/metrics"
)

type exp struct {
	expvarLock sync.Mutex
	registry   metrics.Registry
}

func (exp *exp) expHandler(w http.ResponseWriter, r *http.Request) {
	exp.syncToExpvar()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

func Exp(r metrics.Registry) {
	h := ExpHandler(r)
	http.Handle("/debug/metrics", h)
}

func ExpHandler(r metrics.Registry) http.Handler {
	e := exp{sync.Mutex{}, r}
	return http.HandlerFunc(e.expHandler)
}

func (exp *exp) getInt(name string) *expvar.Int {
	var v *expvar.Int
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p != nil {
		v = p.(*expvar.Int)
	} else {
		v = new(expvar.Int)
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
	return v
}

func (exp *exp) getFloat(name string) *expvar.Float {
	var v *expvar.Float
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p != nil {
		v = p.(*expvar.Float)
	} else {
		v = new(expvar.Float)
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
	return v
}

func (exp *exp) getString(name string) *expvar.String {
	var v *expvar.String
	exp.expvarLock.Lock()
	p := expvar.Get(name)
	if p != nil {
		v = p.(*expvar.String)
	} else {
		v = new(expvar.String)
		expvar.Publish(name, v)
	}
	exp.expvarLock.Unlock()
	return v
}

func (exp *exp) publishCounter(name string, metric metrics.Counter) {
	v := exp.getInt(name)
	v.Set(metric.Count())
}

func (exp *exp) publishValues(name string, metric metrics.Values) {
	s := metric.Snapshot()
	for _, k := range s.Keys() {
		exp.getString(name + "." + k).Set(fmt.Sprintf("%s", s.Get(k)))
	}
}

func (exp *exp) publishString(name string, metric metrics.String) {
	v := exp.getString(name)
	v.Set(metric.Value())
}

func (exp *exp) publishGauge(name string, metric metrics.Gauge) {
	v := exp.getInt(name)
	v.Set(metric.Value())
}

func (exp *exp) publishGaugeFloat64(name string, metric metrics.GaugeFloat64) {
	exp.getFloat(name).Set(metric.Value())
}

func (exp *exp) publishHistogram(name string, metric metrics.Histogram) {
	h := metric.Snapshot()
	//ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	exp.getInt(name + ".count").Set(h.Count())
	exp.getFloat(name + ".min").Set(float64(h.Min()))
	exp.getFloat(name + ".max").Set(float64(h.Max()))
	exp.getFloat(name + ".mean").Set(float64(h.Mean()))
	exp.getFloat(name + ".std-dev").Set(float64(h.StdDev()))
	//exp.getFloat(name + ".50-percentile").Set(float64(ps[0]))
	//exp.getFloat(name + ".75-percentile").Set(float64(ps[1]))
	//exp.getFloat(name + ".95-percentile").Set(float64(ps[2]))
	//exp.getFloat(name + ".99-percentile").Set(float64(ps[3]))
	//exp.getFloat(name + ".999-percentile").Set(float64(ps[4]))
}

func (exp *exp) publishMeter(name string, metric metrics.Meter) {
	m := metric.Snapshot()
	exp.getInt(name + ".count").Set(m.Count())
	exp.getFloat(name + ".one-minute").Set(float64(m.Rate1()))
	exp.getFloat(name + ".five-minute").Set(float64(m.Rate5()))
	exp.getFloat(name + ".fifteen-minute").Set(float64((m.Rate15())))
	exp.getFloat(name + ".mean").Set(float64(m.RateMean()))
}

func (exp *exp) publishTimer(name string, metric metrics.Timer) {
	t := metric.Snapshot()
	//ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
	exp.getInt(name + ".count").Set(t.Count())
	exp.getFloat(name + ".min").Set(float64(t.Min()))
	exp.getFloat(name + ".max").Set(float64(t.Max()))
	exp.getFloat(name + ".mean").Set(float64(t.Mean()))
	exp.getFloat(name + ".std-dev").Set(float64(t.StdDev()))
	//exp.getFloat(name + ".50-percentile").Set(float64(ps[0]))
	//exp.getFloat(name + ".75-percentile").Set(float64(ps[1]))
	//exp.getFloat(name + ".95-percentile").Set(float64(ps[2]))
	//exp.getFloat(name + ".99-percentile").Set(float64(ps[3]))
	//exp.getFloat(name + ".999-percentile").Set(float64(ps[4]))
	//exp.getFloat(name + ".one-minute").Set(float64(t.Rate1()))
	//exp.getFloat(name + ".five-minute").Set(float64(t.Rate5()))
	//exp.getFloat(name + ".fifteen-minute").Set(float64((t.Rate15())))
	//exp.getFloat(name + ".mean-rate").Set(float64(t.RateMean()))
}

func (exp *exp) syncToExpvar() {
	exp.registry.Each(func(name string, i interface{}) {
		switch i := i.(type) {
		case metrics.Counter:
			exp.publishCounter(name, i)
		case metrics.Values:
			exp.publishValues(name, i)
		case metrics.String:
			exp.publishString(name, i)
		case metrics.Gauge:
			exp.publishGauge(name, i)
		case metrics.GaugeFloat64:
			exp.publishGaugeFloat64(name, i)
		case metrics.Histogram:
			exp.publishHistogram(name, i)
		case metrics.Meter:
			exp.publishMeter(name, i)
		case metrics.Timer:
			exp.publishTimer(name, i)
		default:
			panic(fmt.Sprintf("unsupported type for '%s': %T", name, i))
		}
	})
}
