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
    //. "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
//    "unsafe"
    "strconv"
    //"fmt"
)

type iCache interface {
    Gets(key string,val interface{}) (cas uint64, flag uint16,err error)

    Cas(
        key string, val interface{},cas uint64, timeout int64, flag uint16) error
}

// map fromcid,cid to uid
// find uid by fromcid or cid
type LGUidMap struct {
    uidMap *LGMap
    casCache iCache
}

func LGNewUidMap(cache iCache) *LGUidMap {
    c := &LGUidMap{casCache:cache}

    c.uidMap = NewLGMap()

    return c
}

func (self *LGUidMap) GetUid(gateid, fromCid, cid int) (uid int,cas uint64) {
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

    var v2 []uint64
    var ok bool
    var err error
    var v interface{}
    if v, ok = self.uidMap.Get(kid); !ok {
        //if not exists in local map object ,then read from cache read
        cas,_,err = self.casCache.Gets(
            "uid_" + strconv.Itoa(kid), &v2)

        if err == nil {
            ok = true
        }
        //fmt.Println("v2,err:",v2,err,cas)
    } else {
        //fmt.Println("has key,v,ok:",v,ok)
        v2, ok = v.([]uint64)
    }

    if !ok {
        return
    }

    if cas > 0 {
        v2 = append(v2,cas)
        self.uidMap.Set(kid, v2)
    }

    uid_ ,co,cas := v2[0],v2[1],v2[2]
    //fmt.Println("co,checkcode:",co,checkcode,v2)
    if fromCid != 0 {
        if co == uint64(checkcode) {
            return int(uid_),cas
        } else {
            return 0,0
        }
    } else {
        return int(uid_),cas
    }
}

func (self *LGUidMap) SaveUid(gateid, fromCid, cid, uid int,cas uint64) error {
    var kid,checkcode int

    if fromCid > 0 {

        //下面是我独特设计，~_~
        //1.为了防止同一个gate服务器分配到同样的cid（cid==fromcid）的玩家身份
        //混淆，必须加上checkcode验证
        //2.将checkcode剥离出来用cid 作为key，就可以将uidmap的数据项控制在
        //65536（32768）个以内，因此几乎可以不用清理uidmap数据项

        kid , checkcode = LGParseID(fromCid)
        kid = LGCombineID(kid, gateid)
        kid = 0 - kid

    } else {
        kid = cid
    }

    v2 := []int{uid,checkcode}
    self.uidMap.Delete(kid)

    //set to cache
    return self.casCache.Cas("uid_" + strconv.Itoa(kid),
        v2, cas, 0,0)
}

func (self *LGUidMap) Clear() {
    self.uidMap.Clear()
}
