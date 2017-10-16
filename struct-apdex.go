// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

type structApdex struct {
	count_satisfying  int32
	count_tolerating  int32
	count_frustrating int32
	apdex_t_in_millis int32
}

func (p *structApdex) IntSlice() []int32 {
	r := make([]int32, 4)
	r[0] = p.count_satisfying
	r[1] = p.count_tolerating
	r[2] = p.count_frustrating
	r[3] = p.apdex_t_in_millis
	return r
}
func (a *structApdex) Reset() {
	a.count_satisfying = 0
	a.count_tolerating = 0
	a.count_frustrating = 0
	a.apdex_t_in_millis = 0
}
func (a *structApdex) Init(apdex_t int32) *structApdex {
	a.apdex_t_in_millis = apdex_t
	return a
}
func (a *structApdex) Add(perf_data int32, has_error bool) {
	if a.apdex_t_in_millis > 0 {
		if has_error {
			a.count_frustrating++
			return
		}
		if perf_data < a.apdex_t_in_millis {
			a.count_satisfying++
		} else if perf_data < a.apdex_t_in_millis*4 {
			a.count_tolerating++
		} else {
			a.count_frustrating++
		}
	}
}
func (a *structApdex) Merge(b *structApdex) {
	if a.apdex_t_in_millis <= 0 {
		a.apdex_t_in_millis = b.apdex_t_in_millis
	}
	if a.apdex_t_in_millis == b.apdex_t_in_millis {
		a.count_frustrating += b.count_frustrating
		a.count_satisfying += b.count_satisfying
		a.count_tolerating += b.count_tolerating
	}
}
