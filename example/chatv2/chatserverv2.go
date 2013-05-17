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
    "./protos"
)
/*
//对数据进行拆包
func (c *Client) ProcessDPs(dps []*lnet.DataPacket) {
    for _, dp := range dps {
        msg := lnet.NewMessageReader(dp.Data)
        lnet.Log("msg.code:",msg.Code,msg.Ver)

        protos.Handlers[msg.Code](c,msg)
    }

    for _,dp:=range dps {
        md := string(dp.Data)

        if md == "/quit" {
            c.Close()
            return
        }
    }
}
*/

func main() {
    datagram := &lnet.Datagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(protos.MakeClient,datagram,config)

    serv.Start("",4444)
}
