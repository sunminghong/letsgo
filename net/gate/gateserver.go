/*=============================================================================
#     FileName: gateserver.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-23 15:42:14
#      History:
=============================================================================*/
package gate

import (
    "strconv"
    "time"
    //"net"
    goconf "github.com/hgfischer/goconf"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/log"
)

//Dispatcher
type LGIDispatcher interface {
    Init()
    //Add(client Client,protocols []int)
    Add(gridID int, messageCodes string)
    Dispatch(dp *LGDataPacket) (gridID int, ok bool)
    GroupCode(messageCode int) int
}

type LGGateServer struct {
    *LGServer

    Grids       *LGClientPool

    Dispatcher LGIDispatcher

    Name string

    //makeclient NewGateClientFunc
}


/*
//define client
type LGNewGateClientFunc func(name string, transport *LGTransport,gate *LGGateServer) LGIClient


func LGNewGateServer(
    makePlayerClient NewGateClientFunc, datagram LGIDatagram,
    MakeLGGridClient NewGateClientFunc,
    dispatcher IDispatcher) *LGGateServer {

        //Server:&Server{Clients:NewLGClientMap(),datagram:datagram,boardcast_chan_num:10,read_buffer_size:1024},
    gs := &GateServer{
        Server:NewServer(nil,datagram),
    }

    //gs.Clients = NewLGClientMap()
    gs.makeclient = makePlayerClient
    //gs.datagram = datagram
    //gs.boardcast_chan_num = 10
    //gs.read_buffer_size = 1024

    gs.Dispatcher = dispatcher

    gs.Grids = NewClientPool(defaultMakeLGGridClient,datagram)

    return gs
}
*/

func LGNewGateServer(
    makePlayerClient LGNewClientFunc, datagram LGIDatagram,
    makeGridClient LGNewClientFunc,
    dispatcher LGIDispatcher) *LGGateServer {

        //Server:&Server{Clients:NewLGClientMap(),datagram:datagram,boardcast_chan_num:10,read_buffer_size:1024},
    gs := &LGGateServer{
        LGServer:LGNewServer(makePlayerClient,datagram),
    }

    //gs.Server = NewServer(makePlayerClient,datagram)

    gs.Dispatcher = dispatcher

    gs.Grids = LGNewClientPool(makeGridClient,datagram)

    return gs
}
/*
//该函数主要是接受新的连接和注册用户在transport list
func (gs *LGGateServer) transportHandler(newcid int, connection net.Conn) {
    transport := LGNewTransport(newcid, connection, gs,gs.Datagram)
    name := "c_"+strconv.Itoa(newcid)
    client := gs.makeclient(name,transport)
    gs.Clients.Add(newcid, name, client)

    //创建go的线程 使用Goroutine
    go gs.transportSender(transport, client)
    go gs.transportReader(transport, client)

    LGDebug("has clients:",s.Clients.Len())
}
*/

func (gs *LGGateServer) Start(configfile string) {
    //parse config ini file
    gs.ConnectGrids(configfile)
    gs.StartGate(configfile)
}

func (gs *LGGateServer) ConnectGrids(configfile string) {
    c, err := goconf.ReadConfigFile(configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    //make some connection to game server
    for i:=1; i<50; i++ {
        section := "GateServer" + strconv.Itoa(i)
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

        endian, err := c.GetInt("GateServer","endian")
        if err == nil {
            da := gs.Datagram.Clone(endian)
            gs.ConnectGrid(gname, host, messageCodes,da)
        } else {
            gs.ConnectGrid(gname, host, messageCodes,nil)
        }

    }
}

func (gs *LGGateServer) StartGate(configfile string) {
    c, err := goconf.ReadConfigFile(configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    //start gate service
    gatename, err := c.GetString("GateServer","name")
    if err != nil {
        LGError(err.Error())
        return
    }

    gatehost, err := c.GetString("GateServer","host")
    if err != nil {
        LGError(err.Error())
        return
    }

    maxConnections, err := c.GetInt("GateServer","maxConnections")
    if err != nil {
        LGError(err.Error())
        return
    }

    endian, err := c.GetInt("GateServer","endian")
    if err == nil {
        gs.Datagram.SetEndian(endian)
    }

    gs.Name = gatename

    gs.LGServer.Start(gatehost,maxConnections)
}

func (gs *LGGateServer) ConnectGrid(name string,host string,messageCodes string,datagram LGIDatagram) {

        pool := gs.Grids
        go pool.Start(name, host, datagram)
        time.Sleep(1)

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

