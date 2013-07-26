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
    goconf "github.com/sunminghong/goconf"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

type LGGridServer struct {
    *LGServer
    GateMap map[int][]int

    AllowDirectConnection bool

    Dispatcher LGIDispatcher

    //makeclient NewGridClientFunc
}

func (gs *LGGridServer) InitFromConfig (
    configfile string, newGridClient LGNewClientFunc,
    datagram LGIDatagram) {

    c, err := goconf.ReadConfigFile(configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    section := "Default"
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

    loglevel, err := c.GetInt(section,"logLevel")
    if err != nil {
        loglevel = 0
    }
    LGSetLevel(loglevel)

    gs.Init( name,serverid,allowDirectConnection,host,maxConnections,
        newGridClient,datagram)
}

func (gs *LGGridServer) Init(
    name string,gridid int,allowDirectConnection bool,host string,
    maxConnections int,
    newGridClient LGNewClientFunc, datagram LGIDatagram) {


    gs.LGServer = LGNewServer(
        name,gridid,host,maxConnections,newGridClient,datagram)

    gs.GateMap = make(map[int][]int)

    gs.AllowDirectConnection = allowDirectConnection

    gs.SetParent(gs)
}

func (gs *LGGridServer) RegisterGate(gridname string,gridid int,c LGIClient) {
    if cs,ok := gs.GateMap[gridid]; ok {
        cs = append(cs,gridid)

        gs.GateMap[gridid] = cs
    } else {
        gs.GateMap[gridid] = []int {c.GetTransport().Cid}
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
            for _, c := range gs.Clients.All() {
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
            gs.Clients.Get(cid).GetTransport().Outgoing <- dp
        }
        LGTrace("broadcastHandler: Handle end!")
    }
}

