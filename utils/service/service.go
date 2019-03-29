// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

//类线程封装
package service

import (
	"sync/atomic"
	"time"
)

type Service struct {
	used int32
}

func (s *Service) Stop() {
	for s.used == 0 {
		time.Sleep(1 * time.Microsecond)
	}
	atomic.AddInt32(&s.used, 1)
	for s.used <= 2 {
		time.Sleep(1 * time.Millisecond)
	}
}
func (s *Service) Start(worker func(running func() bool)) {
	s.used = 0
	go func() {
		inRunning := func() bool {
			return s.used == 1
		}
		atomic.AddInt32(&s.used, 1)
		worker(inRunning)
		atomic.AddInt32(&s.used, 1)
	}()
}
