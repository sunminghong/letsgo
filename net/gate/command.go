/*=============================================================================
#     FileName: command
#         Desc: gate command
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-28 18:20:23
#      History:
=============================================================================*/
package gate

import (
    //. "github.com/sunminghong/letsgo/net"
    "strings"
    "strconv"
)

type LGCmd struct { }

func (c LGCmd) Pack(line string) []byte {
        return []byte(line)
}

func (c LGCmd) UnPack(data []byte) string {
        return string(data)
}

func (c LGCmd) Register(name string,serverid int) []byte {
    sid :=strconv.Itoa(serverid)
    return c.Pack(name + ":" + sid)
}

func (c LGCmd) UnRegister(data []byte) (name string,serverid int) {
    line := c.UnPack(data)

    lines := strings.Split(line,":")
    name = lines[0]
    serverid,_ = strconv.Atoi(lines[1])

    return
}

var cmd *LGCmd = new(LGCmd)

