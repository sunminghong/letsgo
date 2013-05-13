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
)
// Idatagram
type Datagram struct {

}


//对数据进行拆包
func (d *Datagram) Fetch(c *lnet.Transport) (n int,msgs []*lnet.DataPacket) {
    msgs = []*lnet.DataPacket{}

    ilen := c.Stream.Len()
    if ilen == 0 {
        return
    }
    lnet.Log("Fetch",c.Stream.Bytes())
    msg := &lnet.DataPacket{Data: c.Stream.Bytes()}
    msgs = append(msgs,msg)
    n += 1

    //send to channel for consume
    c.InitBuff()

    return
}

//对数据进行封包
func (d *Datagram) Pack(dp *lnet.DataPacket) []byte {
    return dp.Data
}

// IClient  
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

func main() {
    datagram := &Datagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(MakeClient,datagram,config)

    serv.Start("",4444)
}

