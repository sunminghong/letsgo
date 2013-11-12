/*=============================================================================
#     FileName: proc201.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-08-15 09:49:50
#      History:
=============================================================================*/
package protos

import (
    "fmt"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/log"
)


func Process201(msgReader LGIMessageReader,c LGIConnection,fromCid int) {
    LGTrace("process 201 is called")

    subcode := msgReader.ReadUint()
    msg := msgReader.ReadString()

    fmt.Println()
    fmt.Println(msg)

    if subcode == 0 {
        //close this connection
        c.Close()
    }
}

