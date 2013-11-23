/*=============================================================================
#     FileName: uidmap.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-11-19 11:41:58
#      History:
=============================================================================*/

package gate

import (
    "fmt"
    "strconv"

    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    //    "unsafe"
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

func getUidKey(uid int ) string {
    return "uid_" + strconv.Itoa(uid)
}

func getKid(gateid,fromCid,cid int) (kidstr string,checkcode int) {
    var kid string

    if fromCid > 0 {

        //下面是我独特设计，~_~
        //1.为了防止同一个gate服务器分配到同样的cid（cid==fromcid）的玩家身份
        //混淆，必须加上checkcode验证
        //2.将checkcode剥离出来用cid 作为key，就可以将uidmap的数据项控制在
        //65536（32768）个以内，因此几乎可以不用清理uidmap数据项

        _kid, _checkcode := LGParseID(fromCid)
        //kid = LGCombineID(kid, gateid)
        //kid = 0 - kid
        kid = fmt.Sprintf("%d_%d",_kid,gateid)
        checkcode = _checkcode

    } else {
        kid = strconv.Itoa(cid)
    }

    kidstr = "kid_" + kid
    return
}

func (self *LGUidMap) RemoveUid(uid int) {
    self.casCache.Delete(getUidKey(uid))
}

func (self *LGUidMap) RemoveConnectionIdByUid(uid int) {
    if uid == 0 {
        return
    }

    gateid,fromCid,cid,_ := self.CheckUid(uid)
    kidstr,_ := getKid(gateid,fromCid,cid)

    self.casCache.Deletes(kidstr,getUidKey(uid))
}

func (self *LGUidMap) RemoveConnectionId(gateid, fromCid, cid int) {
    kidstr,_ := getKid(gateid,fromCid,cid)

    uid := self.GetUid(gateid,fromCid,cid)
    if uid > 0 {
        self.casCache.Deletes(kidstr,getUidKey(uid))
    } else {
        self.casCache.Delete(kidstr)
    }
}

func (self *LGUidMap) CheckUid(uid int) (gateid, fromCid, cid int, cas uint64) {
    var v2 []int
    var err error
    //if not exists in local map object ,then read from cache read
    cas, _, err = self.casCache.Gets(getUidKey(uid), &v2)

    if err != nil {
        cas = 0
        return
    }
    gateid, fromCid, cid = v2[0], v2[1], v2[2]
    return
}

func (self *LGUidMap) GetUid(gateid, fromCid, cid int) (uid int) {
    kidstr,checkcode := getKid(gateid,fromCid,cid)

    var v2 []int
    var err error

    _, _, err = self.casCache.Gets(kidstr, &v2)

    fmt.Println("cache read error:",kidstr,err)
    if err != nil {
        return
    }

    uid, co := v2[0], v2[1]
    fmt.Println("co,checkcode:",co,checkcode,v2)
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
    kidstr,checkcode := getKid(gateid,fromCid,cid)

    v2 := []int{uid, checkcode}

    //set to cache
    err :=self.casCache.Set(kidstr, v2, 0, 0)
    fmt.Println("cache set kidstr:",kidstr,err)
    if err !=nil {
        LGTrace("saveuid():",err)
    }

    v3 := []int{gateid, fromCid, cid}
    err = self.casCache.Set(getUidKey(uid), v3, 0, 0)
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
    err := self.casCache.Cas(getUidKey(uid), v3, cas, 0, 0)
    if err != nil {
        return err
    }

    kidstr,checkcode := getKid(gateid,fromCid,cid)

    v2 := []int{uid, checkcode}

    //set to cache
    return self.casCache.Set(kidstr, v2, 0, 0)
}

func (self *LGUidMap) Clear() {
    //todo: clear all uid from memcache
}

