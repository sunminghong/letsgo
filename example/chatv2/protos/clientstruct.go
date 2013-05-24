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
)

var Endian int = lnet.BigEndian

// IClient  
type Client struct {
    Transport *lnet.Transport
    Name string
    Username *string
}

func MakeClient (name string,transport *lnet.Transport) lnet.IClient {
    username := "someone"
    return &Client{transport,name,&username}
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*lnet.DataPacket) {
    for _, dp := range dps {
        msg := lnet.NewMessageReader(dp.Data,Endian)
        log.Trace("msg.code:",msg.Code,msg.Ver)

        Handlers[msg.Code](c,msg)
    }
}

//对数据进行拆包
func (c *Client) GetTransport() *lnet.Transport {
    return c.Transport
}

func (c *Client) GetName() string {
    return c.Name
}

func (c *Client) Close() {
    c.Transport.Close()
}

func (c *Client) Closed() {
    msg := "system: " + (*c.Username) + " is leave!"
    mw := lnet.NewMessageWriter(Endian)
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

