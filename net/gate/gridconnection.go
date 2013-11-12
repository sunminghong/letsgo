/*=============================================================================
#     FileName: gridclient.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-11-06 19:25:47
#      History:
=============================================================================*/


/*

*/
package gate

import (
    "reflect"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
)


//define client
type LGGridProcessHandleFunc func(
    msg LGIMessageReader,c *LGGridConnection,fromCid int)


// LGIConnection  
type LGGridConnection struct {
    *LGBaseConnection

    Process LGGridProcessHandleFunc

    GateId int
    Grid *LGGridServer

    //parent
    Parent interface{}
    parentMethodsMap map[string]reflect.Value
}

func (s *LGGridConnection) SetParent(p interface{},methods ...string) {
    s.Parent = p
    if len(methods) == 0 {
        methods = []string{"ConnectionByGateClosed"}
    }

    methodmap := make(map[string]reflect.Value)
    parent := reflect.ValueOf(s.Parent)
    for _,mname := range methods {
        method := parent.MethodByName(mname)
        if method.IsValid() {
            methodmap[mname] = method
        }
    }
    s.parentMethodsMap = methodmap
}

func (c *LGGridConnection) Closed() {
    if c.GateId > 0 {
        c.Grid.RemoveGate(c.GateId,c.GetTransport().Cid)
    }
}

func (c *LGGridConnection) ConnectionByGateClosed(gateid int, fromCid int) {
    if method,ok := c.parentMethodsMap["ConnectionByGateClosed"]; ok {
        args := []reflect.Value{
            reflect.ValueOf(gateid),
            reflect.ValueOf(fromCid),
        }

        method.Call(args)
        return
    }

    panic("LGpanic:LGGridConnection's Method ConnectionByGateClosed need override write by sub object")
}

//对数据进行拆包
func (c *LGGridConnection) ProcessDPs(dps []*LGDataPacket) {
    stream := c.Transport.Stream
    endianer := stream.Endianer
    for _, dp := range dps {

        switch dp.Type {
        case LGDATAPACKET_TYPE_DELAY:
            LGTrace("msg.code(delay):",int(endianer.Uint16(dp.Data)),len(dp.Data))
            msg := LGNewMessageReader(dp.Data,stream.Endian)
            c.Process(msg,c,dp.FromCid)

        case LGDATAPACKET_TYPE_GENERAL:
            LGTrace("msg.code:",int(endianer.Uint16(dp.Data)),len(dp.Data))
            msg := LGNewMessageReader(dp.Data,stream.Endian)
            c.Process(msg,c,0)

        case LGDATAPACKET_TYPE_CLOSED:
            c.ConnectionByGateClosed(c.GateId,dp.FromCid)

        case LGDATAPACKET_TYPE_GATE_REGISTER:
            gatename,gateid := cmd.UnRegister(dp.Data)
            c.SetType(LGCLIENT_TYPE_GATE)

            if c.Grid == nil {
                //LGTrace("transport.server type is ",reflect.TypeOf(c.Transport.Server))
                if grid,ok := c.Transport.Server.(*LGGridServer) ;ok {
                    c.Grid = grid
                } else {
                    LGError("gridserver client init error:transport.Server is not GridServer type")
                }
            }
            c.GateId = gateid
            c.Grid.RegisterGate(gatename,gateid,c)

            LGInfo(c.GetTransport().Conn.RemoteAddr()," is register to gate,gateid=",gateid)
        }
    }
}

func (c *LGGridConnection) SendMessage(fromCid int,msg LGIMessageWriter) {
    dp := &LGDataPacket{
        FromCid: fromCid,
        Data: msg.ToBytes(),
    }

    if fromCid == 0 {
        dp.Type = LGDATAPACKET_TYPE_GENERAL
    } else {
        dp.Type = LGDATAPACKET_TYPE_DELAY
    }

    c.Transport.SendDP(dp)
}

func (c *LGGridConnection) SendBytes(ifCompress bool,fromCid int,data []byte) {
    if (fromCid == 0) {
        c.Transport.SendBytes(data)
        return
    }

    dp := &LGDataPacket{
        Type: LGDATAPACKET_TYPE_DELAY_DATAS,
        FromCid: fromCid,
        Data: data,
    }
    if ifCompress {
        dp.Type = LGDATAPACKET_TYPE_DELAY_DATAS_COMPRESS
    }

    LGTrace("sendbytes from gridclient:%s,%v",fromCid,dp.Data)
    c.Transport.SendDP(dp)
}

