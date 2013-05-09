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
    datagram := &lnet.EchoDatagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(lnet.NewEchoClient,datagram,config)
    serv.Start("",4444)
}

