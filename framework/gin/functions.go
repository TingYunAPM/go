// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_gin

import (
	"github.com/TingYunAPM/go"
	"github.com/gin-gonic/gin"
)

const (
	context_id = "!@#$^_tingyun"
)

//通过一个wrap过的gin.Context 获取这个Context对应的Action
func FindAction(c *gin.Context) *tingyun.Action {
	if p, found := c.Get(context_id); found {
		return p.(*tingyun.Action)
	}
	return nil
}
