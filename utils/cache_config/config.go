// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//配置信息缓存的性能优化版本,无锁访问, array存储
package cache_config

import (
	"errors"
	"fmt"
)

type Configuration struct {
	CStrings  Strings
	CBools    Bools
	CIntegers Integers
}

func (c *Configuration) Init(stringMax, boolMax, integerMax int) {
	c.CStrings.Init(stringMax)
	c.CBools.Init(boolMax)
	c.CIntegers.Init(integerMax)
}

func readInt(v interface{}) (error, int) {
	switch r := v.(type) {
	case float64:
		return nil, int(r)
	case float32:
		return nil, int(r)
	case int:
		return nil, r
	case int32:
		return nil, int(r)
	case int64:
		return nil, int(r)
	case uint32:
		return nil, int(r)
	case uint64:
		return nil, int(r)
	default:
		return errors.New(fmt.Sprint(v, ":  not int value.")), 0
	}
}

func (c *Configuration) Update(strings, bools, integers map[string]int, key string, value interface{}) bool {
	if index, found := strings[key]; found {
		if v, ok := value.(string); !ok {
			return false
		} else {
			return c.CStrings.Update(index, v)
		}
	}
	if index, found := bools[key]; found {
		if v, ok := value.(bool); !ok {
			return false
		} else {
			return c.CBools.Update(index, v)
		}
	}
	if index, found := integers[key]; found {
		if err, vl := readInt(value); err == nil {
			return c.CIntegers.Update(index, int64(vl))
		}
		return false
	}
	return false
}
func (c *Configuration) Commit() {
	c.CStrings.Commit()
	c.CIntegers.Commit()
	c.CBools.Commit()
}
