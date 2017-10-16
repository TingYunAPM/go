// Copyright 2017 冯立强 fenglq@tingyun.com.  All rights reserved.

//配置信息缓存的性能优化版本,无锁访问, array存储
package cache_config

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
		if v, ok := value.(int); !ok {
			return false
		} else {
			return c.CIntegers.Update(index, int64(v))
		}
	}
	return false
}
func (c *Configuration) Commit() {
	c.CStrings.Commit()
	c.CIntegers.Commit()
	c.CBools.Commit()
}
