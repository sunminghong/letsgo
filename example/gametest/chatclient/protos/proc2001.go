/*=============================================================================
#     FileName: proc2001.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-08-07 18:17:15
#      History:
=============================================================================*/
package protos

import (
    "fmt"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/log"
)


func Process2001(msgReader LGIMessageReader,c LGIClient,fromCid int) {
    LGTrace("process 2001 is called")

    subcode := msgReader.ReadUint()
    msg := msgReader.ReadString()

    fmt.Println()
    fmt.Println(msg)

    if subcode == 0 {
        //close this connection
        c.Close()
    }
}

