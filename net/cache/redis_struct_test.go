/*=============================================================================
#     FileName: memcache_struct_test.go
#         Desc: memcache test struct set/get
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-10 10:41:38
#      History:
=============================================================================*/

package cache

import (
    //"os/exec"
    //. "github.com/sunminghong/letsgo/helper"
    "testing"
    . "strconv"
    "fmt"
)

type C struct {
    I1 int
    I2 int
    I3 int
    I4 int
    I5 int
    I6 int
    I7 int
    I8 int
    I9 int
    I10 int
    I11 int
    I12 int
    I13 int
    I14 int
    I15 int

    Str1 string
}

func (a1 *C) Eq(a2 *C) bool {
    if a1.I1 != a2.I1 {
        return false
    }
    if a1.Str1 != a2.Str1 {
        return false
    }

    return true
}

func NewC() (a1 *C) {
    a1 = &C{
        I1:1,
        I2:12,
        I3:13,
        I4:1234,
        I5:1342,
        I6:123423,
        I7:123,
        I8:2341,
        I9:452234,
        I10:690,
        I11:2342,
        I12:56878,
        I13:213134,
        I14:2352,
        I15:234,

        Str1:"abc",
    }

    return a1
}


func TestRedis3(t *testing.T) {
    //fmt.Println("////////////////////////test 3 //////////////////////////////")
    //cmd := exec.Command("memcached", "-s", "/tmp/vtocc_cache.sock")
    //if err =cmd.Start(); err != nil {
    //    t.Errorf("Memcache start: %v", err)
    //    return
    //}
    //defer cmd.Process.Kill()
    //time.Sleep(time.Second)



    //c, err =Connect("/tmp/vtocc_cache.sock")
    //
    c := LGNewRedis("192.168.18.18:6379")
    _, err := c.Connect(0)
    if err != nil {
        t.Errorf("Connect: %v", err)
        return
    }

    a1:= NewC()

    // Set
    err =c.Hmset("Hello", a1)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }

    var b C
    var ok bool

	err=c.Hgetall("Hello",&b)
	if err != nil {
		t.Errorf("Get: %v", ok)
		return
	}

    if !b.Eq(a1) {
		t.Errorf("Get: %v", ok)
		return
    }

    i1,err := c.Hget("Hello","I1")
	if err != nil {
		t.Errorf("hget : %v", err)
	}

    if ii,err:=Atoi(string(i1)); err!=nil || ii!=a1.I1 {
		t.Errorf("Hget hello.i1: %v",a1.I1,i1)
    }

    //test HGet
    vals,err := c.Hmget("Hello","I1","I2","Str1")
	if err != nil {
		t.Errorf("hget : %v", err)
	}
    
    ii1,ok := vals["I1"]
    if !ok || Itoa(a1.I1)!=ii1 {
        fmt.Println("ad:"+ii1)
		t.Errorf("Hget hello.i1: %d%v",a1.I1, vals)
    }
    ii2,ok := vals["I2"]
    if !ok || Itoa(a1.I2)!=ii2 {
		t.Errorf("Hget hello.i1: %v",vals)
    }
    ii3,ok := vals["Str1"]
    if !ok || a1.Str1!=ii3 {
		t.Errorf("Hget hello.i1: %v",vals)
    }

	// Delete
	ok,err =c.Del("Hello")
	if err != nil {
		t.Errorf("Delete: %v", err)
	}

    ok,err = c.Exists("Hello")
    if ok {
		t.Errorf("Delete is error!")
    }

    /*
    //var slabs [] byte
    slabs, err := c.Stats("slabs")
	if err != nil {
		t.Errorf("Stats: %v", err)
		return
	}
    //fmt.Println("slabs:",slabs)
    */

	//FlushAll
	// Set
	err = c.Set("Flush", "sdflsadf")
	if err != nil {
		t.Errorf("Set: %v", err)
	}

	err = c.FlushAll()
	if err != nil {
		t.Errorf("FlushAll: err %v", err)
		return
	}

    _,err = c.Get("Flush")
	if err == nil {
		t.Errorf("Get: %v after FlushAll", err)
		return
	}
    //fmt.Println("FlushAll2 ...")
    //fmt.Println("////////////////////////test 3 //////////////////////////////")
}
