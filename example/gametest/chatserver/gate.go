/*=============================================================================
#     FileName: gate.go
#         Desc: game gate server
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
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net/gate"
)

const (
    endian = LGLittleEndian
)

var gateserver *LGGateServer=&LGGateServer{}

//
//func newPlayerConnection (name string,transport *LGTransport) LGIConnection {
//    LGTrace("gateclient is connect:",name)
//
//    c := &LGConnection{
//        LGBaseConnection:&LGBaseConnection{Transport:transport,Name:name},
//        Gate : gateserver,
//    }
//
//    c.Init()
//    return c
//}
//
func newGridConnection (name string,transport *LGTransport) LGIConnection {
    LGTrace("gridclient is connect:",name)

    c := &LGGateToGridConnection{LGBaseConnection:&LGBaseConnection{Transport:transport,Name:name}}
    c.Gate = gateserver

    c.Register()

    return c
}


var (
    loglevel = flag.Int("loglevel",0,"log level")
    gateconf = flag.String("gateconf","gate.conf","gate server config file")
    gridsconf = flag.String("gridconf","gate.conf","grid server config file")
)



func main() {
    flag.Parse()

    //todo: server endian
    datagram := LGNewDatagram(endian)
    //gateserver := LGNewGateServer(
    //    LGNewConnection,datagram,newGridConnection,LGNewDispatcher())



    gateserver.InitFromConfig(
        *gateconf,LGNewConnection,datagram,newGridConnection,LGNewDispatcher())

    LGSetLevel(*loglevel)

    quit := make(chan bool)
    go gateserver.StartConsole(quit)

    gateserver.ConnectGrids(gridsconf)
    go gateserver.Start()

    <-quit
}
