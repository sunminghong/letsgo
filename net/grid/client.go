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
package grid

import (
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)


// LGIClient  
type LGGridClient struct {
    *LGBaseClient

    Process LGProcessHandleFunc
}

/*
func LGProccessHandle(code int,msg *MessageReader,c LGIClient,fromCid int) {
    fmt.Println("message is request")
}

func LGMakeLGBaseClient (name string,transport *LGTransport) LGIClient {
    c := &LGBaseClient{
        LGBaseClient:&LGBaseClient{transport,name,LGCLIENT_TYPE_GENERAL},
    }
    c.Process = LGProcessHandleFunc
}*/

//对数据进行拆包
func (c *LGGridClient) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        LGTrace("msg.code:",code,len(dp.Data))

        msg := LGNewMessageReader(dp.Data,c.Transport.Stream.Endian)

        switch dp.Type {
        case LGDATAPACKET_TYPE_DELAY:
            c.Process(code,msg,c,dp.FromCid)

        case LGDATAPACKET_TYPE_GATECONNECT:
            c.SetType(LGCLIENT_TYPE_GATE)

        case LGDATAPACKET_TYPE_GENERAL:
            c.Process(code,msg,c,0)
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

