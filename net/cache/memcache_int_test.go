/*=============================================================================
#     FileName: memcache_int_test.go
#         Desc: memcache test int type
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-10 10:41:10
#      History:
=============================================================================*/

package cache

import (
    //"os/exec"
    "testing"
    "time"
    //"fmt"
)

func TestMemcache2(t *testing.T) {
    //fmt.Println("////////////////////////test int //////////////////////////////")
    //cmd := exec.Command("memcached", "-s", "/tmp/vtocc_cache.sock")
    //if err =cmd.Start(); err != nil {
    //    t.Errorf("Memcache start: %v", err)
    //    return
    //}
    //defer cmd.Process.Kill()
    //time.Sleep(time.Second)



    //c, err =Connect("/tmp/vtocc_cache.sock")
    //
    c := LGNewMemcache(NewSerialize())
    count, err := c.Connect("192.168.18.18:11211")
    if count == 0 || err != nil {
        t.Errorf("Connect: %v", err)
        return
    }

    // Set
    err =c.Set("Hello", 1223, 0, 0)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    expect2("Set",t, c, "Hello", 1223)


	// Delete
	err =c.Delete("Hello")
	if err != nil {
		t.Errorf("Delete: %v", err)
	}
	//expect2("Delete", t, c, "Hello",0) 

	// Flags
	err = c.Set("Hello", 1223, 0xFF3F, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
    var b int
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
	expect2("Flags", t, c, "Hello", 1223)

	// timeout
    //fmt.Println("timeout...")
	err = c.Set("Lost", 4567, 0, 1)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect2("timeout",t, c, "Lost", 4567)
	time.Sleep(2 * time.Second)

    b = 0
	_, ok =c.Get("Lost",&b)
    ////fmt.Printf("timeout get is ",ok,b)
    if ok {
        t.Errorf("timeout : %v", ok)
		return
    } else if b == 4567 {
        t.Errorf("timeout : value is read out", err)
        return
    }
    //fmt.Printf("timeout2 ....")

	// cas
	err = c.Set("Data", 5555, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect2("cas1",t, c, "Data", 5555)

    b =0
    var cas uint64
	cas,f,ok =c.Gets("Data",&b)
	if !ok {
		t.Errorf("Gets: %v", err)
		return
	}
	if cas == 0 {
		t.Errorf("Expecting non-zero for cas")
	}
	err = c.Cas("Data",("not set"), 12345,0 ,0)
	if err == nil {
		t.Errorf("Cas: %v", err)
		return
	}
	expect2("cas2", t, c, "Data", 5555)

	err = c.Cas("Data", 3333, cas, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect2("cas3",t, c, "Data", 3333)

	err = c.Set("Data",22222, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect2("cas4", t, c, "Data", 22222)

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
	err = c.Set("Flush", 1111, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
	}
	expect2("FlushAll1",t, c, "Flush", 1111)

	err = c.FlushAll()
	if err != nil {
		t.Errorf("FlushAll: err %v", err)
		return
	}

    b = 0
	f, ok = c.Get("Flush",&b)
	if ok {
		t.Errorf("Get: %v after FlushAll", ok)
		return
	}
	if b == 1111 {
		t.Errorf("FlushAll failed")
		return
	}
    //fmt.Println("FlushAll2 ...")
}

func expect2(cmd string, t *testing.T, c *LGMemcache, key string, value int) {
    //fmt.Println(cmd,"。。。")
    var b int
    _, ok :=c.Get(key,&b)
	if !ok {
        //fmt.Println(cmd,"///")
		t.Errorf("Get: %v", ok)
		return
	}
	if b != value {
        //fmt.Println(cmd,"///.")
		t.Errorf("Expecting %s, Received %s", value, b)
	}
}
