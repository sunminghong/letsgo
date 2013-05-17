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
    lnet "github.com/sunminghong/letsgo/net"
)

func init() {
    Handlers[1011] = Process1011
}

func Process1011(c *Client, reader *lnet.MessageReader) {
    lnet.Log("process 1011 is called")

    md := reader.ReadString()

    var msg string
    if *c.Username == "someone" {
        c.Username = &md

        msg = "system: welcome to " + md + "!"
    } else {
        msg = (*c.Username) + "> " + md
    }

    if md == "/quit" {
        c.Close()
        return
    }
    c.GetTransport().SendBoardcast([]byte(msg))

}
