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
package protos

import (
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)


var Handlers map[int]LGProcessHandleFunc= make(map[int]LGProcessHandleFunc)

func proccessHandle(code int,msg LGIMessageReader,c LGIClient,fromCid int) {
    LGTrace("message is request")

    h, ok := Handlers[code]
    if ok {
        h(code,msg,c,fromCid)
    }
}

func newMessageWriter(c LGIClient) LGIMessageWriter {
    return LGNewMessageWriter(c.GetTransport().Stream.Endian)
}

func init() {

}
