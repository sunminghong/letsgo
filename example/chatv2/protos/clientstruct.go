/*=============================================================================
#     FileName: clientstruct.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-17 18:33:10
#      History:
=============================================================================*/
package protos

import (
    lnet "github.com/sunminghong/letsgo/net"
    "github.com/sunminghong/letsgo/helper"
    "github.com/sunminghong/letsgo/log"
)

var Endian int = helper.BigEndian

// IClient  
type Client struct {
    *lnet.DefaultClient
}

func MakeClient (name string,transport *lnet.Transport) lnet.IClient {
    username := "someone"
    c := &Client{}
    c.Username = &username
    c.Transport = transport
    c.Name = name

    return c
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*lnet.DataPacket) {
    for _, dp := range dps {
        msg := lnet.NewMessageReader(dp.Data,Endian)
        log.Trace("msg.code:",msg.Code,msg.Ver)

        //todo: route don't execute
        Handl(msg.Code,c,msg)
    }
}

func (c *Client) Closed() {
    msg := "system: " + (*c.Username) + " is leave!"
    mw := lnet.NewMessageWriter(c.Transport.Stream.Endian)
    mw.SetCode(2011,0)
    mw.WriteString(msg,0)

    c.Transport.SendBoardcast(mw.ToBytes())
}

func (c *Client) SendMessage(msg lnet.IMessageWriter) {
    c.Transport.SendDP(0,msg.ToBytes())
}

func (c *Client) SendBoardcast(msg lnet.IMessageWriter) {
    c.Transport.SendBoardcast(msg.ToBytes())
}

func NewMessageWriter(c lnet.IClient) *lnet.MessageWriter {
    return lnet.NewMessageWriter(c.GetTransport().Stream.Endian)
}
