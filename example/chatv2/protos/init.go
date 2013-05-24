/*=============================================================================
#     FileName: init.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-17 17:55:51
#      History:
=============================================================================*/
package protos

import (
    lnet "github.com/sunminghong/letsgo/net"
)

type processHandler func(c *Client,reader *lnet.MessageReader)

var Handlers map[int]processHandler = make(map[int]processHandler)

func Handl(code int,c *Client,reader *lnet.MessageReader) {
    h, ok := Handlers[code]
    if ok {
        h(c,reader)
    }
}

func init() {

}
/*
func Process(c *Client,reader *MessageReader) {

    rw := lnet.RWStream(body,lnet.BigEndian)

    msg = rw.ReadString()
        md := string(dp.Data)

        fmt.Println()
        fmt.Println(md)
        fmt.Print("you> ")
    }
}*/

