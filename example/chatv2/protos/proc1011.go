/*=============================================================================
#     FileName: proc1011.go
#         Desc: server base
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-17 18:43:13
#      History:
=============================================================================*/
package protos

import (
    //"fmt"
    "github.com/sunminghong/letsgo/log"
    lnet "github.com/sunminghong/letsgo/net"
    "strconv"
)

func init() {
    Handlers[1011] = Process1011
}

func LGProcess1011(c *Client, reader *lnet.MessageReader) {
    log.LGTrace("process 1011 is called")

    md := reader.ReadString()

    switch {
    case md == "/quit":
        c.Close()
        return
    case len(md)>8 && md[:8] == "/setmax=":
        _max := md[8:]
        max, err := strconv.Atoi(_max)
        if err != nil {
            log.Warn("setmax is error:",err)
            return
        }
        c.GetTransport().Server.SetMaxConnections(max)
        return
    }

    var msg string
    if *c.Username == "someone" {
        c.Username = &md

        msg = "system: welcome to " + md + "!"
    } else {
        msg = (*c.Username) + "> " + md
    }

    log.Debug("1011 write out:", msg)

    mw := NewMessageWriter(c)
    mw.SetCode(2011, 0)
    mw.WriteString(msg, 0)
    c.SendBroadcast(mw)
}
