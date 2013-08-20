/*=============================================================================
#     FileName: uidmap.go
#         Desc: client of default grid server receive (process player or gate connection on common)
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 18:03:17
#      History:
=============================================================================*/
package gate

import (
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    //    "unsafe"
    //"fmt"
    "strconv"
)

type iCache interface {
    Gets(key string, val interface{}) (cas uint64, flag uint16, err error)

    Set(key string, val interface{}, flag uint16, timeout int64) error

    Delete(key string) error
    Deletes(key ...string) error

    Cas(
        key string, val interface{}, cas uint64, flag uint16, timeout int64) error
}

// map fromcid,cid to uid
// find uid by fromcid or cid
type LGUidMap struct {
    casCache iCache
}

func LGNewUidMap(cache iCache) *LGUidMap {
    c := &LGUidMap{casCache: cache}
    return c
}

func (self *LGUidMap) RemoveUid(uid int) {
    self.casCache.Delete("uid_" + strconv.Itoa(uid))
}

func (self *LGUidMap) RemoveKid(gateid, fromCid, cid int) {

    var kid int

    if fromCid > 0 {
        //fromcid = gate-to-grid-clientid + checkcode
        kid, _= LGParseID(fromCid)
        kid = LGCombineID(kid, gateid)
        kid = 0 - kid
    } else {
        kid = cid
    }

    kkid := "kid_" + strconv.Itoa(kid)
    self.casCache.Delete(kkid)
}

func (self *LGUidMap) CheckUid(uid int) (gateid, fromCid, cid int, cas uint64) {
    var v2 []int
    var err error
    //if not exists in local map object ,then read from cache read
    cas, _, err = self.casCache.Gets(
        "uid_"+strconv.Itoa(uid), &v2)

    if err != nil {
        cas = 0
        return
    }
    gateid, fromCid, cid = v2[0], v2[1], v2[2]

    return
}

func (self *LGUidMap) GetUid(gateid, fromCid, cid int) (uid int) {
    var kid int
    var checkcode int

    if fromCid > 0 {
        //fromcid = gate-to-grid-clientid + checkcode
        kid, checkcode = LGParseID(fromCid)
        kid = LGCombineID(kid, gateid)
        kid = 0 - kid
    } else {
        kid = cid
    }

    var v2 []int
    var ok bool
    var err error

    _, _, err = self.casCache.Gets(
        "kid_"+strconv.Itoa(kid), &v2)

    if err == nil {
        ok = true
    }

    if !ok {
        return
    }

    uid, co := v2[0], v2[1]
    //fmt.Println("co,checkcode:",co,checkcode,v2)
    if fromCid != 0 {
        if co == checkcode {
            return uid
        } else {
            return 0
        }
    } else {
        return uid
    }
}

func (self *LGUidMap) SaveUid(gateid, fromCid, cid, uid int) error {
    var kid, checkcode int

    if fromCid > 0 {

        //下面是我独特设计，~_~
        //1.为了防止同一个gate服务器分配到同样的cid（cid==fromcid）的玩家身份
        //混淆，必须加上checkcode验证
        //2.将checkcode剥离出来用cid 作为key，就可以将uidmap的数据项控制在
        //65536（32768）个以内，因此几乎可以不用清理uidmap数据项

        kid, checkcode = LGParseID(fromCid)
        kid = LGCombineID(kid, gateid)
        kid = 0 - kid

    } else {
        kid = cid
    }

    v2 := []int{uid, checkcode}

    //set to cache
    err :=self.casCache.Set("kid_"+strconv.Itoa(kid), v2, 0, 0)
    if err !=nil {
        LGTrace("saveuid():",err)
    }

    v3 := []int{gateid, fromCid, cid}
    err = self.casCache.Set("uid_"+strconv.Itoa(uid), v3, 0, 0)
    if err !=nil {
        LGTrace("saveuid():",err)
    }

    return err
}

func (self *LGUidMap) CasUid(gateid, fromCid, cid, uid int, cas uint64) error {
    //if cas == 0 {
    //    return self.SaveUid(gateid,fromCid,cid,uid)
    //}
    v3 := []int{gateid, fromCid, cid}
    err := self.casCache.Cas("uid_"+strconv.Itoa(uid), v3, cas, 0, 0)
    if err != nil {
        return err
    }

    var kid, checkcode int

    if fromCid > 0 {

        //下面是我独特设计，~_~
        //1.为了防止同一个gate服务器分配到同样的cid（cid==fromcid）的玩家身份
        //混淆，必须加上checkcode验证
        //2.将checkcode剥离出来用cid 作为key，就可以将uidmap的数据项控制在
        //65536（32768）个以内，因此几乎可以不用清理uidmap数据项

        kid, checkcode = LGParseID(fromCid)
        kid = LGCombineID(kid, gateid)
        kid = 0 - kid

    } else {
        kid = cid
    }

    v2 := []int{uid, checkcode}

    //set to cache
    return self.casCache.Set("kid_"+strconv.Itoa(kid), v2, 0, 0)
}

func (self *LGUidMap) Clear() {
    //todo: clear all uid from memcache
}

