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
package grid

import (
    "strconv"
    "time"
    "net"
    goconf "github.com/sunminghong/goconf"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/log"
)

//Dispatcher
type LGIDispatcher interface {
    Init()
    //Add(client Client,protocols []int)
    Add(gridID int, messageCodes *string)
    Dispatch(messageCode int) (gridID int, ok bool)
    GroupCode(messageCode int) int
}

type LGGridServer struct {
    *LGServer
    GateMap map[int][]int

    AllowDirectConnection bool

    Dispatcher LGIDispatcher

    //makeclient NewGridClientFunc
}

func (gs *LGGridServer) InitFromConfig (
    configfile string, newGridClient LGNewClientFunc, datagram LGIDatagram) *LGGridServer {

    c, err := goconf.ReadConfigFile(*configfile)
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

    serverid, err := c.GetString(section,"serverid")
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
        allowDirectConnection = 0
    }

    endian, err := c.GetInt(section,"endian")
    if err == nil {
        datagram.SetEndian(endian)
    } else {
        datagram.SetEndian(LGLittleEndian)
    }

    gs.Init( name,serverid,allwDirectConnection,host,maxConnections,
        newGridClient,datagram)
}

func (gs *LGGridServer) Init(
    name string,gridid int,allowDirectConnection bool,host string, maxConnections int,
    newGridClient LGNewClientFunc, datagram LGIDatagram) *LGGridServer {


    gs := &LGGridServer{
        LGServer:LGNewServer(name,gridid,host,maxConnections,newPlayerClient,datagram),
    }

    gs.GridMap = make(map[int][]int)

    gs.AllowDirectConnection = allowDirectConnection

    gs.SetParent(gs)
}

func (gs *LGGridServer) RegisterGate(gridname string,gridid int,c *LGIClient) {
    if cs,ok := gs.GateMap[gridid]; ok {
        cs = append(cs,gridid)

        gs.GateMap[gridid] = cs
    } else {
        gs.GateMap[gridid] = []int {c.GetTransport().Cid}
    }
}

func (gs *LGGridServer) NewTransport(
    newcid int, conn net.Conn) *LGTransport {

    LGTrace("gridserver's newtransport is run")
    return LGNewTransport(newcid, conn, gs,gs.Datagram)
}

func (gs *LGGridServer) NewTransport(
    newcid int, conn net.Conn) *LGTransport {

    LGTrace("gridserver's newtransport is run")
    return LGNewTransport(newcid, conn, gs,gs.Datagram)
}

func (gs *LGGridServer) BroadcastHandler(broadcastChan <-chan *LGDataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        LGTrace("broadcastHandler: chan Waiting for input")
        dp := <-broadcastChan

        //fromCid := dp.FromCid
        dp0 := &LGDataPacket{
            Type: LGDATAPACKET_TYPE_GENERAL,
            FromCid: 0,
            Data: dp.Data,
        }

        if gs.AllowDirectConnection {
            for _, c := range s.Clients.All() {
                //if fromCid == Cid {
                //    continue
                //}
                if c.GetType() == LGCLIENT_TYPE_GATE {
                    continue
                }
                c.GetTransport().outgoing <- dp0
            }
        }

        //broadcast to grid server
        for _, cs := range gs.GateMap {
            LGTrace("broadcastHandler: client.type",c.GetType())

            cid := cs[0]
            s.Clients.Get(cid).GetTransport().outgoing <- dp
        }
        LGTrace("broadcastHandler: Handle end!")
    }
}

