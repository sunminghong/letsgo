/*=============================================================================
#     FileName: proc1011.go
#         Desc: server base
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:40:02
#      History:
=============================================================================*/
package protos

import (
    //"fmt"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    "strconv"
)

func Process1011(msgReader LGIMessageReader,c LGIClient,fromCid int,session *Session) {
    LGTrace("process 1011 is called")

    md := msgReader.ReadString()

    switch {
    case md == "/quit":
        c.Close()
        return
    case len(md)>8 && md[:8] == "/setmax=":
        _max := md[8:]
        max, err := strconv.Atoi(_max)
        if err != nil {
            LGWarn("setmax is error:",err)
            return
        }
        c.GetTransport().Server.SetMaxConnections(max)
        return
    }

    var msg string
    if session == nil {
        //this user is not register
        session = &Session{&md}

        cid := c.GetTransport().Cid
        uid := (fromCid << 5) + cid
        Uidmap.SaveUid(fromCid,cid,uid)
        SetSession(uid,session)

        LGTrace("1011p():fromcid,cid,uid",fromCid,cid,uid)

        msg = "system: welcome to " + md + "!"
    } else {
        msg = (*session.Username) + "> " + md
    }

    LGDebug("1011 write out:", msg)

    mw := newMessageWriter(c)
    mw.SetCode(2011, 0)
    mw.WriteString(msg, 0)

    c.SendBroadcast(fromCid,mw)
}

