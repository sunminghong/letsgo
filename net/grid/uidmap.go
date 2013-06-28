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
package grid

import (
    //. "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
//    "unsafe"
)

//const (
    //kidMask = 1 << 31
    //kidMask = 1 << (unsafe.Sizeof(int(1))<<3 - 1)
//)

// map fromcid,cid to uid
// find uid by fromcid or cid
type LGUidMap struct {
    uidMap *LGMap
}

func NewLGUidMap() *LGUidMap {
    c := &LGUidMap{}

    c.uidMap = NewLGMap()
    //LGTrace("kidMask", unsafe.Sizeof(int(1))<<3-1)

    return c
}

func (self *LGUidMap) GetUid(fromCid int, cid int) int {
    var kid int
    var checkcode int

    if fromCid > 0 {
        //fromcid = gate-to-grid-clientid + checkcode
        kid, checkcode = LGParseID(fromCid)
        //kid = fromcid | kidMask
        kid = 0 - kid
    } else {
        kid = cid
    }

    if v, ok := self.uidMap.Get(kid); ok {
        if v2, ok := v.([]int); ok {
            if fromCid > 0 {
                //uid, co := LGParseID(v2)
                uid ,co := v2[0],v2[1]
                if co == checkcode {
                    return uid
                } else {
                    return 0
                }
            } else {
                return v2[0]
            }
        } else {
            return 0
        }
    }
    return 0
}

func (self *LGUidMap) SaveUid(fromCid int, cid int, uid int) {
    var kid,checkcode int

    if fromCid > 0 {
        //为了防止一个gate服务器不同的玩家分配到同样的socketid（cid==fromcid），必须加上checkcode验证
        kid , checkcode = LGParseID(fromCid)
        kid = 0 - kid

        //uid = LGCombineID(uid, checkcode)
    } else {
        kid = cid
    }
    self.uidMap.Set(kid, []int{uid,checkcode})

}
