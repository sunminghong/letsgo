/*=============================================================================
#     FileName: servermain.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-06 19:15:13
#      History:
=============================================================================*/
package main

import (
    "net"
    lnet "github.com/sunminghong/letsgo/net"
)


func main() {
    datagram := &TextDatagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(NewClient,datagram,config)
    serv.Start("",4444)
}

type TextDatagram struct {

}


//对数据进行拆包
func (d *TextDatagram) Fetch(c *lnet.Client) {

    ilen := len(c.Buff)
    if ilen == 0 {
        return
    }
    lnet.Log("Fetch",c.Buff)
    msg := &lnet.DataPacket{Data: c.Buff}

    //send to channel for consume
    c.ProcessMsg(msg)
    c.Buff = make([]byte,1024)
}

//对数据进行封包
func (d *TextDatagram) Pack(dp *lnet.DataPacket) []byte {
    return dp.Data
}


func(c *lnet.Client) ProcessMsg(msg *lnet.DataPacket) {
    lnet.Log("processmsg:",len(msg.Data))
    c.Outgoing <- msg
}


// new Client object
func NewClient(newclientid int, conn net.Conn, server *lnet.Server) *lnet.Client {
    c := &lnet.Client{
            Cid:      newclientid,
            Conn:     conn,
            Server:   server,
            Outgoing: make(chan *lnet.DataPacket, 10),
            Quit:     make(chan bool),
    }

    c.InitBuff()

    return c
}
