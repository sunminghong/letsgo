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

// Connection  
type LGBaseConnection struct {
    Transport *LGTransport
    Name string
    connectionType int
}

//对数据进行拆包
func (c *LGBaseConnection) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        LGTrace("msg.code:",code,len(dp.Data))
    }
}

func (c *LGBaseConnection) GetTransport() *LGTransport {
    return c.Transport
}

func (c *LGBaseConnection) GetName() string {
    return c.Name
}

func (c *LGBaseConnection) GetType() int{
    return c.connectionType
}

func (c *LGBaseConnection) SetType(t int) {
    c.connectionType = t
}

func (c *LGBaseConnection) Close() {
    c.Transport.Close()
}

func (c *LGBaseConnection) Closed() {
    panic("LGpanic:Closed need override write by sub object")
}

func (c *LGBaseConnection) SendMessage(fromcid int,msg LGIMessageWriter) {
    LGTrace("sendmessage:fromcid",fromcid)
    dp := &LGDataPacket{
        Type: LGDATAPACKET_TYPE_GENERAL,
        FromCid: fromcid,
        Data: msg.ToBytes(),
    }

    c.Transport.SendDP(dp)
}

func (c *LGBaseConnection) SendBroadcast(fromcid int,msg LGIMessageWriter) {
    LGTrace("broadcast:fromcid",fromcid)
    dp := &LGDataPacket{
        Type: LGDATAPACKET_TYPE_BROADCAST,
        Data: msg.ToBytes(),
        FromCid: fromcid,
    }

    c.Transport.SendBroadcast(dp)
}

