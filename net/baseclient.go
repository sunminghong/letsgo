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
    "github.com/sunminghong/letsgo/log"
)

// Client  
type BaseClient struct {
    Transport *Transport
    Name string
    connectionType int
}

//对数据进行拆包
func (c *BaseClient) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        log.Trace("msg.code:",code,len(dp.Data))
    }
}

func (c *BaseClient) GetTransport() *Transport {
    return c.Transport
}

func (c *BaseClient) GetName() string {
    return c.Name
}

func (c *BaseClient) GetType() int{
    return c.connectionType
}

func (c *BaseClient) SetType(t int) {
    c.connectionType = t
}

func (c *BaseClient) Close() {
    c.Transport.Close()
}

func (c *BaseClient) Closed() {
    log.Trace("this client is closed!")
    //todo: override write by sub object
    panic("Closed need override write by sub object")
}

func (c *BaseClient) SendMessage(fromcid int,msg IMessageWriter) {
    dp := &DataPacket{
        Type: DATAPACKET_TYPE_GENERAL,
        FromCid: fromcid,
        Data: msg.ToBytes(),
    }

    c.Transport.SendDP(dp)
}

func (c *BaseClient) SendBoardcast(fromcid int,msg IMessageWriter) {
    dp := &DataPacket{
        Type: DATAPACKET_TYPE_BOARDCAST,
        Data: msg.ToBytes(),
        FromCid: fromcid,
    }

    c.Transport.SendBoardcast(dp)
}

