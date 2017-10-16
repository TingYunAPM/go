// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package cache_config

type string_item struct {
	used  bool
	value string
}

func (i *string_item) Init() { i.used = false }
func (i *string_item) Set(value string) {
	i.used = true
	i.value = value
}

//string cache
type Strings struct {
	item_count int
	current    int
	arrays     [4][]string_item
}

func (s *Strings) Init(item_count int) *Strings {
	s.item_count = item_count
	s.current = 3
	for i := 0; i < 4; i++ {
		a := make([]string_item, item_count)
		for j := 0; j < item_count; j++ {
			a[j].Init()
		}
		s.arrays[i] = a
	}
	return s
}
func in_range(id int, id_range int) bool {
	return id < id_range && id >= 0
}

func (s *Strings) Find(id int) (string, bool) {
	if !in_range(id, s.item_count) {
		return "", false
	}
	item := &s.arrays[s.current][id]
	if item.used {
		return item.value, true
	}
	return "", false
}

func (s *Strings) Read(id int, default_value string) string {
	if v, found := s.Find(id); found {
		return v
	}
	return default_value
}
func (s *Strings) Update(id int, value string) bool {
	if !in_range(id, s.item_count) {
		return false
	}
	s.arrays[(s.current+1)%4][id].Set(value)
	return true
}
func (s *Strings) Commit() { s.current = (s.current + 1) % 4 }
