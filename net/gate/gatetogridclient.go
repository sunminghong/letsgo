/*=============================================================================
#     FileName: gatetogridclient.go
#         Desc: default client of gate server receive grid server(process gridserver connection return data)
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-28 18:40:35
#      History:
=============================================================================*/
package gate

import (
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)

// Client  
type LGGateToGridClient struct {
    *LGBaseClient

    Gate *LGGateServer
    clients *LGClientMap
}
/*
func LGNewGateToGridClient (name string,transport *LGTransport) LGIClient {
    LGTrace("gridclient is connect:",name)

    c := &LGGateToGridClient{LGBaseClient:&LGBaseClient{Transport:transport,Name:name}}

    c.Register()

    return c
}*/

func (c *LGGateToGridClient) Closed() {
    gridID := c.GetTransport().Cid
    c.Gate.Dispatcher.Remove(gridID)
}

func (c *LGGateToGridClient) Register() {
    c.clients = c.Gate.Clients

    line := cmd.Register(c.Gate.Name,c.Gate.Serverid)
    //register to grid server
    dp := &LGDataPacket{
        FromCid: 0,
        Data: line,
        Type : LGDATAPACKET_TYPE_GATECONNECT,
    }

    c.Transport.SendDP(dp)
}

//对数据进行拆包
func (c *LGGateToGridClient) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        LGTrace("gategridclient.ProcessDPs():dp.type",dp.Type)
        LGTrace("c.clients",c.clients.All())
        switch dp.Type {
        case LGDATAPACKET_TYPE_DELAY:
            LGTrace("delay")

            dp.Type = LGDATAPACKET_TYPE_GENERAL
            cli := c.clients.Get(dp.FromCid)
            if cli!=nil {
                cli.GetTransport().SendDP(dp)
            }

        case LGDATAPACKET_TYPE_CLOSE:
            LGTrace("gatetogridclient.ProcessDps():close player connection:",dp.FromCid)

            if cc := c.clients.Get(dp.FromCid); cc!=nil {
                LGTrace("gatetogridclient.ProcessDps():.....closed")
                cc.Close()
            }

        case LGDATAPACKET_TYPE_BROADCAST:
            LGTrace("broadcast")
            //c.gate.SendBroadcast(c.gate.Clients.Get(dp.FromCid).GetTransport(),dp)
            c.Gate.SendBroadcast(dp)

        default:
            //process msg ,eg:command line
            c.clients.Get(dp.FromCid).GetTransport().SendDP(dp)
        }
    }
}

