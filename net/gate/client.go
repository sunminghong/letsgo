/*=============================================================================
#     FileName: client.go
#         Desc: default gate server receive client(process player connection)
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-07 14:50:34
#      History:
=============================================================================*/
package gate

import (
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)

// Client  
type LGClient struct {
    *LGBaseClient

    //gate *LGGateServer

    dispatcher LGIDispatcher

    Gate *LGGateServer
    grids *LGClientMap
}


func LGNewClient (name string,transport *LGTransport) LGIClient {
    LGTrace("gateclient is connect:",name)

    c := &LGClient{
        LGBaseClient:&LGBaseClient{Transport:transport,Name:name},
    }

    if gate,ok := c.Transport.Server.(*LGGateServer) ;ok {
        c.Gate = gate
        c.grids = c.Gate.Grids.Clients
    } else {
        LGError("gateserver client init error:transport.Server is not GateServer type")
    }

    c.dispatcher = c.Gate.Dispatcher
    return c
}

func (c *LGClient) Closed() {
    LGTrace("this client is closed!")

    //dispatch to one grid
    gridID,ok := c.dispatcher.Dispatch(0)
    if ok {
        gridClient := c.grids.Get(gridID)
        if gridClient != nil {
            dp := &LGDataPacket{
                Type: LGDATAPACKET_TYPE_CLOSED,
                FromCid: c.Transport.Cid,
                Data: []byte{1},
            }

            gridClient.GetTransport().SendDP(dp)
        }
    }
    return

    //在线连接数更新统计
}

//对数据进行拆包
func (c *LGClient) ProcessDPs(dps []*LGDataPacket) {
    defer func() {
        if r:=recover(); r!=nil {
            LGError("grid 服务出错：",r)
        }
    }()

    for _, dp := range dps {
        var code int
        if len(dp.Data) > 2 {
            code = int(c.Transport.Stream.Endianer.Uint16(dp.Data))
            //LGTrace("msg.code:",code)

        } else {
            code = 0
        }
        //msg := NewMessageReader(dp.Data,c.Transport.Stream.Endian)
        //dispatch to one grid
        gridID,ok := c.dispatcher.Dispatch(code)
        if ok {
            LGTrace("dispatch to gridID",gridID)
            gridClient := c.grids.Get(gridID)
            if gridClient != nil {

                dp.Type = LGDATAPACKET_TYPE_DELAY
                dp.FromCid = c.Transport.Cid

				buf := make([]byte,len(dp.Data))
				copy(buf,dp.Data)
				dp.Data = buf

                gridClient.GetTransport().SendDP(dp)

                //todo: 当grid超时处理是需要返回原协议失败
            } else {
                //todo: 是否需要缓存没有处理的数据包
                LGError("分配的grid 服务器不存在:",code)
            }
        } else {
            LGError("messageCode has not grid process:",code)
        }
    }
}
