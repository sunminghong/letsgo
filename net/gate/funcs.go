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

func LGGetConnection(c *LGGridConnection, gateid, cid int) *LGGridConnection {
    if gateid==0 {
        //close direct connection
        if server, ok := c.GetTransport().Server.(*LGGridServer); ok {
            if dc := server.Connections.Get(cid); dc != nil {
                if dgc,ok := dc.(*LGGridConnection);ok {
                    return dgc
                } else {
                    return nil
                }
            }
            return nil
        } else {
            return nil
        }
    }

    if c.GateId != gateid{
        //要从一个直连clientA断开一个非直连clientGB，就必须通过gateid去找到连接clientGB的clientG
        if gridserver, ok := c.GetTransport().Server.(*LGGridServer); ok {
            LGTrace("gatemap:", gridserver.GateMap)
            if cs, ok := gridserver.GateMap[gateid]; ok {
                if dc := gridserver.Connections.Get(cs[0]); dc != nil {
                    if dgc,ok := dc.(*LGGridConnection);ok {
                        return dgc
                    } else {
                        return nil
                    }
                }
            }
        }

    } else {
        return c
    }
    return nil
}

func LGDisconnect(c *LGGridConnection, gateid, fromCid, cid int, prefunc func(disconnect LGIConnection)) {
    LGTrace("Disconnect is called:",gateid,fromCid,cid)

    dgc := LGGetConnection(c,gateid,cid)
    if dgc == nil {
        LGTrace("disconnect is lost:gate client is lost")

        prefunc(nil)
        return
    }

    if fromCid == 0 {
        prefunc(dgc)

        time.Sleep(200 * time.Millisecond)
        dgc.Close()
        return
    }

    prefunc(dgc)
    //this is connection is conneted by a gate

    //wait one second then send a command to gate to close client player connection
    time.Sleep(200 * time.Millisecond)

    dp := &LGDataPacket{
        Type:    LGDATAPACKET_TYPE_CLOSE,
        FromCid: fromCid,
        Data:    []byte{1},
    }

    dgc.GetTransport().SendDP(dp)
}

func LGSendMessage(c *LGGridConnection,gateid int,fromCid int,cid int,msg LGIMessageWriter) {
    dgc := LGGetConnection(c,gateid,cid)
    if dgc == nil {
        return
    }

    dp := &LGDataPacket{
        FromCid: fromCid,
        Data: msg.ToBytes(),
    }

    if fromCid == 0 {
        dp.Type = LGDATAPACKET_TYPE_GENERAL
    } else {
        dp.Type = LGDATAPACKET_TYPE_DELAY
    }

    dgc.GetTransport().SendDP(dp)
}
