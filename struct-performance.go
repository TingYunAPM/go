// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

type structPerformance struct {
	sum          float64 //累加和
	exclusive    float64
	valueMax     float64 //最大值
	valueMin     float64 //最小值
	sumSquare    float64 //平方和
	access_count int32   //累加次数
}
type structSqlTracePerf struct {
	Count int32
	Sum   float64
	Max   float64
	Min   float64
}

func (s *structSqlTracePerf) Reset() {
	s.Count = 0
	s.Sum = 0
	s.Max = 0
	s.Min = 0
}
func (s *structSqlTracePerf) Append(t *structSqlTracePerf) {
	s.Sum += t.Sum
	if t.Count > 0 {
		if s.Count == 0 || t.Max > s.Max {
			s.Max = t.Max
		}
		if s.Count == 0 || t.Min < s.Min {
			s.Min = t.Min
		}
	}
	s.Count += t.Count
}
func (s *structSqlTracePerf) Add(value float64) {
	if s.Count == 0 {
		s.Max = value
		s.Min = value
	} else {
		if value > s.Max {
			s.Max = value
		} else if value < s.Min {
			s.Min = value
		}
	}
	s.Count++
	s.Sum += value
}
func newStructSqlTracePerf() *structSqlTracePerf {
	return &structSqlTracePerf{0, 0, 0, 0}
}
func (p *structPerformance) IntSlice() []int64 {
	r := make([]int64, 6)
	r[0] = int64(p.access_count)
	r[1] = int64(p.sum)
	r[2] = int64(p.exclusive)
	r[3] = int64(p.valueMax)
	r[4] = int64(p.valueMin)
	r[5] = int64(p.sumSquare)
	return r
}
func (p *structPerformance) FloatSlice() []interface{} {
	r := make([]interface{}, 6)
	r[0] = p.access_count
	r[1] = p.sum
	r[2] = p.exclusive
	r[3] = p.valueMax
	r[4] = p.valueMin
	r[5] = p.sumSquare
	return r
}

func newStructPerformance() *structPerformance {
	r := &structPerformance{}
	r.Reset()
	return r
}
func (p *structPerformance) Reset() *structPerformance {
	p.sum = 0
	p.exclusive = 0
	p.valueMax = 0
	p.valueMin = 0
	p.sumSquare = 0
	p.access_count = 0
	return p
}

func (p *structPerformance) Append(q *structPerformance) {
	if q.access_count > 0 {
		p.sum += q.sum
		if p.access_count == 0 || q.valueMax > p.valueMax {
			p.valueMax = q.valueMax
		}
		if p.access_count == 0 || q.valueMin < p.valueMin {
			p.valueMin = q.valueMin
		}
		p.sumSquare += q.sumSquare
		p.exclusive += q.exclusive
		p.access_count += q.access_count
	}
}
func (p *structPerformance) AddComponent(value float64, excl float64) {
	p.sum += value
	if p.access_count == 0 {
		p.valueMax = excl
		p.valueMin = excl
	} else {
		if p.valueMax < excl {
			p.valueMax = excl
		}
		if p.valueMin > excl {
			p.valueMin = excl
		}
	}
	p.access_count++
	p.sumSquare += excl * excl
	p.exclusive += excl
}
func (p *structPerformance) AddValue(value float64, excl float64) {
	p.sum += value
	if p.access_count == 0 {
		p.valueMax = value
		p.valueMin = value
	} else {
		if p.valueMax < value {
			p.valueMax = value
		}
		if p.valueMin > value {
			p.valueMin = value
		}
	}
	p.access_count++
	p.sumSquare += value * value
	p.exclusive += excl
}

//!used
func (p *structPerformance) AppendValue(value float64, count int32, excl float64) {
	if count > 0 {
		avg := value / float64(count)
		if p.access_count == 0 {
			p.valueMax = avg
			p.valueMin = avg
		} else {
			if p.valueMax < avg {
				p.valueMax = avg
			}
			if p.valueMin > avg {
				p.valueMin = avg
			}
		}
		p.access_count += count
		p.sum += value
		p.sumSquare += avg * value
		p.exclusive += excl
	}
}
