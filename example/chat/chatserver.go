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

// LGIConnection  
type Connection struct {
    Transport *lnet.Transport
    Name string
    Username *string
}

func LGMakeConnection (name string,transport *lnet.Transport) lnet.LGIConnection {
    username := "someone"
    return &Connection{transport,name,&username}
}

//对数据进行拆包
func (c *Connection) ProcessDPs(dps []*lnet.LGDataPacket) {
    for _,dp:=range dps {
        md := string(dp.Data)

        if md == "/quit" {
            c.Close()
            return
        }

        var msg string
        if *c.Username == "someone" {
            c.Username = &md

            msg = "system: welcome to " + md + "!"
        } else {
            msg = (*c.Username) + "> "+ md
        }
        c.Transport.SendBroadcast([]byte(msg))
    }
}

//对数据进行拆包
func (c *Connection) GetTransport() *lnet.Transport {
    return c.Transport
}

func (c *Connection) GetName() string {
    return c.Name
}

func (c *Connection) Close() {
    c.Transport.Close()
}

func (c *Connection) Closed() {
    msg := "system: " + (*c.Username) + " is leave!"
    c.Transport.SendBroadcast([]byte(msg))
}

func (c *Connection) SendMessage(msg lnet.LGIMessageWriter) {
    c.Transport.SendDP(0,msg.ToBytes())
}

func (c *Connection) SendBroadcast(msg lnet.LGIMessageWriter) {
    c.Transport.SendBroadcast(msg.ToBytes())
}

func main() {
    datagram := &lib.Datagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(MakeConnection,datagram,config)

    serv.Start("",4444)
}
