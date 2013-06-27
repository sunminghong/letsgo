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
    "strconv"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

type LGServer struct {
    Parent interface{}
    parentMethodsMap map[string]reflect.Value


    broadcast_chan_num int
    read_buffer_size   int

    maxConnections int
    makeclient LGNewClientFunc
    Datagram   LGIDatagram

    host string
    port int

    //define transport dict/map set
    Clients *LGClientMap

    TransportNum int

    broadcastChan chan *LGDataPacket

    idassign *LGIDAssign
}

func LGNewServer(makeclient LGNewClientFunc, datagram LGIDatagram) *LGServer {
    s := &LGServer{Clients: LGNewClientMap()}

    s.makeclient = makeclient

    s.Datagram = datagram

    s.broadcast_chan_num = 10
    s.read_buffer_size = 1024

    s.idassign = LGNewIDAssign(1<<16)

    return s
}

func (s *LGServer) Start(addr string, maxConnections int) {
    LGInfo("Hello Server!")

    s.maxConnections = maxConnections
    //todo: maxConnections don't proccess
    //addr := host + ":" + strconv.Itoa(port)

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.broadcastChan = make(chan *LGDataPacket, s.broadcast_chan_num)
    go s.broadcastHandler(s.broadcastChan)

    LGInfo("listen with :", addr)
    netListen, error := net.Listen("tcp", addr)
    if error != nil {
        LGError(error)
    } else {
        //defer函数退出时执行
        defer netListen.Close()
        for {
            LGTrace("Waiting for connection")
            connection, error := netListen.Accept()
            if error != nil {
                LGError("Transport error: ", error)
            } else {
                LGDebug(connection.RemoteAddr()," is connection!")

                newcid := s.AllocTransportid()
                if newcid == 0 {
                    LGWarn("connection num is more than ",s.maxConnections)
                } else {
                    go s.transportHandler(newcid, connection)
                }
            }
        }
    }
}

func (s *LGServer) SetMaxConnections(max int) {
    s.maxConnections = max
}

func (s *LGServer) removeClient(cid int) {
    s.Clients.Remove(cid)
    s.Freeid()
}

func (s *LGServer) Freeid(cid int) {
    asdf
    //todo: .....
    s.idassign.Free(cid)
}

func (s *LGServer) Allocid() int {
    if (s.Clients.Len() >= s.maxConnections) {
        return 0
    }

    return s.idassign.GetFree()
}

func (s *LGServer) SetParent(p interface{}) {
    s.Parent = p
    methods := []string{"AllocTransportid","NewTransport"}

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

func (s *LGServer) AllocTransportid() int {
    if method,ok := s.parentMethodsMap["AllocTransportid"]; ok {
        id := method.Call(nil)[0].Int()
        return int(id)
    }

    return s.Allocid()
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
    s.Clients.Add(newcid, name, client)

    //创建go的线程 使用Goroutine
    go s.transportSender(transport, client)
    go s.transportReader(transport, client)

    LGDebug("has clients:",s.Clients.Len())
}

func (s *LGServer) transportReader(transport *LGTransport, client LGIClient) {
    buffer := make([]byte, s.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {

            client.Closed()
            transport.Closed()
            transport.Conn.Close()
            s.removeClient(transport.Cid)
            LGError(err)
            break
        }

        LGTrace("read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        LGTrace("transport.Buff", transport.Stream.Bytes())
        n, dps := transport.Fetch()
        LGTrace("fetch message number", n)
        if n > 0 {
            client.ProcessDPs(dps)
        }
    }
    LGTrace("TransportReader stopped for ", transport.Cid)
}

func (s *LGServer) transportSender(transport *LGTransport, client LGIClient) {
    for {
        select {
        case dp := <-transport.outgoing:
            LGTrace("transportSender outgoing:",dp.Type, len(dp.Data))
            //buf := s.Datagram.Pack(dp)
            //transport.Conn.Write(buf)

            transport.PackWrite(dp)
        case <-transport.Quit:
            LGDebug("Transport ", transport.Cid, " quitting")

            transport.Closed()
            transport.Conn.Close()
            s.removeClient(transport.Cid)
            break
        }
    }
}

func (s *LGServer) broadcastHandler(broadcastChan <-chan *LGDataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        LGTrace("broadcastHandler: chan Waiting for input")
        dp := <-broadcastChan
        //buf := s.Datagram.pack(dp)

        fromCid := dp.FromCid
        for _, c := range s.Clients.All() {
            //if fromCid == Cid {
            //    continue
            //}
            if c.GetType() != LGCLIENT_TYPE_GATE {
                dp.FromCid = 0
                dp.Type = LGDATAPACKET_TYPE_GENERAL
            } else {
                dp.FromCid = fromCid
                dp.Type = LGDATAPACKET_TYPE_BROADCAST
            }
            c.GetTransport().outgoing <- dp
        }
        LGTrace("broadcastHandler: Handle end!")
    }
}

//send broadcast message data for other object
func (s *LGServer) SendBroadcast(transport *LGTransport, dp *LGDataPacket) {
    s.broadcastChan <- dp
}

