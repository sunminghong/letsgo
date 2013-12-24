/*=============================================================================
#     FileName: gridserver.go
#         Desc: grid server
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-28 19:39:02
#      History:
=============================================================================*/
package gate

import (
    "net"
    iniconfig "github.com/sunminghong/iniconfig"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

type LGGridServer struct {
    *LGServer
    GateMap map[int]LGSliceInt

    AllowDirectConnection bool

    //Dispatcher LGIDispatcher

    //makeclient NewGridConnectionFunc
}

func (gs *LGGridServer) InitFromConfig (
    configfile string, newGridConnection LGNewConnectionFunc,
    datagram LGIDatagram) {

    c, err := iniconfig.ReadConfigFile(configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    section := "Default"

    logconf, err := c.GetString(section,"logConfigFile")
    if err != nil {
        logconf = ""
    }
    LGSetLogger(&logconf)


    //start grid service
    name, err := c.GetString(section,"name")
    if err != nil {
        LGError(err.Error())
        return
    }

    host, err := c.GetString(section,"host")
    if err != nil {
        LGError(err.Error())
        return
    }

    serverid, err := c.GetInt(section,"serverid")
    if err != nil {
        LGError(err.Error())
        return
    }

    maxConnections, err := c.GetInt(section,"maxConnections")
    if err != nil {
        maxConnections = 1000
    }

    allowDirectConnection, err := c.GetBool(section,"allowDirectConnection")
    if err != nil {
        allowDirectConnection = false
    }

    endian, err := c.GetInt(section,"endian")
    if err == nil {
        datagram.SetEndian(endian)
    } else {
        datagram.SetEndian(LGLittleEndian)
    }

    profile, err := c.GetInt(section,"profile")
    if err != nil {
        profile = 0
    }

    if profile > 0 {
		//LGSetMemProfileRate(1)
		//defer func() {
		//    LGGC()
		//    LGDumpHeap()
		//    LGPrintGCSummary()
		//}()
    }

    gs.Init( name,serverid,allowDirectConnection,host,maxConnections,
        newGridConnection,datagram)
}

func (gs *LGGridServer) Init(
    name string,gridid int,allowDirectConnection bool,host string,
    maxConnections int,
    newGridConnection LGNewConnectionFunc, datagram LGIDatagram) {


    gs.LGServer = LGNewServer(
        name,gridid,host,maxConnections,newGridConnection,datagram)

    gs.GateMap = make(map[int]LGSliceInt)

    gs.AllowDirectConnection = allowDirectConnection

    gs.SetParent(gs)
}

func (gs *LGGridServer) RegisterGate(gatename string,gateid int,c LGIConnection) {
    if cs,ok := gs.GateMap[gateid]; ok {
        cs = append(cs,gateid)

        gs.GateMap[gateid] = cs
    } else {
        gs.GateMap[gateid] = LGSliceInt {c.GetTransport().Cid}
    }
}

func (gs *LGGridServer) RemoveGate(gateid ,cid int) {
    LGTrace("RemoveGate")
    if cs,ok := gs.GateMap[gateid]; ok {

        cs.RemoveValue(cid)
        if len(cs) > 0 {
            gs.GateMap[gateid] = cs
        } else {
            delete(gs.GateMap,gateid)
        }
    }
}

func (gs *LGGridServer) NewTransport(newcid int, conn net.Conn) *LGTransport {

    LGTrace("gridserver's newtransport is run")
    return LGNewTransport(newcid, conn, gs,gs.Datagram)
}

func (gs *LGGridServer) BroadcastHandler(broadcastChan <-chan *LGDataPacket) {
    for {
        LGTrace("broadcastHandler: chan Waiting for input")
        dp := <-broadcastChan

        //fromCid := dp.FromCid
        dp0 := &LGDataPacket{
            Type: LGDATAPACKET_TYPE_GENERAL,
            FromCid: 0,
            Data: dp.Data,
        }

        //broadcast to dplayer client of irect connect to this server 
        if gs.AllowDirectConnection {
            for _, c := range gs.Connections.All() {
                //if fromCid == Cid {
                //    continue
                //}
                if c.GetType() == LGCLIENT_TYPE_GATE {
                    continue
                }
                c.GetTransport().Outgoing <- dp0
            }
        }

        //broadcast to grid server
        for _, cs := range gs.GateMap {
            LGTrace("broadcastHandler: gatemap",cs)

            cid := cs[0]
            gs.Connections.Get(cid).GetTransport().Outgoing <- dp
        }
        LGTrace("broadcastHandler: Handle end!")
    }
}

