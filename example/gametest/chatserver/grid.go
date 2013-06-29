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
    . "github.com/sunminghong/letsgo/net/gate"
    "./grids"
)

var serv *LGGridServer

var (
    loglevel = flag.Int("loglevel", 0, "log level")
    conf = flag.String("conf","grid1.conf","grid server config file")
)
func main() {
    flag.Parse()

    LGSetLevel(*loglevel)

    datagram := LGNewDatagram(LGLittleEndian)

    serv = &LGGridServer{}
    serv.InitFromConfig(*conf,grids.NewClient, datagram)

    serv.Start()
}

