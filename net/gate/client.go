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
    "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)

// Client  
type Client struct {
    *BaseClient

    //gate *GateServer

    dispatcher IDispatcher
    grids *ClientMap
}

func MakeClient (name string,transport *Transport,gate *GateServer) IClient {
    log.Trace("gateclient is connect:",name)

    c := &Client{
        BaseClient:&BaseClient{Transport:transport,Name:name},
        dispatcher : gate.Dispatcher,
        grids : gate.Grids.Clients,
    }

    return c
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        //msg := NewMessageReader(dp.Data,c.Transport.Stream.Endian)
        code := c.Transport.Stream.Endianer.Uint16(dp.Data)
        log.Trace("msg.code:",code)

        //dispatch to one grid
        gridID,ok := c.dispatcher.Dispatch(dp)
        if ok {
            log.Trace("dispatch to gridID",gridID)
            gridClient := c.grids.Get(gridID)

            dp.Type = DATAPACKET_TYPE_DELAY
            dp.FromCid = c.Transport.Cid

            gridClient.GetTransport().SendDP(dp)

            //todo: 当grid超时处理是需要返回原协议失败
        } else {
            log.Error("messageCode has not grid process:",code)
        }
    }
}

