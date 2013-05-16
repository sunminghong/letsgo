/*=============================================================================
#     FileName: protocol.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-13 17:50:08
#      History:
=============================================================================*/
package lib

import (
    lnet "github.com/sunminghong/letsgo/net"
)

type Message struct {

}

func Process(c *Client, body []byte) {

    rw := lnet.RWStream(body,lnet.BigEndian)

    msg = rw.ReadString()
        md := string(dp.Data)

        fmt.Println()
        fmt.Println(md)
        fmt.Print("you> ")
    }
}
