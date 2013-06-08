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
    "github.com/sunminghong/letsgo/log"
)

// Client  
type DefaultClient struct {
    *BaseClient

    Process ProcessHandleFunc
}

/*
 need write blow func
func ProccessHandle(code int,msg *MessageReader,c IClient,fromCid int) {
    fmt.Println("message is request")
}

func MakeDefaultClient (name string,transport *Transport) IClient {
    c := &BaseClient{
        BaseClient:&BaseClient{transport,name,CLIENT_TYPE_GENERAL},
    }
    c.Process = ProcessHandle
}
*/

//对数据进行拆包
func (c *DefaultClient) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        log.Trace("msg.code:",code,len(dp.Data))

        msgReader := NewMessageReader(dp.Data,c.Transport.Stream.Endian)

        c.Process(code, msgReader,c,0)
    }
}

