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
    "flag"

    lnet "github.com/sunminghong/letsgo/net"
    "github.com/sunminghong/letsgo/log"
    "./protos"
)


var (
    loglevel = flag.Int("loglevel",0,"log level")
)

func main() {
    flag.Parse()

    log.SetLevel(*loglevel)

    datagram := lnet.NewDatagram(protos.Endian)

    serv := lnet.NewServer(protos.MakeClient,datagram)

    serv.Start(":4444",2)
}
