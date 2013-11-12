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
    . "github.com/sunminghong/letsgo/net"
)

func newMessageWriter(c LGIConnection) LGIMessageWriter {
    return LGNewMessageWriter(c.GetTransport().Stream.Endian)
}

func init() {

}
