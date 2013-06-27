/*=============================================================================
#     FileName: session.go
#         Desc: 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:52:38
#      History:
=============================================================================*/
package protos

import (
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/net/grid"
)


var Uidmap *LGUidMap = NewLGUidMap()

type Session struct {
    Username *string
}

var sessionMap *LGMap = NewLGMap()

func GetSession(uid int) *Session {
    if v,ok := sessionMap.Get(uid);ok {
        if v2,ok1 := v.(*Session); ok1 {
            return v2
        } else {
            return nil
        }
    }
    return nil
}

func SetSession(uid int,sess *Session) {
    sessionMap.Set(uid,sess)
}

