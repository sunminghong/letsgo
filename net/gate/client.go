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
    *DefaultClient

    //gate *GateServer

    dispatcher IDispatcher
    grids *ClientMaps
}

func MakeClient (name string,transport *Transport,gate *GateServer) Client {
    log.Trace("gateclient is connect:",name)

    c := &Client{transport,name,&username}
    //c.gate = gate

    c.dispatcher = gate.Dispatcher
    c.grids = gate.Grids.Clients
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        //msg := NewMessageReader(dp.Data,c.Transport.Stream.Endian)
        code := c.Transport.Stream.Endianer.Uint16(dp.Data)
        log.Trace("msg.code:",code)

        //dispatch to one grid
        gridID,ok = c.dispatcher.Dispatch(dp)
        if ok {
            log.Trace("dispatch to gridID",gridID)
            gridClient := c.grids.Get(gridID)

            dp := &DataPacket{Type: DATAPACKET_TYPE_DELAY, Data: data}
            dp.FromCid = c.Cid
            gridClient.Transport.SendDP(dp)
        } else {
            log.Error("messageCode has not grid process:",code)
        }
    }
}

