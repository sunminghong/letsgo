/*=============================================================================
#     FileName: playerconnection.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-11-06 19:25:58
#      History:
=============================================================================*/

package gate
import (
    . "github.com/sunminghong/letsgo/net"
)

//defined a struct for connection of one player
type LGPlayerConnection struct {
	*LGGridConnection

	FromCid int
	GateId  int

	Endian int
}

func (self *LGPlayerConnection) SendMessage(msg LGIMessageWriter) {
	self.LGGridConnection.SendMessage(self.FromCid, msg)
}

func (self *LGPlayerConnection) SendBroadcast(msg LGIMessageWriter) {
	self.LGGridConnection.SendBroadcast(self.FromCid, msg)
}

func LGGetPlayerConnection(c *LGGridConnection, gateId, fromCid, cid int) *LGPlayerConnection {

	newc := LGGetConnection(c, gateId, cid)
	if newc == nil {
		return nil
	}

	pc := &LGPlayerConnection{
		LGGridConnection: newc,
		FromCid:      fromCid,
		GateId:       gateId,
		Endian:       newc.GetTransport().Stream.Endian,
	}

	return pc
}
