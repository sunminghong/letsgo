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
)

//define a struct or class of rec transport connection
type DataPacket struct {
    Type  byte
    Data  []byte
    Other interface{}
}

type ClientMap struct {
    maplock sync.RWMutex

    maps map[int]IClient
    mapsByName map[string]int
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
    c, ok := tm.maps[cid]
    if ok {
        return c
    }
    return nil
}

func (tm *ClientMap) GetByName(name string) IClient {
    cid, ok := tm.mapsByName[name]
    if ok {
        return tm.maps[cid]
    }
    return nil
}

func (tm *ClientMap) All() map[int]IClient {
    return tm.maps
}

func NewClientMap() *ClientMap { return &ClientMap{maps: make(map[int]IClient),mapsByName: make(map[string]int)} }

type Server struct {
    boardcast_chan_num int
    read_buffer_size   int

    makeclient NewClientFunc
    datagram   IDatagram

    host string
    port int

    //define transport dict/map set
    Clients *ClientMap

    TransportNum int

    boardcastChan chan *DataPacket
}

func NewServer(makeclient NewClientFunc, datagram IDatagram, config map[string]interface{}) *Server {
    s := &Server{Clients: NewClientMap()}

    s.makeclient = makeclient

    s.datagram = datagram

    s.boardcast_chan_num = 10
    s.read_buffer_size = 1024

    return s
}

func (s *Server) Start(host string, port int) {
    Log("Hello Server!")

    addr := host + ":" + strconv.Itoa(port)

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.boardcastChan = make(chan *DataPacket, s.boardcast_chan_num)
    go s.boardcastHandler(s.boardcastChan)

    Log("listen with :", addr)
    netListen, error := net.Listen("tcp", addr)
    if error != nil {
        Log(error)
    } else {
        //defer函数退出时执行
        defer netListen.Close()
        for {
            Log("Waiting for transports")
            connection, error := netListen.Accept()
            if error != nil {
                Log("Transport error: ", error)
            } else {
                newcid := s.allocTransportid()
                go s.transportHandler(newcid, connection)
            }
        }
    }
}

func (s *Server) removeClient(cid int) {
    s.Clients.Remove(cid)
}

func (s *Server) allocTransportid() int {
    s.TransportNum += 1
    return s.TransportNum
}

//该函数主要是接受新的连接和注册用户在transport list
func (s *Server) transportHandler(newcid int, connection net.Conn) {
    transport := NewTransport(newcid, connection, s)
    name := "c_"+strconv.Itoa(newcid)
    client := s.makeclient(name,transport)
    s.Clients.Add(newcid, name, client)

    //创建go的线程 使用Goroutine
    go s.transportSender(transport)
    go s.transportReader(transport, client)

}

func (s *Server) transportReader(transport *Transport, client IClient) {
    buffer := make([]byte, s.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {
            client.Closed()
            transport.Closed()
            s.removeClient(transport.Cid)
            Log(err)
            break
        }

        Log("read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        Log("transport.Buff", transport.Stream.Bytes())
        n, dps := s.datagram.Fetch(transport)
        Log("fetch message number", n)
        if n > 0 {
            client.ProcessDPs(dps)
        }
    }
    Log("TransportReader stopped for ", transport.Cid)
}

func (s *Server) transportSender(transport *Transport) {
    for {
        select {
        case dp := <-transport.Outgoing:
            Log(dp.Type, dp.Data)
            buf := s.datagram.Pack(dp)
            transport.Conn.Write(buf)
        case <-transport.Quit:
            Log("Transport ", transport.Cid, " quitting")
            transport.Conn.Close()
            break
        }
    }
}

func (s *Server) boardcastHandler(boardcastChan <-chan *DataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        Log("boardcastHandler: chan Waiting for input")
        dp := <-boardcastChan
        //buf := s.datagram.pack(dp)

        sendCid, ok := dp.Other.(int)
        if !ok {
            sendCid = 0
        }

        for Cid, c := range s.Clients.All() {
            if sendCid == Cid {
                continue
            }
            c.GetTransport().Outgoing <- dp
        }
        Log("boardcastHandler: Handle end!")
    }
}

//send boardcast message data for other object
func (s *Server) SendBoardcast(transport *Transport, data []byte) {
    dp := &DataPacket{Type: DATAPACKET_TYPE_BOARDCAST, Data: data, Other: transport.Cid}
    s.boardcastChan <- dp
}

//send message data for other object
func (s *Server) SendDP(transport *Transport, dataType byte, data []byte) {
    dp := &DataPacket{Type: dataType, Data: data}
    transport.Outgoing <- dp

}
