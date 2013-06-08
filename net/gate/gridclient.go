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
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)

// Client  
type LGGridClient struct {
    *LGBaseClient

    gate *LGGateServer
}

func LGMakeLGGridClient (name string,transport *LGTransport,gate *LGGateServer) LGIClient {
    LGTrace("gridclient is connect:",name)

    c := &LGGridClient{LGBaseClient:&LGBaseClient{Transport:transport,Name:name}}
    c.gate = gate

    //register to grid server
    dp := &LGDataPacket{
        FromCid: 0,
        Data: []byte{1},
        Type : LGDATAPACKET_TYPE_GATECONNECT,
    }

    transport.SendDP(dp)

    return c
}

//对数据进行拆包
func (c *LGGridClient) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        code := c.Transport.Stream.Endianer.Uint16(dp.Data)
        LGTrace("msg.code:",code)

        switch dp.Type {
        case LGDATAPACKET_TYPE_DELAY:
            dp.Type =LGDATAPACKET_TYPE_GENERAL

            c.gate.Clients.Get(dp.FromCid).GetTransport().SendDP(dp)

        case LGDATAPACKET_TYPE_BROADCAST:
            //c.gate.SendBroadcast(c.gate.Clients.Get(dp.FromCid).GetTransport(),dp)
            c.gate.SendBroadcast(nil,dp)

        default:
            //process msg ,eg:command line
        }
    }
}

