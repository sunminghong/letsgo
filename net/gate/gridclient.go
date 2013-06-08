/*=============================================================================
#     FileName: gridclient.go
#         Desc: default client of gate server receive grid server(process gridserver connection return data)
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-07 14:51:54
#      History:
=============================================================================*/
package gate

import (
    "github.com/sunminghong/letsgo/log"
)

// Client  
type GridClient struct {
    *DefaultClient

    gate *GateServer
}

func MakeGridClient (name string,transport *Transport,gate *GateServer) Client {
    log.Trace("gridclient is connect:",name)

    c := &GridClient{transport,name,&username}
    c.gate = gate
}

//对数据进行拆包
func (c *GridClient) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        code := c.Transport.Stream.Endianer.Uint16(dp.Data)
        log.Trace("msg.code:",code)

        switch dp.Type {
        case DATAPACKET_TYPE_DELAY:
            dp.dataType =DATAPACKET_TYPE_COMMON

            c.gate.Clients.Get(dp.FromCid).GetTransport().SendDP(dp)

        case DATAPACKET_TYPE_BROADCAST:
            c.gate.SendBroadcast(dp)

        default:// DATAPACKET_TYPE_COMMON
            //process msg ,eg:command line
        }
    }
}

