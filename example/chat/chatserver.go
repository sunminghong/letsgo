/*=============================================================================
#     FileName: echoserver.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-07 18:50:33
#      History:
=============================================================================*/
package main

import (
    lnet "github.com/sunminghong/letsgo/net"
    "./lib"
)

// IProtocol  
type Client struct {
    Transport *lnet.Transport
    Name *string
}

func MakeClient (name string,transport *lnet.Transport) lnet.IProtocol {
    name_ := "someone"
    return &Client{transport,&name_}
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*lnet.DataPacket) {
    for _,dp:=range dps {
        md := string(dp.Data)

        if md == "/quit" {
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

func main() {
    datagram := &lib.Datagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(MakeClient,datagram,config)

    serv.Start("",4444)
}
