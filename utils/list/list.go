// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//A better List then native
package list

type node struct {
	pre   *node
	nxt   *node
	value interface{}
}

func (n *node) reset() {
	n.pre = nil
	n.nxt = nil
	n.value = nil
}

type Iterator struct {
	root *List
	n    *node
}

func (i Iterator) IsEnd() bool {
	return i.n == nil
}
func (i Iterator) Set(value interface{}) bool {
	if i.n == nil {
		return false
	}
	i.n.value = value
	return true
}
func (i Iterator) Valid() bool {
	if i.root == nil {
		return false
	}
	if i.n == nil {
		return true
	}
	if i.root.count == 0 {
		return false
	}
	if i.n.pre == nil {
		return i.root.first == i.n
	}
	if i.n.nxt == nil {
		return i.root.last == i.n
	}
	return true
}
func (i Iterator) Destroy() {
	i.root = nil
	i.n = nil
}
func (i Iterator) Remove() (interface{}, bool) {
	if !i.Valid() {
		return nil, false
	}
	if i.root.count == 0 {
		return nil, false
	}
	defer i.Destroy()
	if i.n.pre == nil {
		return i.root.PopFront()
	}
	if i.n.nxt == nil {
		return i.root.PopBack()
	}
	i.n.pre.nxt = i.n.nxt
	i.n.nxt.pre = i.n.pre
	i.root.count -= 1
	ret := i.n.value
	i.n.reset()
	return ret, true
}

func (i Iterator) InsertFront(value interface{}) (Iterator, bool) {
	if i.root.count == 0 {
		return i.root.PushFront(value), true
	}
	if !i.Valid() || i.IsEnd() {
		return Iterator{i.root, nil}, false
	}
	if i.n.pre == nil {
		return i.root.PushFront(value), true
	}
	n := &node{i.n.pre, i.n, value}
	n.pre.nxt = n
	i.n.pre = n
	i.root.count += 1
	return Iterator{i.root, n}, true
}

func (i Iterator) InsertBack(value interface{}) (Iterator, bool) {
	if i.root.count == 0 {
		return i.root.PushBack(value), true
	}
	if !i.Valid() || i.IsEnd() {
		return Iterator{i.root, nil}, false
	}
	if i.n.nxt == nil {
		return i.root.PushBack(value), true
	}
	n := &node{i.n, i.n.nxt, value}
	n.nxt.pre = n
	i.n.nxt = n
	i.root.count += 1
	return Iterator{i.root, n}, true
}

func (i Iterator) Value() (interface{}, bool) {
	if i.n == nil {
		return nil, false
	}
	return i.n.value, true
}
func (i Iterator) Equal(it Iterator) bool {
	return i.n == it.n && i.root == it.root
}

func (i *Iterator) MoveBack() {
	if i.root != nil && i.n != nil {
		i.n = i.n.nxt
	}
}
func (i *Iterator) MoveFront() {
	if i.root != nil && i.n != nil {
		i.n = i.n.pre
	}
}
func (i Iterator) Front() Iterator {
	if i.n == nil {
		return Iterator{i.root, nil}
	}
	return Iterator{i.root, i.n.pre}
}
func (i Iterator) Back() Iterator {
	if i.n == nil {
		return Iterator{i.root, nil}
	}
	return Iterator{i.root, i.n.nxt}
}

type List struct {
	first *node
	last  *node
	count int
}

func (l *List) Init() *List {
	l.first = nil
	l.last = nil
	l.count = 0
	return l
}
func (l *List) Size() int {
	return l.count
}
func (l *List) Front() Iterator {
	return Iterator{l, l.first}
}
func (l *List) Back() Iterator {
	return Iterator{l, l.last}
}

func (l *List) PushBack(value interface{}) Iterator {
	var n *node
	if l.count == 0 {
		n = &node{nil, nil, value}
		l.first = n
		l.last = n
	} else {
		n = &node{l.last, nil, value}
		l.last.nxt = n
		l.last = n
	}
	l.count += 1
	return Iterator{l, n}
}
func (l *List) PushFront(value interface{}) Iterator {
	var n *node
	if l.count == 0 {
		n = &node{nil, nil, value}
		l.first = n
		l.last = n
	} else {
		n = &node{nil, l.first, value}
		l.first.pre = n
		l.first = n
	}
	l.count += 1
	return Iterator{l, n}
}
func (l *List) PopFront() (interface{}, bool) {
	if l.count == 0 {
		return nil, false
	}
	rnode := l.first
	ret := rnode.value
	l.first = rnode.nxt
	rnode.value = nil
	l.count -= 1
	if l.count == 0 {
		l.last = nil
	} else {
		l.first.pre = nil
		rnode.nxt = nil
	}
	return ret, true
}
func (l *List) PopBack() (interface{}, bool) {
	if l.count == 0 {
		return nil, false
	}
	rnode := l.last
	ret := rnode.value
	l.last = rnode.pre
	rnode.value = nil
	l.count -= 1
	if l.count == 0 {
		l.first = nil
	} else {
		l.last.nxt = nil
		rnode.pre = nil
	}
	return ret, true
}
