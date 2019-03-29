// Copyright 2016-2019 冯立强 fenglq@tingyun.com.  All rights reserved.

//无锁模式的配置信息读写
package configBase

//configBase 创建一个配置对象
//  配置对象特性: 无锁更新, 不能执行并发写操作,为多读单写无锁方案
//  配置项的更新(写) n个update + 一个 commit =>一个写事务

type configInfo map[string]interface{}
type ConfigBase struct {
	array [2]configInfo
	index int
}

//读配置项
func (c *ConfigBase) Value(name string) (interface{}, bool) {
	v, ok := c.array[c.index][name]
	return v, ok
}

//更新配置项
func (c *ConfigBase) Update(name string, value interface{}) {
	c.array[c.index^1][name] = value
}

//提交更新
func (c *ConfigBase) Commit() {
	for k, v := range c.array[c.index] {
		if _, exist := c.array[c.index^1][k]; !exist {
			c.array[c.index^1][k] = v
		}
	}
	c.index = c.index ^ 1
}

//创建一个配置对象
func CreateConfig() *ConfigBase {
	res := &ConfigBase{}
	res.array[0] = make(configInfo)
	res.array[1] = make(configInfo)
	res.index = 0
	return res
}
