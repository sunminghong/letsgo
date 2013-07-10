/*=============================================================================
#     FileName: memcache_struct_b_test.go
#         Desc: memcache test struct set/get
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-10 18:41:02
#      History:
=============================================================================*/

package cache

import (
    //"os/exec"
    "testing"
    "fmt"
)

var c *LGMemcache

func Init() {
    c = LGNewMemcache(NewSerialize())
    count, err := c.Connect("192.168.18.18:11211")
    if count == 0 || err != nil {
        fmt.Println("Init is error:",err)
        return
    }

	c.FlushAll()
}

func Benchmark_Get(t *testing.B) {
    Init()

    a1,a2 := NewA1A2()

    c.Set("Hello", a1, 0, 0)
    for i:=0;i<t.N;i++ {
        _, ok :=c.Get("Hello",a2)
        if !ok {
            t.Errorf("Get: %v", ok)
            return
        }
    }
}

func Benchmark_Gets(t *testing.B) {
    Init()
    a1,a2 := NewA1A2()

    c.Set("Hello", a1, 0, 0)
    for i:=0;i<t.N;i++ {
        _,_, ok :=c.Gets("Hello",a2)
        if !ok {
            t.Errorf("Get: %v", ok)
            return
        }
    }
}

func Benchmark_Set(t *testing.B) {
    Init()
    a1,_ := NewA1A2()

    for i:=0;i<t.N;i++ {
        err :=c.Set("Hello", a1, 0, 0)
        if err != nil {
            t.Errorf("Set: %v", err)
            return
        }
    }
}

func Benchmark_Delete(t *testing.B) {
    Init()

    for i:=0;i<t.N;i++ {
        err :=c.Delete("Hello")
        if err != nil {
            t.Errorf("Delete: %v", err)
            return
        }
    }
}

func Benchmark_CasNo(t *testing.B) {
    Init()
    a1,a2 := NewA1A2()

    // cas
    err := c.Set("Data", a1, 0, 0)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    _,_,ok :=c.Gets("Data",a2)
    if !ok {
        t.Errorf("Gets: %v", err)
        return
    }
    for i:=0;i<t.N;i++ {
        err := c.Cas("Data",a2, 12345,0 ,0)
        if err == nil {
            t.Errorf("Cas: %v", err)
            return
        }
    }
}

func Benchmark_CasYes(t *testing.B) {
    Init()
    a1,a2 := NewA1A2()

    // cas
    err := c.Set("Data", a1, 0, 0)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    cas,_,ok :=c.Gets("Data",a2)
    if !ok {
        t.Errorf("Gets: %v", err)
        return
    }
    for i:=0;i<t.N;i++ {
        c.Cas("Data", a2, cas, 0, 0)
    }
}

