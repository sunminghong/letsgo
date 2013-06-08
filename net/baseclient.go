/*=============================================================================
#     FileName: defaultclient.go
#         Desc: default dispatcher
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-06 14:56:05
#      History:
=============================================================================*/
package net

import (
    . "github.com/sunminghong/letsgo/log"
)

// Client  
type LGBaseClient struct {
    Transport *LGTransport
    Name string
    connectionType int
}

//对数据进行拆包
func (c *LGBaseClient) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        LGTrace("msg.code:",code,len(dp.Data))
    }
}

func (c *LGBaseClient) GetTransport() *LGTransport {
    return c.Transport
}

func (c *LGBaseClient) GetName() string {
    return c.Name
}

func (c *LGBaseClient) GetType() int{
    return c.connectionType
}

func (c *LGBaseClient) SetType(t int) {
    c.connectionType = t
}

func (c *LGBaseClient) Close() {
    c.Transport.Close()
}

func (c *LGBaseClient) Closed() {
    LGTrace("this client is closed!")
    //todo: override write by sub object
    panic("Closed need override write by sub object")
}

func (c *LGBaseClient) SendMessage(fromcid int,msg LGIMessageWriter) {
    dp := &LGDataPacket{
        Type: LGDATAPACKET_TYPE_GENERAL,
        FromCid: fromcid,
        Data: msg.ToBytes(),
    }

    c.Transport.SendDP(dp)
}

func (c *LGBaseClient) SendBroadcast(fromcid int,msg LGIMessageWriter) {
    dp := &LGDataPacket{
        Type: LGDATAPACKET_TYPE_BROADCAST,
        Data: msg.ToBytes(),
        FromCid: fromcid,
    }

    c.Transport.SendBoardcast(dp)
}

