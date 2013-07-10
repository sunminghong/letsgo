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
package net

import (
    //. "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    mc "code.google.com/p/vitess/go/memcache"
//    "unsafe"
)


type LGOnlineManager {
    mc.Connection
}

func (self *LGOnlineManager) Online() bool {

}

func (self *LGOnlineManager) Online() bool {

}


