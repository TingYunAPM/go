// Copyright 2018 冯立强 fenglq@tingyun.com.  All rights reserved.

package cache_config

type array_item struct {
	used  bool
	value []int64
}

func (a *array_item) Init() { a.used = false }
func (a *array_item) Set(value []int64) {
	a.used = true
	a.value = make([]int64, len(value))
	for i, v := range value {
		a.value[i] = v
	}
}

//Array cache
type Arrays struct {
	item_count int
	current    int
	arrays     [4][]array_item
	cleared    bool
}

func (s *Arrays) cleanNext() {
	if !s.cleared {
		for j := 0; j < s.item_count; j++ {
			s.arrays[(s.current+1)%4][j].Init()
		}
		s.cleared = true
	}
}
func (s *Arrays) Init(item_count int) *Arrays {
	s.item_count = item_count
	s.current = 3
	s.cleared = true
	for i := 0; i < 4; i++ {
		a := make([]array_item, item_count)
		for j := 0; j < item_count; j++ {
			a[j].Init()
		}
		s.arrays[i] = a
	}
	return s
}
func (s *Arrays) Find(id int) ([]int64, bool) {
	if !in_range(id, s.item_count) {
		return nil, false
	}
	item := &s.arrays[s.current][id]
	if item.used {
		return item.value, true
	}
	return nil, false
}

func (s *Arrays) Read(id int, default_value []int64) []int64 {
	if v, found := s.Find(id); found {
		return v
	}
	return default_value
}
func (s *Arrays) Update(id int, value []int64) bool {
	if !in_range(id, s.item_count) {
		return false
	}
	s.cleanNext()
	s.arrays[(s.current+1)%4][id].Set(value)
	return true
}
func (s *Arrays) Commit() {
	s.cleanNext()
	s.current = (s.current + 1) % 4
	s.cleared = false
}
