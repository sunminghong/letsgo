/*=============================================================================
#     FileName: defaultclient.go
#         Desc: default dispatcher
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-06 14:56:05
#      History:
=============================================================================*/
package net

import (
    . "github.com/sunminghong/letsgo/log"
)

// Client  
type LGDefaultClient struct {
    *LGBaseClient

    Process LGProcessHandleFunc
}

/*
 need write blow func
func LGProccessHandle(code int,msg *MessageReader,c LGIClient,fromCid int) {
    fmt.Println("message is request")
}

func LGMakeDefaultClient (name string,transport *LGTransport) LGIClient {
    c := &LGBaseClient{
        LGBaseClient:&LGBaseClient{transport,name,LGCLIENT_TYPE_GENERAL},
    }
    c.Process = ProcessHandle
}
*/

//对数据进行拆包
func (c *LGDefaultClient) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        LGTrace("msg.code:",code,len(dp.Data))

        msgReader := LGNewMessageReader(dp.Data,c.Transport.Stream.Endian)

        c.Process(code, msgReader,c,0)
    }
}

