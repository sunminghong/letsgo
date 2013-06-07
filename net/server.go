/*=============================================================================
#     FileName: server.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-03 14:27:57
#      History:
=============================================================================*/
package net

import (
    "net"
    "strconv"
    "sync"
    "github.com/sunminghong/letsgo/log"
)


type ClientMap struct {
    maplock *sync.RWMutex

    maps map[int]IClient
    mapsByName map[string]int
}

func NewClientMap() *ClientMap {
    return &ClientMap{
        maplock: new(sync.RWMutex),
        maps: make(map[int]IClient),
        mapsByName: make(map[string]int),
    }
}

func (tm *ClientMap) Add(cid int,name string, client IClient) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    tm.maps[cid] = client
    if len(name) > 0 {
        tm.mapsByName[name] = cid
    }
}

func (tm *ClientMap) Remove(cid int) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    _, ok := tm.maps[cid]
    if ok {
        name := tm.maps[cid].GetName()
        if len(name)>0 {
            _,ok :=tm.mapsByName[name]
            if ok {
                delete(tm.mapsByName,name)
            }
        }
        delete(tm.maps, cid)
    }
}

func (tm *ClientMap) RemoveByName(name string) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    cid,ok :=tm.mapsByName[name]
    if ok {
        _, ok := tm.maps[cid]
        if ok {
            delete(tm.maps, cid)
        }
        delete(tm.mapsByName,name)
    }
}

func (tm *ClientMap) Get(cid int) IClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    c, ok := tm.maps[cid]
    if ok {
        return c
    }
    return nil
}

func (tm *ClientMap) GetByName(name string) IClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    cid, ok := tm.mapsByName[name]
    if ok {
        return tm.maps[cid]
    }
    return nil
}

func (tm *ClientMap) All() map[int]IClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return tm.maps
}

func (tm *ClientMap) Len() int {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return len(tm.maps)
}

type Server struct {
    boardcast_chan_num int
    read_buffer_size   int

    maxConnections int
    makeclient NewClientFunc
    datagram   IDatagram

    host string
    port int

    //define transport dict/map set
    Clients *ClientMap

    TransportNum int

    boardcastChan chan *DataPacket
}

func NewServer(makeclient NewClientFunc, datagram IDatagram) *Server {
    s := &Server{Clients: NewClientMap()}

    s.makeclient = makeclient

    s.datagram = datagram

    s.boardcast_chan_num = 10
    s.read_buffer_size = 1024

    return s
}

func (s *Server) Start(addr string, maxConnections int) {
    log.Info("Hello Server!")

    s.maxConnections = maxConnections
    //todo: maxConnections don't proccess
    //addr := host + ":" + strconv.Itoa(port)

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.boardcastChan = make(chan *DataPacket, s.boardcast_chan_num)
    go s.boardcastHandler(s.boardcastChan)

    log.Info("listen with :", addr)
    netListen, error := net.Listen("tcp", addr)
    if error != nil {
        log.Error(error)
    } else {
        //defer函数退出时执行
        defer netListen.Close()
        for {
            log.Trace("Waiting for connection")
            connection, error := netListen.Accept()
            if error != nil {
                log.Error("Transport error: ", error)
            } else {
                newcid := s.allocTransportid()
                if newcid == 0 {
                    log.Warn("connection num is more than ",s.maxConnections)
                } else {
                    go s.transportHandler(newcid, connection)
                }
            }
        }
    }
}

func (s *Server) SetMaxConnections(max int) {
    s.maxConnections = max
}

func (s *Server) removeClient(cid int) {
    s.Clients.Remove(cid)
}

func (s *Server) allocTransportid() int {
    if (s.Clients.Len() >= s.maxConnections) {
        return 0
    }
    s.TransportNum += 1
    return s.TransportNum
}

//该函数主要是接受新的连接和注册用户在transport list
func (s *Server) transportHandler(newcid int, connection net.Conn) {
    transport := NewTransport(newcid, connection, s,s.datagram)
    name := "c_"+strconv.Itoa(newcid)
    client := s.makeclient(name,transport)
    s.Clients.Add(newcid, name, client)

    //创建go的线程 使用Goroutine
    go s.transportSender(transport, client)
    go s.transportReader(transport, client)

    log.Debug("has clients:",s.Clients.Len())
}

func (s *Server) transportReader(transport *Transport, client IClient) {
    buffer := make([]byte, s.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {

            client.Closed()
            transport.Closed()
            transport.Conn.Close()
            s.removeClient(transport.Cid)
            log.Error(err)
            break
        }

        log.Trace("read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        log.Trace("transport.Buff", transport.Stream.Bytes())
        n, dps := s.datagram.Fetch(transport)
        log.Trace("fetch message number", n)
        if n > 0 {
            client.ProcessDPs(dps)
        }
    }
    log.Trace("TransportReader stopped for ", transport.Cid)
}

func (s *Server) transportSender(transport *Transport, client IClient) {
    for {
        select {
        case dp := <-transport.outgoing:
            log.Trace("transportSender outgoing:",dp.Type, len(dp.Data))
            //buf := s.datagram.Pack(dp)
            //transport.Conn.Write(buf)

            s.datagram.PackWrite(transport.Conn.Write,dp)
        case <-transport.Quit:
            log.Debug("Transport ", transport.Cid, " quitting")

            transport.Closed()
            transport.Conn.Close()
            s.removeClient(transport.Cid)
            break
        }
    }
}

func (s *Server) boardcastHandler(boardcastChan <-chan *DataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        log.Trace("boardcastHandler: chan Waiting for input")
        dp := <-boardcastChan
        //buf := s.datagram.pack(dp)

        fromCid := dp.FromCid
        for _, c := range s.Clients.All() {
            //if fromCid == Cid {
            //    continue
            //}
            if c.GetType() != CLIENT_TYPE_GATE {
                dp.FromCid = 0
                dp.Type = DATAPACKET_TYPE_GENERAL
            } else {
                dp.FromCid = fromCid
                dp.Type = DATAPACKET_TYPE_BOARDCAST
            }
            c.GetTransport().outgoing <- dp
        }
        log.Trace("boardcastHandler: Handle end!")
    }
}

//send boardcast message data for other object
func (s *Server) SendBoardcast(transport *Transport, dp *DataPacket) {
    s.boardcastChan <- dp
}

