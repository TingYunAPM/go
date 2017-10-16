// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.
// Use of this source code is governed by a MIT license: https://opensource.org/licenses/MIT

package tingyun_beego

import (
	"io"
	"net/http"

	"github.com/TingYunAPM/go"
)

//对http.Request.Body的hook,用来保证tingyun.Action在Web过程中的唯一性
type bodyWrapper struct {
	Body   io.ReadCloser
	action *tingyun.Action
}

func (b *bodyWrapper) Close() error {
	return b.Body.Close()
}
func (b *bodyWrapper) Read(p []byte) (n int, err error) {
	return b.Body.Read(p)
}
func wrapRequest(r *http.Request, a *tingyun.Action) bool {
	if isWraped(r) {
		return false
	}
	body := r.Body
	r.Body = &bodyWrapper{Body: body, action: a}
	return true
}
func unWrapRequest(r *http.Request) {
	if p, ok := r.Body.(*bodyWrapper); ok {
		r.Body = p.Body
		p.Body = nil
		p.action = nil
	}
}
func getActionByRequest(r *http.Request) *tingyun.Action {
	if p, ok := r.Body.(*bodyWrapper); ok {
		return p.action
	}
	return nil
}
func isWraped(r *http.Request) bool {
	if _, ok := r.Body.(*bodyWrapper); ok {
		return true
	}
	return false
}
