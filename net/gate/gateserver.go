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
    goconf "github.com/hgfischer/goconf"
    //"github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/net"
    "github.com/sunminghong/letsgo/log"
)

//Dispatcher
type IDispatcher interface {
    Init()
    //Add(client Client,protocols []int)
    Add(gridID int, messageCodes string)
    Dispatch(dp *DataPacket) (gridID int, ok bool)
    GroupCode(messageCode int) int
}

type GateServer struct {
    *Server

    Grids       *ClientPool

    Dispatcher IDispatcher

    //makeclient NewGateClientFunc
}


/*
//define client
type NewGateClientFunc func(name string, transport *Transport,gate *GateServer) IClient


func NewGateServer(
    makePlayerClient NewGateClientFunc, datagram IDatagram,
    MakeGridClient NewGateClientFunc,
    dispatcher IDispatcher) *GateServer {

        //Server:&Server{Clients:NewClientMap(),datagram:datagram,boardcast_chan_num:10,read_buffer_size:1024},
    gs := &GateServer{
        Server:NewServer(nil,datagram),
    }

    //gs.Clients = NewClientMap()
    gs.makeclient = makePlayerClient
    //gs.datagram = datagram
    //gs.boardcast_chan_num = 10
    //gs.read_buffer_size = 1024

    gs.Dispatcher = dispatcher

    gs.Grids = NewClientPool(defaultMakeGridClient,datagram)

    return gs
}
*/

func NewGateServer(
    makePlayerClient NewClientFunc, datagram IDatagram,
    makeGridClient NewClientFunc,
    dispatcher IDispatcher) *GateServer {

        //Server:&Server{Clients:NewClientMap(),datagram:datagram,boardcast_chan_num:10,read_buffer_size:1024},
    gs := &GateServer{
        Server:NewServer(makePlayerClient,datagram),
    }

    //gs.Server = NewServer(makePlayerClient,datagram)

    gs.Dispatcher = dispatcher

    gs.Grids = NewClientPool(makeGridClient,datagram)

    return gs
}

//该函数主要是接受新的连接和注册用户在transport list
func (gs *GateServer) transportHandler(newcid int, connection net.Conn) {
    transport := NewTransport(newcid, connection, gs,gs.Datagram)
    name := "c_"+strconv.Itoa(newcid)
    client := gs.makeclient(name,transport)
    gs.Clients.Add(newcid, name, client)

    //创建go的线程 使用Goroutine
    go gs.transportSender(transport, client)
    go gs.transportReader(transport, client)

    log.Debug("has clients:",s.Clients.Len())
}

func (gs *GateServer) Start(configfile string) {
    //parse config ini file
    gs.ConnectGrids(configfile)
    gs.StartGate(configfile)
}

func (gs *GateServer) ConnectGrids(configfile string) {
    c, err := goconf.ReadConfigFile(configfile)
    if err != nil {
        log.Error(err.Error())
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
                log.Error(err.Error())
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
            da := gs.datagram.Clone(endian)
            gs.ConnectGrid(gname, host, messageCodes,da)
        } else {
            gs.ConnectGrid(gname, host, messageCodes,nil)
        }

    }
}

func (gs *GateServer) StartGate(configfile string) {
    c, err := goconf.ReadConfigFile(configfile)
    if err != nil {
        log.Error(err.Error())
        return
    }

    //start gate service
    gatename, err := c.GetString("GateServer","name")
    if err != nil {
        log.Error(err.Error())
        return
    }

    gatehost, err := c.GetString("GateServer","host")
    if err != nil {
        log.Error(err.Error())
        return
    }

    maxConnections, err := c.GetInt("GateServer","maxConnections")
    if err != nil {
        log.Error(err.Error())
        return
    }

    endian, err := c.GetInt("GateServer","endian")
    if err == nil {
        gs.datagram.SetEndian(endian)
    }

    gs.Name = gatename
    gs.Gate.Start(gatehost,maxConnections)
}

func (gs *GateServer) ConnectGrid(name string,host string,messageCodes string,datagram Datagram) {

        pool := gs.Grids
        go pool.Start(name, host, datagram)
        time.Sleep(1)

        //if Pool don't find it ,then that is no success!
        c := pool.Clients.GetByName(name)
        if c == nil {
            log.Error(host + " can't connect")
            return
        }

        //add dispatche
        gridID := c.GetTransport().Cid
        gs.Dispatcher.Add(gridID,messageCodes)
}

