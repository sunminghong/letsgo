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
    "bytes"
    "container/list"
    "fmt"
    "net"
)


type Server struct {
    boardcast_chan_num int
    read_buffer_size   int

    newclient newClientFunc
    datagram  IDatagram

    //define client dict/map set
    ClientMap map[int]Client

    ClientNum int

    boardcastChan chan dataPacket
}

func newServer(newclient newClientFunc, datagram IDatagram, config map[string]interface{}) {
    s := &Server{}

    if newclient == nil {
        s.newclient = newClient
    } else {
        s.newclient = newclient
    }

    if s.datagram == nil {
        s.datagram = &Datagram{}
    } else {
        s.datagram = datagram
    }

    boardcast_chan_num = 10
    read_buffer_size = 1024
}

func (s *Server) start(host string, port int) {
    Log("Hello Server!")

    addr = host + ":" + string(port)
    s.ClientMap = make(map[int]Client)
    s.newclient = newclient

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.BoardcastChan = make(chan dataPacket, s.boardcast_chan_num)
    go s.boardcastHandler(s.BoardcastChan)

    netListen, error := net.Listen("tcp", addr)
    if error != nil {
        Log(error)
    } else {
        //defer函数退出时执行
        defer netListen.Close()
        for {
            Log("Waiting for clients")
            connection, error := netListen.Accept()
            if error != nil {
                Log("Client error: ", error)
            } else {
                newclientdi := s.allocClientid()
                go ClientHandler(newclientid, connection)
            }
        }
    }
}
func (s *Server) allocClientid() {
    s.ClientNum += 1
    return s.ClientNum
}

//该函数主要是接受新的连接和注册用户在client list
func (s *Server) ClientHandler(newclientid int, connection *net.Conn) {
    Log("one new player connectting ! ")

    newClient := s.newclient(newclientid, connection, s)

    //创建go的线程 使用Goroutine
    go s.clientSender(newClient)
    go s.clientReader(newClient)

    s.ClientMap[newclientid] = newClient

}

func (s *Server) clientReader(client *Client) {
    buffer := make([]byte, s.read_buffer_size)
    for {

        bytesRead, err := client.Conn.Read(buffer)

        if err != nil {
            client.Close()
            s.RemoveClient(client.CID)
            Log(err)
            break
        }
        ByteApend(client.Buff, buffer)

        s.datagram.Fetch(client)
    }
    Log("ClientReader stopped for ", c.Clientid)
}

func (s *Server) clientSender(client *Client) {
    for {
        select {
        case dp := <-client.Outgoing:
            Log(dp.Type, dp.Data)
            buf := s.datagram.pack(dp)
            client.Conn.Write(buf)
        case <-c.Quit:
            Log("Client ", client.Name, " quitting")
            client.Conn.Close()
            break
        }
    }
}

func (s *Server) boardcastHandler(boardcastChan <-chan dataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        Log("boardcastHandler: chan Waiting for input")
        dp := <-boardcastChan
        //buf := s.datagram.pack(dp)

        sendcid = int(dp.Other)
        for cid, c = range s.ClientMap {
            if sendcid == cid {
                continue
            }
            c.Outgoing <- dp
        }
        Log("boardcastHandler: Handle end!")
    }
}

func (s *Server) SendBoardcast(client Client, data []byte) {
    dp = &DataPacket{Type: DATAPACKET_TYPE_BOARDCAST, Data: data, Other: client.CID}
    s.boardcastChan <- dp
}

func (s *Server) Write(client Client, dataType int, data []byte) {
    dp = &DataPacket{Type: dataType, Data: data}
    c.Outgoing <- dp

}

func (s *Server) removeClient(cid int) {

    a,ok := s.ClientMap[c.CID]
    if ok {
        delete(s.ClientMap,c.CID)
    }
}
