/*=============================================================================
#     FileName: redis_string_test.go
#         Desc: memcache test
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-16 10:56:19
#      History:
=============================================================================*/

package cache

import (
    //"os/exec"
    //. "github.com/sunminghong/letsgo/helper"
    "testing"
    //"time"
    //"fmt"
)

func TestRedis(t *testing.T) {
    //fmt.Println("////////////////////////test string //////////////////////////////")
    //cmd := exec.Command("memcached", "-s", "/tmp/vtocc_cache.sock")
    //if err =cmd.Start(); err != nil {
    //    t.Errorf("Redis start: %v", err)
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

    // Set
    err =c.Set("Hello", "world")
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    expectRedis("Set",t, c, "Hello", "world")


	// Del
    ok,err :=c.Del("Hello")
	if err != nil {
		t.Errorf("Del: %v", err)
	}
	//expectRedis("Delete", t, c, "Hello", "")

    ok,err = c.Exists("Hello")
    if ok {
		t.Errorf("Delete is error!")
    }


	//// timeout
    ////fmt.Println("timeout...")
	//err = c.Set("Lost", "World", 1)
	//if err != nil {
	//	t.Errorf("Set: %v", err)
	//	return
	//}
	//expectRedis("timeout",t, c, "Lost", "World")
	//time.Sleep(2 * time.Second)

    //b = ""
	//_, ok =c.Get("Lost",&b)
    //////fmt.Printf("timeout get is ",ok,b)
    //if ok {
    //    t.Errorf("timeout : %v", ok)
	//	return
    //} else if b == "World" {
    //    t.Errorf("timeout : value is read out", err)
    //    return
    //}
    ////fmt.Printf("timeout2 ....")


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
	err = c.Set("Flush", "Test")
	if err != nil {
		t.Errorf("Set: %v", err)
	}
	expectRedis("FlushAll1",t, c, "Flush", "Test")

	err = c.FlushAll()
	if err != nil {
		t.Errorf("FlushAll: err %v", err)
		return
	}

    b := ""
    b,err = c.Get("Flush")
	if b == "Test" {
		t.Errorf("FlushAll failed")
		return
	}
    //fmt.Println("FlushAll2 ...")
}

func expectRedis(cmd string, t *testing.T, c *LGRedis, key, value string) {
    //fmt.Println(cmd,"...")
    b,err :=c.Get(key)
    if err != nil {
		t.Errorf("Expecting %s, Received %s", err)
    }
	if string(b) != value {
        //fmt.Println(cmd,"///.")
		t.Errorf("Expecting %s, Received %s", value, b)
	}
}
