// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

type mapPerformance map[string]*structPerformance

func newMapPerformance() mapPerformance {
	return mapPerformance{}
}
func (p mapPerformance) Reset() {
	for k := range p {
		delete(p, k)
	}
}
func (p mapPerformance) Read(onElement func(string, *structPerformance)) {
	for k, v := range p {
		onElement(k, v)
	}
}
func (p mapPerformance) Add(name string, perf float64) {
	v, ok := p[name]
	if !ok {
		v = newStructPerformance()
		p[name] = v
	}
	v.AddValue(perf, 0)
}

func (p mapPerformance) ExclAdd(name string, perf float64, excl float64) {
	v, ok := p[name]
	if !ok {
		v = newStructPerformance()
		p[name] = v
	}
	v.AddComponent(perf, excl)
}
func (p mapPerformance) Merge(q mapPerformance) {
	for k, v := range q {
		if s, ok := p[k]; ok {
			s.Append(v)
		} else {
			p[k] = v
		}
	}
}
