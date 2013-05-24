/*=============================================================================
#     FileName: defaultrouter.go
#         Desc: default router
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-23 17:43:26
#      History:
=============================================================================*/
package net

import (
    "strconv"
    "strings"
    "github.com/sunminghong/letsgo/helper"
)

type DefaultRouter struct {
    protomaps map[int]int
}

func (r *DefaultRouter) Init() {
    r.protomaps = make(map[int]int)
}

func (r *DefaultRouter) Add(cid int, protocols string) {
    protos := strings.Split(protocols,",")
    for _,p_ := range protos {
        p := strings.Trim(p_," ")
        if len(p) == 0 {
            continue
        }
        pcode, err := strconv.Atoi(p)
        if err ==nil {
            r.protomaps[pcode] = cid
        }
    }
}

func (r *DefaultRouter) Handler(dp DataPacket) (cid int,ok bool) {
    if dp.Type == DATAPACKET_TYPE_DELAY {
        return 0,false
    }

    proto := r.ParseProtos(dp.Code)

    cid,ok = r.protomaps[proto]

    log.Trace("router Handler func messageCode,proto,cid:",dp.Code,proto,cid)
    return cid,ok
}

func (r *DefaultRouter) ParseProtos(code uint16) int {
    return int(code / 100)
}

