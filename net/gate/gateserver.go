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
    . "github.com/sunminghong/letsgo/helper"
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

func (gs *LGGateServer) InitFromConfig (
    configfile string,
    newPlayerClient LGNewClientFunc, datagram LGIDatagram,
    newGridClient LGNewClientFunc, dispatcher LGIDispatcher) {

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

    endian, err := c.GetInt(section,"endian")
    if err == nil {
        datagram.SetEndian(endian)
    } else {
        datagram.SetEndian(LGLittleEndian)
    }

    gs.Init( name,serverid,host,maxConnections,
    newPlayerClient,datagram,newGridClient,dispatcher)

}

func (gs *LGGateServer) Init(
    name string,gridid int, host string, maxConnections int,
    newPlayerClient LGNewClientFunc, datagram LGIDatagram,
    newGridClient LGNewClientFunc, dispatcher LGIDispatcher) {

    gs.LGServer = LGNewServer(name,gridid,host,maxConnections,newPlayerClient,datagram)

    gs.Grids = LGNewClientPool(newGridClient,datagram)

    gs.Dispatcher = LGNewDispatcher()

    gs.SetParent(gs)
}

func (gs *LGGateServer) NewTransport(
    newcid int, conn net.Conn) *LGTransport {

    LGTrace("gateserver's newtransport is run")
    return LGNewTransport(newcid, conn, gs,gs.Datagram)
}

func (gs *LGGateServer) Start(gateconfigfile *string,gridsconfigfile *string) {
    //parse config ini file
    gs.connectGrids(gridsconfigfile)
    gs.LGServer.Start()
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

func (gs *LGGateServer) ConnectGrid(
    name string,host string,messageCodes *string,datagram LGIDatagram) {

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

