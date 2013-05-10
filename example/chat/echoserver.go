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


func main() {
    datagram := EchoDatagram{ }

    config := make(map[string]interface{})

    serv := lnet.NewServer(NewEchoClient,datagram,config)

    serv.Start("",4444)
}

