/*=============================================================================
#     FileName: memcache.go
#         Desc: client of memcached client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-08 11:51:07
#      History:
=============================================================================*/
package cache

import (
    "code.google.com/p/vitess/go/memcache"
    . "github.com/sunminghong/letsgo/helper"
    "strings"
	"errors"
)

type LGMemcache struct {
	c *memcache.Connection
    serialize LGISerialize
	conninfo string
}

func LGNewMemcache( serialize LGISerialize) *LGMemcache {
    return &LGMemcache{serialize:serialize}
}

func (self *LGMemcache) GetRaw(key string) (val []byte, flag uint16,err error) {
	if self.c == nil {
		self.c = self.connectInit()
	}

	val, flag, err = self.c.Get(key)
	if err != nil {
        return
	}
    return
}

func (self *LGMemcache) Get(key string,val interface{}) (flag uint16,err error) {
	if self.c == nil {
		self.c = self.connectInit()
	}
	v, flag, err := self.c.Get(key)
	if err != nil {
        err = errors.New("memache read error:" + err.Error() + string(v) + "!!!")
        return
	}

    err = self.serialize.Deserialize(v,val)
    if err != nil {
        return
    }
    return
}

func (self *LGMemcache) Gets(key string,val interface{}) (cas uint64, flag uint16,err error) {
	if self.c == nil {
		self.c = self.connectInit()
	}
	v, flag, cas, err := self.c.Gets(key)
	if err != nil {
        err = errors.New("memache read error:" + err.Error() + string(v) + "!!!")
		return
	}
    err = self.serialize.Deserialize(v,val)
    if err != nil {
        return
    }
	return
}

func (self *LGMemcache) Cas(
    key string, val interface{},cas uint64, flag uint16, timeout int64) error {
	if self.c == nil {
		self.c = self.connectInit()
	}

    v,err := self.serialize.Serialize(val)
    if err != nil {
        return err
    }

	stored, err := self.c.Cas(key, flag, uint64(timeout), v,cas)
	if err != nil {
        return err
    }
    if  stored == false {
		return errors.New("stored fail")
	}
	return nil
}

func (self *LGMemcache) Set(
    key string, val interface{}, flag uint16, timeout int64) error {
	if self.c == nil {
		self.c = self.connectInit()
	}

    v,err := self.serialize.Serialize(val)
    if err != nil {
        //return errors.New("don't self.serialize this data")
        return err
    }

	stored, err := self.c.Set(key, flag, uint64(timeout), v)
	if err != nil {
        return err
    }
    if  stored == false {
		return errors.New("stored fail")
	}
	return nil
}

func (self *LGMemcache) SetRaw(
    key string, val []byte,flag uint16, timeout int64) error {
	if self.c == nil {
		self.c = self.connectInit()
	}

	stored, err := self.c.Set(key, flag, uint64(timeout), val)
	if err != nil {
        return err
    }
    if  stored == false {
		return errors.New("stored fail")
	}
	return nil
}

func (self *LGMemcache) Delete(key string) error {
	if self.c == nil {
		self.c = self.connectInit()
	}
	_, err := self.c.Delete(key)
	return err
}

func (self *LGMemcache) Deletes(keys ...string) error {
	if self.c == nil {
		self.c = self.connectInit()
	}

    var errs error
    for _,key:=range keys {
        _, err := self.c.Delete(key)
        if err !=nil {
            errs = err
        }
    }
	return errs
}

//This purges the entire cache.
func (self *LGMemcache) FlushAll() (err error) {
    return self.c.FlushAll()
}

func (self *LGMemcache) IsExist(key string) bool {
	if self.c == nil {
		self.c = self.connectInit()
	}
	v, _, err := self.c.Get(key)
	if err != nil {
		return false
	}
	if len(v) == 0 {
		return false
	} else {
		return true
	}
	return true
}

func (self *LGMemcache) Stats(argument string) (result string, err error) {
    r,err := self.c.Stats(argument)
    if err == nil {
        result = string(r)
    }
    return
}

func (self *LGMemcache) ClearAll() error {
	if self.c == nil {
		self.c = self.connectInit()
	}
	err := self.c.FlushAll()
	return err
}

func (self *LGMemcache) Connect(addrStr string) (count int, err error) {
    addrs := strings.Split(addrStr,",")

    for _,addr := range addrs {
        self.conninfo = addr
        self.c = self.connectInit()

        if self.c == nil {
            err = errors.New("dial tcp conn error")
        } else {
            count ++
        }
    }
    return count,err
}

func (self *LGMemcache) connectInit() *memcache.Connection {
	c, err := memcache.Connect(self.conninfo)
	if err != nil {
		return nil
	}
	return c
}

