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

// Connection  
type LGGateToGridConnection struct {
    *LGBaseConnection

    Gate *LGGateServer
    clients *LGConnectionMap
}
/*
func LGNewGateToGridConnection (name string,transport *LGTransport) LGIConnection {
    LGTrace("gridclient is connect:",name)

    c := &LGGateToGridConnection{LGBaseConnection:&LGBaseConnection{Transport:transport,Name:name}}

    c.Register()

    return c
}*/

func (c *LGGateToGridConnection) Closed() {
    gridID := c.GetTransport().Cid
    c.Gate.Dispatcher.Remove(gridID)
}

func (c *LGGateToGridConnection) Register() {
    c.clients = c.Gate.Connections

    line := cmd.Register(c.Gate.Name,c.Gate.Serverid)
    //register to grid server
    dp := &LGDataPacket{
        FromCid: 0,
        Data: line,
        Type : LGDATAPACKET_TYPE_GATE_REGISTER,
    }

    c.Transport.SendDP(dp)
}

//对数据进行拆包
func (c *LGGateToGridConnection) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        LGTrace("gategridclient.ProcessDPs():dp.type=%d,fromcid=% X,len(data)=%d",dp.Type,dp.FromCid,len(dp.Data))
        //LGTrace("c.clients",c.clients.All())

		buf := make([]byte,len(dp.Data))
		copy(buf,dp.Data)
		dp.Data = buf
        if dp.Type == LGDATAPACKET_TYPE_BROADCAST {
            LGTrace("broadcast")
            //c.gate.SendBroadcast(c.gate.Connections.Get(dp.FromCid).GetTransport(),dp)
            c.Gate.SendBroadcast(dp)
            return
        }

        cli := c.clients.Get(dp.FromCid)
        if cli == nil {
            LGTrace("dp lost,fromcid:% X",dp.FromCid)
            return
        }

        switch dp.Type {
        case LGDATAPACKET_TYPE_DELAY:
            LGTrace("dp.Type = delay msg:% X",dp.Data)

            dp.Type = LGDATAPACKET_TYPE_GENERAL
            cli.GetTransport().SendDP(dp)

        case LGDATAPACKET_TYPE_DELAY_DATAS_COMPRESS:
            LGTrace("delay compress datas:%d,\n% X",dp.FromCid,dp.Data)

            dp.Type = LGDATAPACKET_TYPE_DATAS_COMPRESS
            cli.GetTransport().SendDP(dp)

        case LGDATAPACKET_TYPE_DELAY_DATAS:
            LGTrace("delay datas:%d,\n% X",dp.FromCid,dp.Data)

            LGTrace("delay datas is send")
            cli.GetTransport().SendBytes(dp.Data)

        case LGDATAPACKET_TYPE_CLOSE:
            LGTrace("gatetogridclient.ProcessDps():close player connection:%d",dp.FromCid)
            cli.Close()
        }
    }
}

