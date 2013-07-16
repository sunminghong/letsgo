/*=============================================================================
#     FileName: redis.go
#         Desc: client of memcached client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-15 16:00:08
#      History:
=============================================================================*/
package cache

import (
	"github.com/garyburd/redigo/redis"
    //"github.com/hoisie/redis"
    . "github.com/sunminghong/letsgo/helper"
    "strings"
	"errors"
)

var (
	DefaultKey string = "beecacheRedis"
)

type LGRedis struct {
	c        redis.Conn
	serialize LGISerialize
	conninfo string
}

func LGNewRedis() *LGRedis{
	return &LGRedis{}
}

func (rc *LGRedis) GetH(key string,fieldname string) interface{} {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Do("HGET", key, fieldname)
	if err != nil {
		return nil
	}
	return v
}

func (rc *LGRedis) GetFlat(key string,val interface{}) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Do("HGET", key)
	if err != nil {
		return err
	}

	p, ok := v.([]interface{})
	if !ok {
		return errors.New("redigo: ScanStruct expectes multibulk reply")
	}

    err = redis.ScanStruct(p,val)
	return err
}

// return map[string]string
func (rc *LGRedis) GetHM(key string,fieldnames []string) (map[string]string,error) {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}

    fs := []string{}
    fs = append(fs,key)
    copy(fs[1:],fieldnames[:])

	v, err := rc.c.Do("HMGET", fs)
	if err != nil {
		return nil,err
	}

    return ScanMap(v)
}

//set hmap set 
func (rc *LGRedis) SetHM(key string, val interface{}) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}

    vs := redis.Args{}.AddFlat(val)
    _, err := rc.c.Do("HMSET",vs)
    return err
}

//set hmap set 
func (rc *LGRedis) SetH(key string, val interface{}) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}

    vs := redis.Args{}.AddFlat(val)
    _, err := rc.c.Do("HMSET",vs)
    return err
}

//delte hmap set
func (rc *LGRedis) DeleteH(key string,fieldname string) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("HDEL", key,fieldname)
	return err
}

func (rc *LGRedis) IsExistH(key string,fieldname string) bool {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := redis.Bool(rc.c.Do("HEXISTS", key,fieldname))
	if err != nil {
		return false
	}
	return v
}

func (rc *LGRedis) Get(key string) string {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := rc.c.Do("GET", key)
	if err != nil {
		return ""
	}

    if v2,ok := v.(string); ok {
        return v2
    } else {
        return ""
    }
}

func (rc *LGRedis) Set(key string, val interface{}, timeout int64) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("SET", key, val)
	return err
}

func (rc *LGRedis) Delete(key string) error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("DEL", key)
	return err
}

func (rc *LGRedis) IsExist(key string) bool {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	v, err := redis.Bool(rc.c.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc *LGRedis) FlushAll() error {
	if rc.c == nil {
		rc.c = rc.connectInit()
	}
	_, err := rc.c.Do("FLUSHALL")
	return err
}

func (self *LGRedis) Connect(addrStr string) (count int, err error) {
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

func (rc *LGRedis) connectInit() redis.Conn {
	c, err := redis.Dial("tcp", rc.conninfo)
	if err != nil {
		return nil
	}
	return c
}

func ScanMap(reply interface{}) (map[string]string, error) {
    vs := make(map[string]string)
	p, ok := reply.([]interface{})
	if !ok {
		return nil, errors.New("redigo: ScanStruct expectes multibulk reply")
	}

	for i := 0; i < len(p); i += 2 {
		name, ok := p[i].([]byte)
		if !ok {
			return vs,errors.New("redigo: ScanStruct key not a bulk value")
		}
		value, ok := p[i+1].([]byte)
		if !ok {
			return vs,errors.New("redigo: ScanStruct value not a bulk value")
		}

        vs[string(name)] = string(value)

	}
	return vs,nil
}

