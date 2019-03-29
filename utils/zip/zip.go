// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

//zip封装
package zip

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

//Deflate 数据压缩
func Deflate(b []byte) ([]byte, error) {
	r := bytes.Buffer{}
	w := zlib.NewWriter(&r)
	_, err := w.Write(b)
	w.Close()
	if err == nil {
		return r.Bytes(), nil
	}
	return nil, err
}

//Inflate 解压缩数据
func Inflate(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}
