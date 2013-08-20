/*=============================================================================
#     FileName: funcs.go
#         Desc: net functions
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-08-08 09:08:58
#      History:
=============================================================================*/
package gate

import (
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    "time"
)


func LGDisconnect(c *LGGridClient, gateid, fromCid, cid int, prefunc func(disconnect LGIClient)) {
    LGTrace("Disconnect is called:",gateid,fromCid,cid)

    if fromCid == 0 {
        //close direct connection
        if server, ok := c.GetTransport().Server.(*LGGridServer); ok {
            LGTrace("close direct connection")
            if dc := server.Clients.Get(cid); dc != nil {
                prefunc(dc)

                time.Sleep(300 * time.Millisecond)
                dc.Close()
            } else {
                prefunc(nil)
            }
        } else {
            prefunc(nil)
        }
        return
    }

    var dgc LGIClient
    if c.Gateid == 0 {
        //要从一个直连clientA断开一个非直连clientGB，就必须通过gateid去找到连接clientGB的clientG
        if gridserver, ok := c.GetTransport().Server.(*LGGridServer); ok {
            LGTrace("gatemap:",gridserver.GateMap)
            if cs, ok := gridserver.GateMap[gateid]; ok {
                if dc := gridserver.Clients.Get(cs[0]); dc != nil {
                    dgc = dc
                }
            }
        }
    } else {
        dgc = c
    }

    if dgc == nil {
        LGTrace("disconnect is lost:gate client is lost")

        prefunc(nil)
        return
    }

    prefunc(dgc)
    //this is connection is conneted by a gate

    //wait one second then send a command to gate to close client player connection
    time.Sleep(300 * time.Millisecond)

    dp := &LGDataPacket{
        Type:    LGDATAPACKET_TYPE_CLOSE,
        FromCid: fromCid,
        Data:    []byte{1},
    }

    dgc.GetTransport().SendDP(dp)
}
