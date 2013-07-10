/*=============================================================================
#     FileName: memcache_test.go
#         Desc: memcache test
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-09 11:13:27
#      History:
=============================================================================*/

package cache

import (
    //"os/exec"
    . "github.com/sunminghong/letsgo/helper"
    "testing"
    "time"
    //"fmt"
)

func NewSerialize() LGISerialize {
    return LGJsonSerialize{}
    //return LGGobSerialize{}
}


func TestMemcache(t *testing.T) {
    //fmt.Println("////////////////////////test string //////////////////////////////")
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
    err =c.Set("Hello", "world", 0, 0)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    expect("Set",t, c, "Hello", "world")


	// Delete
	err =c.Delete("Hello")
	if err != nil {
		t.Errorf("Delete: %v", err)
	}
	//expect("Delete", t, c, "Hello", "")

	// Flags
	err = c.Set("Hello", "world", 0xFF3F, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
    var b string
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
	expect("Flags", t, c, "Hello", "world")

	// timeout
    //fmt.Println("timeout...")
	err = c.Set("Lost", "World", 0, 1)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect("timeout",t, c, "Lost", "World")
	time.Sleep(2 * time.Second)

    b = ""
	_, ok =c.Get("Lost",&b)
    ////fmt.Printf("timeout get is ",ok,b)
    if ok {
        t.Errorf("timeout : %v", ok)
		return
    } else if b == "World" {
        t.Errorf("timeout : value is read out", err)
        return
    }
    //fmt.Printf("timeout2 ....")

	// cas
	err = c.Set("Data", "Set", 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect("cas1",t, c, "Data", "Set")

    b = ""
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
	expect("cas2", t, c, "Data", "Set")

	err = c.Cas("Data", ("Changed"), cas, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect("cas3",t, c, "Data", "Changed")

	err = c.Set("Data",("Overwritten"), 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect("cas4", t, c, "Data", "Overwritten")

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
	err = c.Set("Flush", ("Test"), 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
	}
	expect("FlushAll1",t, c, "Flush", "Test")

	err = c.FlushAll()
	if err != nil {
		t.Errorf("FlushAll: err %v", err)
		return
	}

    b = ""
	f, ok = c.Get("Flush",&b)
	if ok {
		t.Errorf("Get: %v after FlushAll", ok)
		return
	}
	if b == "Test" {
		t.Errorf("FlushAll failed")
		return
	}
    //fmt.Println("FlushAll2 ...")
}

func expect(cmd string, t *testing.T, c *LGMemcache, key, value string) {
    //fmt.Println(cmd,"...")
    var b string
    _, ok :=c.Get(key,&b)
	if !ok {
        //fmt.Println(cmd,"///")
		t.Errorf("Get: %v", ok)
		return
	}
	if string(b) != value {
        //fmt.Println(cmd,"///.")
		t.Errorf("Expecting %s, Received %s", value, b)
	}
}
