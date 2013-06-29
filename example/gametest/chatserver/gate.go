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
//func newPlayerClient (name string,transport *LGTransport) LGIClient {
//    LGTrace("gateclient is connect:",name)
//
//    c := &LGClient{
//        LGBaseClient:&LGBaseClient{Transport:transport,Name:name},
//        Gate : gateserver,
//    }
//
//    c.Init()
//    return c
//}
//
func newGridClient (name string,transport *LGTransport) LGIClient {
    LGTrace("gridclient is connect:",name)

    c := &LGGateToGridClient{LGBaseClient:&LGBaseClient{Transport:transport,Name:name}}
    c.Gate = gateserver

    c.Register()

    return c
}


var (
    loglevel = flag.Int("loglevel",0,"log level")
    gateconf = flag.String("gateconf","gate.conf","gate server config file")
    gridsconf = flag.String("gridconf","grids.conf","grid server config file")
)

func main() {
    flag.Parse()

    //todo: server endian
    datagram := LGNewDatagram(endian)
    //gateserver := LGNewGateServer(
    //    LGNewClient,datagram,newGridClient,LGNewDispatcher())


    gateserver.InitFromConfig(
        *gateconf,LGNewClient,datagram,newGridClient,LGNewDispatcher())

    LGSetLevel(*loglevel)

    gateserver.Start(gridsconf)

}
