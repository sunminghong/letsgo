/*=============================================================================
#     FileName: gate.go
#         Desc: game server
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-08 17:16:13
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

var gateserver *LGGateServer


func newPlayerClient (name string,transport *LGTransport) LGIClient {
    LGTrace("gateclient is connect:",name)

    c := &LGClient{
        LGBaseClient:&LGBaseClient{Transport:transport,Name:name},
        Gate : gateserver,
    }

    c.Init()
    return c
}

func newGridClient (name string,transport *LGTransport) LGIClient {
    LGTrace("gridclient is connect:",name)

    c := &LGGridClient{LGBaseClient:&LGBaseClient{Transport:transport,Name:name}}
    c.Gate = gateserver
    c.Register()

    return c
}

func newGateServer() *LGGateServer {
    datagram := LGNewDatagram(endian)
    gs := LGNewGateServer(
        newPlayerClient,datagram,newGridClient,LGNewDispatcher())

   return gs
}


var (
    loglevel = flag.Int("loglevel",0,"log level")
    gateconf = flag.String("gateconf","gate.conf","gate server config file")
    gridsconf = flag.String("gridconf","grids.conf","grid server config file")
)

func main() {
    flag.Parse()

    gateserver = newGateServer()

    LGSetLevel(*loglevel)

    gateserver.Start(gateconf,gridsconf)

}
