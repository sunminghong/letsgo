/*=============================================================================
#     FileName: defaultgridclient.go
#         Desc: client of default grid server receive (process player or gate connection on common)
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-07 10:40:26
#      History:
=============================================================================*/
package gate

import (
//    "reflect"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)


//define client
type LGGridProcessHandleFunc func(
    msg LGIMessageReader,c *LGGridClient,fromCid int)


// LGIClient  
type LGGridClient struct {
    *LGBaseClient

    Process LGGridProcessHandleFunc

    Gateid int
    Grid *LGGridServer
}

func (c *LGGridClient) Closed() {
    c.Grid.RemoveGate(c.Gateid,c.GetTransport().Cid)
}

//对数据进行拆包
func (c *LGGridClient) ProcessDPs(dps []*LGDataPacket) {
    stream := c.Transport.Stream
    endianer := stream.Endianer
    for _, dp := range dps {

        switch dp.Type {
        case LGDATAPACKET_TYPE_DELAY:
            LGTrace("msg.code(delay):",int(endianer.Uint16(dp.Data)),len(dp.Data))
            msg := LGNewMessageReader(dp.Data,stream.Endian)
            c.Process(msg,c,dp.FromCid)

        case LGDATAPACKET_TYPE_GENERAL:
            LGTrace("msg.code:",int(endianer.Uint16(dp.Data)),len(dp.Data))
            msg := LGNewMessageReader(dp.Data,stream.Endian)
            c.Process(msg,c,0)

        case LGDATAPACKET_TYPE_GATECONNECT:
            gatename,gateid := cmd.UnRegister(dp.Data)
            c.SetType(LGCLIENT_TYPE_GATE)

            if c.Grid == nil {
                //LGTrace("transport.server type is ",reflect.TypeOf(c.Transport.Server))
                if grid,ok := c.Transport.Server.(*LGGridServer) ;ok {
                    c.Grid = grid
                } else {
                    LGError("gridserver client init error:transport.Server is not GridServer type")
                }
            }
            c.Gateid = gateid
            c.Grid.RegisterGate(gatename,gateid,c)

            LGInfo(c.GetTransport().Conn.RemoteAddr()," is register to gate!")
        }
    }
}

func (c *LGGridClient) SendMessage(fromcid int,msg LGIMessageWriter) {
    dp := &LGDataPacket{
        FromCid: fromcid,
        Data: msg.ToBytes(),
    }
    if fromcid == 0 {
        dp.Type = LGDATAPACKET_TYPE_GENERAL
    } else {
        dp.Type = LGDATAPACKET_TYPE_DELAY
    }

    c.Transport.SendDP(dp)
}

