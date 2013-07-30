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

package gate

import (
    "testing"
    "fmt"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/net/cache"
)

func NewSerialize() LGISerialize {
    return LGJsonSerialize{}
    //return LGGobSerialize{}
}

func TestUidmap(t *testing.T) {
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
    serialize := NewSerialize()
    c := LGNewMemcache(serialize)
    count, err := c.Connect("192.168.18.18:11211")
    if count == 0 || err != nil {
        t.Errorf("Connect: %v", err)
        return
    }
	err = c.Set("Data", 5555, 0, 0)
	if err != nil {
		t.Errorf("Set: %v", err)
		return
	}
    var ii int
    cas,_,err := c.Gets("Data",&ii)
    fmt.Println("data",err,cas,ii)

    err = c.Cas("Data",6666,cas,0,0)
    fmt.Println(err)
    cas,_,err = c.Gets("Data",&ii)
    fmt.Println("data",err,cas,ii)

    err = c.Cas("Data",7777,3,0,0)
    fmt.Println(err)
    cas,_,err = c.Gets("Data",&ii)
    fmt.Println("data",err,cas,ii)


    uid := 12345
    gateid:=1
    fromcid := LGCombineID(13,545)
    cid := 3

    cas = 0
    uidmap := LGNewUidMap(c)
    uidmap.SaveUid(gateid,fromcid,cid,uid)
    fmt.Println("/////////////////////")

    uid = uidmap.GetUid(gateid,fromcid,cid)
    if uid!=12345 {
        t.Errorf("uid read error:%d,%d",12345,uid)
    }
    fmt.Println("cas==",cas)


    //uidmap.Clear()

    gateid_,fromcid_,cid_,cas := uidmap.CheckUid(uid)
    if gateid_!=gateid || fromcid_ !=fromcid || cid_!=cid {
        t.Errorf("uid read error:%d,%d",12345,uid)
    }
    fmt.Println("cas=2=",gateid_,fromcid_,cid_,cas)


    err = uidmap.CasUid(gateid,fromcid,cid,2345235,32234)
    fmt.Println(err)
    uid = uidmap.GetUid(gateid,fromcid,cid)

    if uid!=12345 {
        t.Errorf("uid read error:",12345,uid)
    }
    fmt.Println("cas==",cas)

    err = uidmap.CasUid(23423,fromcid,cid,2345235,cas)
    fmt.Println(err)
    uid = uidmap.GetUid(gateid,fromcid,cid)

    if uid==2345235{
        t.Errorf("uid read error:",12345,uid)
    }

    fmt.Println("cas==",cas)
    gateid_,fromcid_,cid_,cas = uidmap.CheckUid(uid)
    if gateid_!=23423 || fromcid_ !=fromcid || cid_!=cid {
        t.Errorf("uid read error:%d,%d",12345,uid)
    }
    fmt.Println("cas=2=",gateid_,fromcid_,cid_,cas)
}

