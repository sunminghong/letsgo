/*=============================================================================
#     FileName: proto.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-14 08:04:40
#      History:
=============================================================================*/
package net

import (
    "net"
    "strconv"
    "time"
)

type Client struct {
    read_buffer_size int

    newproto NewProtocolFunc
    datagram IDatagram
    /*runloop IRunLoop*/

    host string
    port int

    Protos       *ProtocolMap
    TransportNum int
    LastCid int

    localhost string
    localport int
   
    Quit    chan bool

    connaddr chan string
}

func NewClient(newproto NewProtocolFunc, datagram IDatagram /*,runloop IRunLoop*/) *Client {
    c := &Client{Protos: NewProtocolMap()}
    c.newproto = newproto
    c.datagram = datagram
    /*
       if runloop != nil {
           c.runloop = runloop
       } else {
           c.runloop = NewRunLoop()
       }
    */

    c.running = make(map[int]bool)
    c.Quit = make(chan bool)
    c.read_buffer_size = 1024
    return c
}

func (c *Client) Init() {

    <-c.Quit
    return

    for {
        addr := <-c.connaddr
        if addr != "exit" {
            c.Close(0)
            time.Sleep(3)
            return
        }
        c.Start(addr, 0)

        /*
           c.running[0] = true
           // wait for quiting (/quit). run until running ic true
           for c.running[0] {
                   time.Sleep(1 * 1e9)
           }
        */
    }
}

func (c *Client) Start(name string,host string, port int) {
    go func() {
        Log("Hello Client!")

        var addr string
        if port == 0 {
            addr = host
        } else {
            addr = host + ":" + strconv.Itoa(port)
        }

        connection, err := net.Dial("tcp", addr)

        // test(err, "dialing")

        mesg := "dialing"
        if err != nil {
            Log("CLIENT: ERROR: ", mesg)
            return
        } else {
            Log("Ok: ", mesg)
        }
        defer connection.Close()
        Log("main(): connected ")

        newcid := c.allocTransportid()

        transport := NewTransport(newcid, connection, c)
        proto := c.newproto(name,transport)
        c.Protos.Add(newcid, proto)
        c.running[newcid] = true

        //创建go的线程 使用Goroutine
        go c.transportSender(transport)
        go c.transportReader(transport, proto)


        time.Sleep(3)

        <-transport.Quit
        /*
           // wait for quiting (/quit). run until running ic true
           for c.running[newcid] {
                   time.Sleep(1 * 1e9)
           }*/
    }()
    <-c.Quit
}

func (c *Client) Close(cid int) {
    if cid == 0 {
        for cid_, _ := range c.running {
            //c.running[cid] = false
            c.Protos.Get(cid_).GetTransport().Quit <- true
        }
        return
    }

    //c.running[cid] = false
    c.Protos.Get(cid).GetTransport().Quit <- true
}

func (c *Client) removeClient(cid int) {
    c.Protos.Remove(cid)
}

func (c *Client) allocTransportid() int {
    c.TransportNum += 1
    c.LastCid = c.TransportNum
    return c.TransportNum
}

func (c *Client) transportReader(transport *Transport, proto IProtocol) {
    buffer := make([]byte, c.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {
            proto.Closed()
            transport.Closed()
            c.removeClient(transport.Cid)
            Log(err)
            break
        }

        Log("read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        Log("transport.Buff", transport.Stream.Bytes())
        n, dps := c.datagram.Fetch(transport)
        Log("fetch message number", n)
        if n > 0 {
            proto.ProcessDPs(dps)
        }
    }
    Log("TransportReader stopped for ", transport.Cid)
}

func (c *Client) transportSender(transport *Transport) {
    for {
        select {
        case dp := <-transport.Outgoing:
            Log(dp.Type, dp.Data)
            buf := c.datagram.Pack(dp)
            transport.Conn.Write(buf)
        case <-transport.Quit:
            Log("Transport ", transport.Cid, " quitting")
            transport.Conn.Close()
            break
        }
    }
}

func (c *Client) boardcastHandler(boardcastChan <-chan *DataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        Log("boardcastHandler: chan Waiting for input")
        dp := <-boardcastChan
        //buf := c.datagram.pack(dp)

        sendCid, ok := dp.Other.(int)
        if !ok {
            sendCid = 0
        }

        for Cid, c := range c.Protos.All() {
            if sendCid == Cid {
                continue
            }
            c.GetTransport().Outgoing <- dp
        }
        Log("boardcastHandler: Handle end!")
    }
}

//send boardcast message data for other object
func (c *Client) SendBoardcast(transport *Transport, data []byte) {
    //dp := &DataPacket{Type: DATAPACKET_TYPE_BOARDCAST, Data: data, Other: transport.Cid}
    //c.boardcastChan <- dp
}

//send message data for other object
func (c *Client) SendDP(transport *Transport, dataType byte, data []byte) {
    dp := &DataPacket{Type: dataType, Data: data}
    transport.Outgoing <- dp

}

/*
type RunLoop struct {
    //running map[int]bool
    running bool
}

func (r *RunLoop) Loop() {
    r.running = true
    // wait for quiting (/quit). run until running ic true
    for r.running {
            time.Sleep(1 * 1e9)
    }
}

func (r *RunLoop) Stop() {
    r.running = false
}

func (r *RunLoop) Status() bool{
    return r.running
}

func NewRunLoop() *RunLoop {
    return &RunLoop{true}
}

type IRunLoop interface{
    Loop()
    Stop()
    Status()
}
*/
