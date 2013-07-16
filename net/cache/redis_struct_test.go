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
    "time"
    //"fmt"
)

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
    c := LGNewRedis()
    count, err := c.Connect("192.168.18.18:6379")
    if count == 0 || err != nil {
        t.Errorf("Connect: %v", err)
        return
    }

    a1,a2 := NewA1A2()

    // Set
    err =c.SetH("Hello", a1)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    expectRedis3("Set",t, c, "Hello", a1)


	// Delete
	err =c.Delete("Hello")
	if err != nil {
		t.Errorf("Delete: %v", err)
	}
	//expectRedis3("Delete", t, c, "Hello",nil) 

	// Flags
	err = c.Set("Hello", a1, 0xFF3F, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
    var b A
    var f uint16
    var ok bool

	f, ok =c.Get("Hello",&b)
	if !ok {
		t.Errorf("Get: %v", ok)
		return
	}
	if f != 0xFF3F {
		t.Errorf("Expecting 0xFF3F, Received %x", f)
        return
	}
	expectRedis3("Flags", t, c, "Hello", a1)

	// timeout
    //fmt.Println("timeout...")
	err = c.Set("Lost", a1, 0, 1)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expectRedis3("timeout",t, c, "Lost", a1)
	time.Sleep(2 * time.Second)

    b = A{}
	_, ok =c.Get("Lost",&b)
    ////fmt.Printf("timeout get is ",ok,b)
    if ok {
        t.Errorf("timeout : %v", ok)
		return
    } else if b.Eq(a1) {
        t.Errorf("timeout : value is read out", err)
        return
    }
    //fmt.Printf("timeout2 ....")

	// cas
	err = c.Set("Data", a1, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expectRedis3("cas1",t, c, "Data", a1)

    b = A{}
    var cas uint64
	cas,f,ok =c.Gets("Data",&b)
	if !ok {
		t.Errorf("Gets: %v", err)
		return
	}
	if cas == 0 {
		t.Errorf("Expecting non-zero for cas")
	}
	err = c.Cas("Data",a2, 12345,0 ,0)
	if err == nil {
		t.Errorf("Cas: %v", err)
		return
	}
	expectRedis3("cas2", t, c, "Data", a1)

	err = c.Cas("Data", a2, cas, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expectRedis3("cas3",t, c, "Data", a2)

	err = c.Set("Data",a1, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expectRedis3("cas4", t, c, "Data", a1)

	// stats
	_, err = c.Stats("")
	if err != nil {
		t.Errorf("Stats: %v", err)
		return
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
	err = c.Set("Flush", a2, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
	}
	expectRedis3("FlushAll1",t, c, "Flush", a2)

	err = c.FlushAll()
	if err != nil {
		t.Errorf("FlushAll: err %v", err)
		return
	}

    b = *a1
	f, ok = c.Get("Flush",&b)
	if ok {
		t.Errorf("Get: %v after FlushAll", ok)
		return
	}
	if b.Eq(a2) {
		t.Errorf("FlushAll failed")
		return
	}
    //fmt.Println("FlushAll2 ...")
    //fmt.Println("////////////////////////test 3 //////////////////////////////")
}

func expectRedis3(cmd string, t *testing.T, c *LGMemcache, key string, value *A) {
    //fmt.Println(cmd,"。。。")
    var b A
    _, ok :=c.Get(key,&b)
	if !ok {
        //fmt.Println(cmd,"///")
		t.Errorf("Get: %v", ok)
		return
	}
	if !b.Eq(value) {
        //fmt.Println(cmd,"///.")
		t.Errorf("Expecting %s, Received %s", value, b)
	}
}

