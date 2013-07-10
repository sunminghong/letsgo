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

type A struct {
    I1 int
    Str1 string
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

    B1 *B
    Barr []*B
    Bmap2 map[string]B
    Bmap1 map[string]*B
}

type B struct {
    I2 int
    Str2 string
}

func (a1 *A) Eq(a2 *A) bool {
    if a1.I1 != a2.I1 {
        return false
    }
    if a1.Str1 != a2.Str1 {
        return false
    }
    for i,B1 := range a1.Barr {
        b2 := a2.Barr[i]
        if !B1.Eq(b2) {
            return false
        }
    }

    for i,B1 := range a1.Bmap1 {
        b2 := a2.Bmap1[i]
        if !B1.Eq(b2) {
            return false
        }
    }

    for i,B1 := range a1.Bmap2 {
        b2 := a2.Bmap2[i]
        if !B1.Eq(&b2) {
            return false
        }
    }
    return true
}

func (B1 *B) Eq(b2 *B) bool {
    if B1.I2 != b2.I2 {
        return false
    }
    if B1.Str2 != b2.Str2 {
        return false
    }
    return true
}

func NewA1A2() (a1,a2 *A) {
    //Bmap1 := make(map[int]*B)
    Bmap1 := map[string]*B {
        "11":&B{2,"def"},
        "22":&B{22,"defdef"},
        "33":&B{222,"defdefdef"},
    }

    Bmap2 := map[string]B {
        "aa":B{3,"ghi"},
        "bb":B{33,"ghighi"},
        "cc":B{333,"ghighighi"},
    }

    a1 = &A{
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
        Barr:[]*B{&B{2,"def"},&B{22,"defdef"},&B{222,"defdefdef"}},
        Bmap1:Bmap1,
        Bmap2:Bmap2,
    }

    a2 = &A{
        I1:11,
        Str1:"abcabc",
        Barr:[]*B{&B{22,"def"},&B{2222,"defdef"},&B{222222,"defdefdef"}},
        Bmap1:Bmap1,
        Bmap2:Bmap2,
    }
    return a1,a2
}

func TestMemcache3(t *testing.T) {
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
    c := LGNewMemcache(NewSerialize())
    count, err := c.Connect("192.168.18.18:11211")
    if count == 0 || err != nil {
        t.Errorf("Connect: %v", err)
        return
    }

    a1,a2 := NewA1A2()

    // Set
    err =c.Set("Hello", a1, 0, 0)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }
    expect3("Set",t, c, "Hello", a1)


	// Delete
	err =c.Delete("Hello")
	if err != nil {
		t.Errorf("Delete: %v", err)
	}
	//expect3("Delete", t, c, "Hello",nil) 

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
	expect3("Flags", t, c, "Hello", a1)

	// timeout
    //fmt.Println("timeout...")
	err = c.Set("Lost", a1, 0, 1)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect3("timeout",t, c, "Lost", a1)
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
	expect3("cas1",t, c, "Data", a1)

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
	expect3("cas2", t, c, "Data", a1)

	err = c.Cas("Data", a2, cas, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect3("cas3",t, c, "Data", a2)

	err = c.Set("Data",a1, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
	expect3("cas4", t, c, "Data", a1)

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
	expect3("FlushAll1",t, c, "Flush", a2)

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

func expect3(cmd string, t *testing.T, c *LGMemcache, key string, value *A) {
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
