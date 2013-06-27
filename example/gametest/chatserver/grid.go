/*=============================================================================
#     FileName: gate.go
#         Desc: game grid server
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:09:28
#      History:
=============================================================================*/
package main

import (
    "flag"
    //    "strconv"
    //"time"
    //"net"
    //goconf "github.com/hgfischer/goconf"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    "./grids"
)

const (
    endian = LGLittleEndian
)

var serv *LGServer

var (
    loglevel = flag.Int("loglevel", 0, "log level")
    addr     = flag.String("add", ":12001", "grid server addr")
)
func main() {
    flag.Parse()

    LGSetLevel(*loglevel)

    datagram := LGNewDatagram(endian)
    serv = LGNewServer(grids.NewClient, datagram)

    serv.Start(*addr, 2)
}

