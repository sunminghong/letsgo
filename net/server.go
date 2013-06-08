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
    . "github.com/sunminghong/letsgo/log"
)

type LGServer struct {
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
}

func LGNewServer(makeclient LGNewClientFunc, datagram LGIDatagram) *LGServer {
    s := &LGServer{Clients: LGNewClientMap()}

    s.makeclient = makeclient

    s.Datagram = datagram

    s.broadcast_chan_num = 10
    s.read_buffer_size = 1024

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
                newcid := s.allocTransportid()
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
}

func (s *LGServer) allocTransportid() int {
    if (s.Clients.Len() >= s.maxConnections) {
        return 0
    }
    s.TransportNum += 1
    return s.TransportNum
}

//该函数主要是接受新的连接和注册用户在transport list
func (s *LGServer) transportHandler(newcid int, connection net.Conn) {
    transport := LGNewTransport(newcid, connection, s,s.Datagram)
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
        n, dps := s.Datagram.Fetch(transport)
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

            s.Datagram.PackWrite(transport.Conn.Write,dp)
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

