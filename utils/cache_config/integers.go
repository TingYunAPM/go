// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package cache_config

type int_item struct {
	used  bool
	value int64
}

func (i *int_item) Init() { i.used = false }
func (i *int_item) Set(value int64) {
	i.used = true
	i.value = value
}

//integer cache
type Integers struct {
	item_count int
	current    int
	arrays     [4][]int_item
}

func (s *Integers) Init(item_count int) *Integers {
	s.item_count = item_count
	s.current = 3
	for i := 0; i < 4; i++ {
		a := make([]int_item, item_count)
		for j := 0; j < item_count; j++ {
			a[j].Init()
		}
		s.arrays[i] = a
	}
	return s
}
func (s *Integers) Find(id int) (int64, bool) {
	if !in_range(id, s.item_count) {
		return -1, false
	}
	item := &s.arrays[s.current][id]
	if item.used {
		return item.value, true
	}
	return -1, false
}

func (s *Integers) Read(id int, default_value int64) int64 {
	if v, found := s.Find(id); found {
		return v
	}
	return default_value
}
func (s *Integers) Update(id int, value int64) bool {
	if !in_range(id, s.item_count) {
		return false
	}
	s.arrays[(s.current+1)%4][id].Set(value)
	return true
}
func (s *Integers) Commit() { s.current = (s.current + 1) % 4 }
