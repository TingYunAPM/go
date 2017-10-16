// Copyright 2016-2017 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun

import (
	"time"
)

type errInfo struct {
	happenTime time.Time
	e          interface{}
	stack      string
	eType      string
}

func (i *errInfo) Destroy() {
	i.e = nil
	i.stack = ""
}
