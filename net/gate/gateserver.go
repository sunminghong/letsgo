/*=============================================================================
#     FileName: gateserver.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:09:20
#      History:
=============================================================================*/
package gate

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

type LGGateServer struct {
    *LGServer

    Grids       *LGClientPool

    Dispatcher LGIDispatcher

    //makeclient NewGateClientFunc
}
/*
func LGNewGateServer(
    name string,gateid int,
    newPlayerClient LGNewClientFunc, datagram LGIDatagram,
    newGridClient LGNewClientFunc,
    dispatcher LGIDispatcher) *LGGateServer {

    gs := &LGGateServer{
        LGServer:LGNewServer(name,gateid,newPlayerClient,datagram),
    }

    gs.Dispatcher = LGNewDispatcher()

    gs.Grids = LGNewClientPool(newGridClient,datagram)

    return gs
}
*/

func (gs *LGGateServer) NewTransport(
    newcid int, conn net.Conn) *LGTransport {

    LGTrace("gateserver's newtransport is run")
    return LGNewTransport(newcid, conn, gs,gs.Datagram)
}

func (gs *LGGateServer) Start(gateconfigfile *string,gridsconfigfile *string) {
    //parse config ini file
    gs.connectGrids(gridsconfigfile)
    gs.startGate(gateconfigfile)
}

func (gs *LGGateServer) connectGrids(configfile *string) {
    c, err := goconf.ReadConfigFile(*configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    //make some connection to game server
    for i:=1; i<50; i++ {
        section := "GridServer" + strconv.Itoa(i)
        if !c.HasSection(section) {
            continue
        }
        gname, err := c.GetString(section,"name")
        if err != nil {
            //if err.Reason == goconf.SectionNotFound {
            //    break
            //} else {
                LGError(err.Error())
            //    continue
            //}
            break
        }

        host, err := c.GetString(section,"host")
        if err != nil {
            continue
        }

        messageCodes, err := c.GetString(section,"messageCodes")
        if err != nil {
            messageCodes = ""
        }

        endian, err := c.GetInt(section,"endian")
        if err == nil {
            da := gs.Datagram.Clone(endian)
            gs.ConnectGrid(gname, host, &messageCodes,da)
        } else {
            gs.ConnectGrid(gname, host, &messageCodes,nil)
        }

    }
}

func (gs *LGGateServer) startGate(configfile *string) {
    c, err := goconf.ReadConfigFile(*configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    section := "GateServer"
    //start gate service
    gatename, err := c.GetString(section,"name")
    if err != nil {
        LGError(err.Error())
        return
    }

    gatehost, err := c.GetString(section,"host")
    if err != nil {
        LGError(err.Error())
        return
    }

    maxConnections, err := c.GetInt(section,"maxConnections")
    if err != nil {
        LGError(err.Error())
        return
    }

    endian, err := c.GetInt(section,"endian")
    if err == nil {
        gs.Datagram.SetEndian(endian)
    }

    gs.Name = gatename

    //gs.LGServer.Start(gatehost,maxConnections)
    gs.LGServer.Start(gatehost,maxConnections)
}

func (gs *LGGateServer) ConnectGrid(name string,host string,messageCodes *string,datagram LGIDatagram) {

        pool := gs.Grids
        go pool.Start(name, host, datagram)
        time.Sleep(2*time.Second)

        LGTrace("clientpool:",pool.Clients.All())
        //if Pool don't find it ,then that is no success!
        c := pool.Clients.GetByName(name)
        if c == nil {
            LGError(host + " can't connect")
            return
        }

        //add dispatche
        gridID := c.GetTransport().Cid
        gs.Dispatcher.Add(gridID,messageCodes)
}

