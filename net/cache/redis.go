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
    "runtime"
    "net"

    "github.com/sunminghong/redis"
)

type LGRedis struct {
    *redis.Client
}

func LGNewRedis(addr string) *LGRedis{
    return &LGRedis{Client:&redis.Client{Addr:addr}}
}

func (self *LGRedis) Connect(db int) (c net.Conn, err error){
    runtime.GOMAXPROCS(2)
    //self.Client.Addr = addr
    self.Db = db

    c, err = net.Dial("tcp", self.Addr)
    if err != nil {
        c.Close()
        return
    }

    c.Close()
    return
}

func (self *LGRedis) Set(key string, val string) (err error){
    bs := []byte(val)
    return self.Client.Set(key,bs)
}

func (self *LGRedis) Get(key string) (val string, err error){
    v,err := self.Client.Get(key)
    if err == nil {
        return string(v),err
    } else {
        return "",err
    }
}

func (self *LGRedis) FlushAll() (err error){
    return self.Client.Flush(true)
}
