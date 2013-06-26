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
    p "./protos"
)

const (
    endian = LGLittleEndian
)

var serv *LGServer


type Client struct {
    *grid.LGClient
}

func (c *Client) Closed() {
    LGTrace("a grid client closed")

    //msg := "system: " + (*c.Username) + " is leave!"
    //mw := lnet.NewMessageWriter(c.Transport.Stream.Endian)
    //mw.SetCode(2011,0)
    //mw.WriteString(msg,0)

    //c.Transport.SendBroadcast(mw.ToBytes())
}
func newClient(name string,transport *LGTransport) LGIClient {
    c := &LGBaseClient{
        LGBaseClient:&LGBaseClient{transport,name,LGCLIENT_TYPE_GENERAL},
    }
    c.Process = p.processHandleFunc
}


var (
    loglevel = flag.Int("loglevel", 0, "log level")
    addr     = flag.String("add", ":4444", "grid server addr")
)
func main() {
    flag.Parse()

    LGSetLevel(*loglevel)

    datagram := LGNewDatagram(endian)
    serv = NewLGServer(p.newClient, datagram)

    serv.Start(addr, 2)
}

