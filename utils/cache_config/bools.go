// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package cache_config

//boolean cache
type Bools struct {
	item_count int
	current    int
	arrays     [4][]byte
}

func (s *Bools) Init(item_count int) *Bools {
	s.item_count = item_count
	s.current = 3
	for i := 0; i < 4; i++ {
		a := make([]byte, item_count)
		for j := 0; j < item_count; j++ {
			a[j] = 0
		}
		s.arrays[i] = a
	}
	return s
}
func (s *Bools) Find(id int) (bool, bool) {
	if !in_range(id, s.item_count) {
		return false, false
	}
	r := s.arrays[s.current][id]
	if (r & 2) == 0 {
		return false, false
	}
	return (r & 1) == 1, true
}

func (s *Bools) Read(id int, default_value bool) bool {
	if v, found := s.Find(id); found {
		return v
	}
	return default_value
}
func (s *Bools) Update(id int, value bool) bool {
	if !in_range(id, s.item_count) {
		return false
	}
	var v byte = 2
	if value {
		v = 3
	}
	s.arrays[(s.current+1)%4][id] = v
	return true
}
func (s *Bools) Commit() { s.current = (s.current + 1) % 4 }
