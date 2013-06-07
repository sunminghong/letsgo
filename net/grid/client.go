/*=============================================================================
#     FileName: defaultgridclient.go
#         Desc: client of default grid server receive (process player or gate connection on common)
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-07 10:40:26
#      History:
=============================================================================*/
package grid

import (
    "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)


// IClient  
type GridClient struct {
    *DefaultClient
}

/*
func MakeGridClient (name string,transport *Transport,process ProcessHandleFunc) IClient {
    log.Trace("gridclient is connect:",name)
    return c
}
*/

//对数据进行拆包
func (c *GridClient) ProcessDPs(dps []*DataPacket) {
    for _, dp := range dps {
        code := int(c.Transport.Stream.Endianer.Uint16(dp.Data))
        log.Trace("msg.code:",code,len(dp.Data))

        msg := NewMessageReader(dp.Data,c.Transport.Stream.Endian)

        switch dp.Type {
        case DATAPACKET_TYPE_DELAY:
            c.Process(code,msg,c,dp.FromCid)

        case DATAPACKET_TYPE_GATECONNECT:
            c.SetType(CLIENT_TYPE_GATE)

        case DATAPACKET_TYPE_GENERAL:
            c.Process(code,msg,c,0)
        }
    }
}

func (c *GridClient) SendMessage(fromcid int,msg IMessageWriter) {
    dp := &DataPacket{
        FromCid: fromcid,
        Data: msg.ToBytes(),
    }
    if fromcid == 0 {
        dp.Type = DATAPACKET_TYPE_GENERAL
    } else {
        dp.Type = DATAPACKET_TYPE_DELAY
    }

    c.Transport.SendDP(dp)
}

