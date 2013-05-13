/*=============================================================================
#     FileName: echoclient.go
#         Desc: echo text server Datagram pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-13 16:34:55
#      History:
=============================================================================*/
package main

import (
    lnet "github.com/sunminghong/letsgo/net"
)

type Client struct {
    Transport *lnet.Transport
    Name *string
}

func MakeClient (transport *lnet.Transport) lnet.IClient {
    name := "someone"
    return &Client{transport,&name}
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*lnet.DataPacket) {
    for _,dp:=range dps {
        md := string(dp.Data)

        if md == "quit" {
            c.Close()
            return
        }

        var msg string
        if *c.Name == "someone" {
            c.Name = &md

            msg = "system: welcome to " + md + "!"
        } else {
            msg = (*c.Name) + "> "+ md
        }
        c.Transport.SendBoardcast([]byte(msg))
    }
}

//对数据进行拆包
func (c *Client) GetTransport() *lnet.Transport {
    return c.Transport
}

func (c *Client) Close() {
    c.Transport.Close()
}

func (c *Client) Closed() {
    msg := "system: " + (*c.Name) + " is leave!"
    c.Transport.SendBoardcast([]byte(msg))
}
