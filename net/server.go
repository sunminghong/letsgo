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
    "reflect"
    "net"
    "time"
    "math/rand"
    "strconv"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

type LGServer struct {
    Name string
    Serverid int
    Addr string

    broadcast_chan_num int
    read_buffer_size   int

    maxConnections int
    makeclient LGNewConnectionFunc
    Datagram   LGIDatagram

    host string
    port int

    //define transport dict/map set
    Connections *LGConnectionMap

    TransportNum int

    broadcastChan chan *LGDataPacket

    exitChan chan bool
    stop bool

    idassign *LGIDAssign

    //parent
    Parent interface{}
    parentMethodsMap map[string]reflect.Value
}

func LGNewServer(
    name string,serverid int,addr string, maxConnections int,
    makeclient LGNewConnectionFunc, datagram LGIDatagram) *LGServer {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    serverid += r.Intn(99) * 100
    s := &LGServer{
        Name:name,
        Serverid:serverid,
        Addr : addr,
        maxConnections : maxConnections,
        Connections: LGNewConnectionMap(),
        stop: false,
        exitChan: make(chan bool),
    }

    s.makeclient = makeclient

    s.Datagram = datagram

    s.broadcast_chan_num = 10
    s.read_buffer_size = 1024

    s.idassign = LGNewIDAssign(1<<16)

    return s
}

func (s *LGServer) SetParent(p interface{},methods ...string) {
    s.Parent = p
    if len(methods) == 0 {
        methods = []string{"NewTransport"}
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

func (s *LGServer) Start() {
    LGInfo(s.Name +" is starting...")

    s.stop=false
    //todo: maxConnections don't proccess
    //addr := host + ":" + strconv.Itoa(port)

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.broadcastChan = make(chan *LGDataPacket, s.broadcast_chan_num)
    go s.broadcastHandler(s.broadcastChan)

    netListen, error := net.Listen("tcp", s.Addr)
    if error != nil {
        LGError(error)
    } else {
        LGInfo("listen with :", s.Addr)
        LGInfo(s.Name +" is started !!!")

        //defer函数退出时执行
        defer netListen.Close()
        for {
            LGTrace("Waiting for connection")
            connection, error := netListen.Accept()
            if s.stop {
                continue
            }

            if error != nil {
                LGError("Transport error: ", error)
            } else {
                LGDebug("%v is connection!",connection.RemoteAddr())

                newcid := s.AllocTransportid()
                if newcid == 0 {
                    LGWarn("connection num is more than ",s.maxConnections)
                } else {
                    newcid = LGGenerateID(newcid)
                    LGTrace("///////////////////////////////////////////////newcid:",newcid)
                    go s.transportHandler(newcid, connection)
                }
            }
        }
    }
}

func (s *LGServer) SetMaxConnections(max int) {
    s.maxConnections = max
}

func (s *LGServer) RemoveConnection(cid int) {
    //if method,ok := s.parentMethodsMap["RemoveConnection"]; ok {
    //    args := []reflect.Value{
    //        reflect.ValueOf(cid),
    //    }

    //    method.Call(args)
    //    return
    //}

    s.Connections.Remove(cid)

    //release id assign
    cid,_= LGParseID(cid)
    s.idassign.Free(cid)
}

func (s *LGServer) AllocTransportid() int {
    if (s.Connections.Len() >= s.maxConnections) {
        return 0
    }

    return s.idassign.GetFree()
}


//for override write by sub struct
func (s *LGServer) NewTransport(newcid int, conn net.Conn) *LGTransport {
    if method,ok := s.parentMethodsMap["NewTransport"]; ok {
        args := []reflect.Value{
            reflect.ValueOf(newcid),
            reflect.ValueOf(conn),
        }

        trans := method.Call(args)[0].Interface().(*LGTransport)
        return trans
    }

    return LGNewTransport(newcid, conn, s,s.Datagram)
}

//该函数主要是接受新的连接和注册用户在transport list
func (s *LGServer) transportHandler(newcid int, connection net.Conn) {
    transport := s.NewTransport(newcid, connection)
    name := "c_"+strconv.Itoa(newcid)
    client := s.makeclient(name,transport)
    s.Connections.Add(newcid, name, client)

    //创建go的线程 使用Goroutine
    go s.transportSender(transport, client)
    go s.transportReader(transport, client)

    LGDebug("has clients:",s.Connections.Len())
}

func (s *LGServer) transportReader(transport *LGTransport, client LGIConnection) {
    buffer := make([]byte, s.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {

            client.Closed()
            transport.Closed()
            transport.Conn.Close()
            s.RemoveConnection(transport.Cid)
            //LGError(err)
            break
        }

        //LGTrace("read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        LGTrace("server transportReader.Buff", transport.Stream.Len())
        n, dps := transport.Fetch()
        //LGTrace("fetch message number", n)
        if n > 0 {
            client.ProcessDPs(dps)
        }
    }
    //LGTrace("TransportReader stopped for ", transport.Cid)
}

func (s *LGServer) transportSender(transport *LGTransport, client LGIConnection) {
    for {
        select {
        case data := <-transport.OutgoingBytes:
            LGTrace("server transportSender OutgoingBytes:",len(data))
            //buf := s.Datagram.Pack(dp)
            transport.Conn.Write(data)

        case dp := <-transport.Outgoing:
            //LGTrace("transportSender Outgoing:type=%d,len=%d,% X",dp.Type, len(dp.Data),dp.Data)
            LGTrace("transportSender Outgoing:type=%d,len=%d",dp.Type, len(dp.Data))
            transport.PackWrite(dp)

        case <-transport.Quit:
            LGDebug("Transport ", transport.Cid, " quitting")

            transport.Closed()
            transport.Conn.Close()
            s.RemoveConnection(transport.Cid)
            break
        }
    }
}

/*
func (s *LGServer) broadcastHandler(broadcastChan <-chan *LGDataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        LGTrace("broadcastHandler: chan Waiting for input")
        dp := <-broadcastChan
        data := s.Datagram.Pack(dp)

        //fromCid := dp.FromCid
        data0 := s.Datagram.Pack(&LGDataPacket{
            Type: LGDATAPACKET_TYPE_GENERAL,
            FromCid: 0,
            Data: dp.Data,
        })
        for _, c := range s.Connections.All() {
            LGTrace("broadcastHandler: client.type",c.GetType())
            //if fromCid == Cid {
            //    continue
            //}
            if c.GetType() == LGCLIENT_TYPE_GATE {
                c.GetTransport().OutgoingBytes <- data
            } else {
                c.GetTransport().OutgoingBytes <- data0
            }
        }
        LGTrace("broadcastHandler: Handle end!")
    }
}
*/

func (s *LGServer) broadcastHandler(broadcastChan <-chan *LGDataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        //LGTrace("broadcastHandler: chan Waiting for input")
        dp := <-broadcastChan

        //fromCid := dp.FromCid
        dp0 := &LGDataPacket{
            Type: LGDATAPACKET_TYPE_GENERAL,
            FromCid: 0,
            Data: dp.Data,
        }
        for _, c := range s.Connections.All() {
            //LGTrace("broadcastHandler: client.type",c.GetType())
            //if fromCid == Cid {
            //    continue
            //}
            if c.GetType() == LGCLIENT_TYPE_GATE {
                c.GetTransport().Outgoing <- dp
            } else {
                c.GetTransport().Outgoing <- dp0
            }
        }
        //LGTrace("broadcastHandler: Handle end!")
    }
}

//send broadcast message data for other object
func (s *LGServer) SendBroadcast(dp *LGDataPacket) {
    s.broadcastChan <- dp
}

func (s *LGServer) Stop() {
    s.stop = true
}

