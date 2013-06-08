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

type ProcessHandleFunc func(
    code int,msg *MessageReader,c IClient,fromCid int)

// Client  
type DefaultClient struct {
    Transport *Transport
    Name string
    connectionType int

    Process ProcessHandleFunc
}

func MakeDefaultClient (name string,transport *Transport,process ProcessHandleFunc) IClient {
    return &DefaultClient{transport,name,CLIENT_TYPE_GENERAL,process}
}

//对数据进行拆包
func (c *DefaultClient) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        log.Trace("msg.code:",code,len(dp.Data))

        msgReader := NewMessageReader(dp.Data,c.Transport.Stream.Endian)

        c.Process(code, msgReader,c,0)
    }
}

func (c *DefaultClient) GetTransport() *Transport {
    return c.Transport
}

func (c *DefaultClient) GetName() string {
    return c.Name
}

func (c *DefaultClient) GetType() int{
    return c.connectionType
}

func (c *DefaultClient) SetType(t int) {
    c.connectionType = t
}

func (c *DefaultClient) Close() {
    c.Transport.Close()
}

func (c *DefaultClient) Closed() {
    log.Trace("this client is closed!")
    //todo: override write by sub object
    panic("Closed need override write by sub object")
}

func (c *DefaultClient) SendMessage(fromcid int,msg IMessageWriter) {
    dp := &DataPacket{
        Type: DATAPACKET_TYPE_GENERAL,
        FromCid: fromcid,
        Data: msg.ToBytes(),
    }

    c.Transport.SendDP(dp)
}

func (c *DefaultClient) SendBoardcast(fromcid int,msg IMessageWriter) {
    dp := &DataPacket{
        Type: DATAPACKET_TYPE_BOARDCAST,
        Data: msg.ToBytes(),
        FromCid: fromcid,
    }

    c.Transport.SendBoardcast(dp)
}

