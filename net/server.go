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
)

type Server struct {
    boardcast_chan_num int
    read_buffer_size   int

    newclient newClientFunc
    datagram  IDatagram

    host    string
    port    int

    //define client dict/map set
    ClientMap map[int]*Client

    ClientNum int

    boardcastChan chan *DataPacket
}

func NewServer(newclient newClientFunc, datagram IDatagram, config map[string]interface{}) *Server {
    s := &Server{}

    s.newclient = newclient

    s.datagram = datagram

    s.boardcast_chan_num = 10
    s.read_buffer_size = 1024

    return s
}

func (s *Server) Start(host string, port int) {
    Log("Hello Server!")

    addr := host + ":" + strconv.Itoa(port)
    s.ClientMap = make(map[int]*Client)

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    s.boardcastChan = make(chan *DataPacket, s.boardcast_chan_num)
    go s.boardcastHandler(s.boardcastChan)

    Log("listen with :",addr)
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
                newclientid := s.allocClientid()
                go s.clientHandler(newclientid, connection)
            }
        }
    }
}

func (s *Server) removeClient(Cid int) {

    _, ok := s.ClientMap[Cid]
    if ok {
        delete(s.ClientMap, Cid)
    }
}
func (s *Server) allocClientid() int {
    s.ClientNum += 1
    return s.ClientNum
}

//该函数主要是接受新的连接和注册用户在client list
func (s *Server) clientHandler(newclientid int, connection net.Conn) {
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
            s.removeClient(client.Cid)
            Log(err)
            break
        }
        Log("read to buff:",bytesRead,buffer)
        client.BuffAppend(buffer[0:bytesRead])

        Log("client.Buff",client.Buff)
        s.datagram.Fetch(client)
    }
    Log("ClientReader stopped for ", client.Cid)
}

func (s *Server) clientSender(client *Client) {
    for {
        select {
        case dp := <-client.Outgoing:
            Log(dp.Type, dp.Data)
            buf := s.datagram.Pack(dp)
            client.Conn.Write(buf)
        case <-client.Quit:
            Log("Client ", client.Cid, " quitting")
            client.Conn.Close()
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

        sendCid,ok := dp.Other.(int)
        if !ok {
            sendCid = 0
        }

        for Cid, c := range s.ClientMap {
            if sendCid == Cid {
                continue
            }
            c.Outgoing <- dp
        }
        Log("boardcastHandler: Handle end!")
    }
}

//send boardcast message data for other object
func (s *Server) SendBoardcast(client *Client, data []byte) {
    dp := &DataPacket{Type: DATAPACKET_TYPE_BOARDCAST, Data: data, Other: client.Cid}
    s.boardcastChan <- dp
}

//send message data for other object
func (s *Server) SendMsg(client *Client, dataType int, data []byte) {
    dp := &DataPacket{Type: dataType, Data: data}
    client.Outgoing <- dp

}

