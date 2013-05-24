/*=============================================================================
#     FileName: clientpool.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-22 14:19:12
#      History:
=============================================================================*/
package net

import (
    "net"
    //"strconv"
    "time"
)

type ClientPool struct {
    read_buffer_size int

    newclient NewClientFunc
    datagram IDatagram
    /*runloop IRunLoop*/

    host string
    port int

    Clients       *ClientMap
    TransportNum int

    localhost string
    localport int

    Quit    chan bool

    connaddr chan string
}

func NewClientPool(newclient NewClientFunc, datagram IDatagram /*,runloop IRunLoop*/) *ClientPool {
    c := &ClientPool{Clients: NewClientMap()}
    c.newclient = newclient
    c.datagram = datagram
    /*
       if runloop != nil {
           c.runloop = runloop
       } else {
           c.runloop = NewRunLoop()
       }
    */

    c.Quit = make(chan bool)
    c.read_buffer_size = 1024
    return c
}

func (c *ClientPool) Start(name string,addr string) {
    //go func() {
        ////Log("Hello Client!")

        //addr = host + ":" + strconv.Itoa(port)

        connection, err := net.Dial("tcp", addr)

        //mesg := "dialing"
        if err != nil {
            //Log("CLIENT: ERROR: ", mesg)
            return
        } else {
            //Log("Ok: ", mesg)
        }
        defer connection.Close()
        //Log("main(): connected ")

        newcid := c.allocTransportid()

        transport := NewTransport(newcid, connection, c)
        client := c.newclient(name,transport)
        c.Clients.Add(newcid,name, client)

        //创建go的线程 使用Goroutine
        go c.transportSender(transport)
        go c.transportReader(transport, client)


        time.Sleep(2)

        <-transport.Quit
    //}()
    <-c.Quit
}

func (c *ClientPool) Close(cid int) {
    if cid == 0 {
        for _, client := range c.Clients.All(){
            //c.running[cid] = false
            client.GetTransport().Quit <- true
        }
        return
    }

    //c.running[cid] = false
    c.Clients.Get(cid).GetTransport().Quit <- true
}

func (c *ClientPool) removeClient(cid int) {
    c.Clients.Remove(cid)
}

func (c *ClientPool) allocTransportid() int {
    c.TransportNum += 1
    return c.TransportNum
}

func (c *ClientPool) transportReader(transport *Transport, client IClient) {
    buffer := make([]byte, c.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {
            client.Closed()
            transport.Closed()
            c.removeClient(transport.Cid)
            //Log(err)
            break
        }

        //Log("read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        //Log("transport.Buff", transport.Stream.Bytes())
        n, dps := c.datagram.Fetch(transport)
        //Log("fetch message number", n)
        if n > 0 {
            client.ProcessDPs(dps)
        }
    }
    //Log("TransportReader stopped for ", transport.Cid)
}

func (c *ClientPool) transportSender(transport *Transport) {
    for {
        select {
        case dp := <-transport.Outgoing:
            Trace("clientpool transportSender:",dp.Type, dp.Data)
            buf := c.datagram.Pack(dp)
            transport.Conn.Write(buf)
        case <-transport.Quit:
            //Log("Transport ", transport.Cid, " quitting")
            transport.Conn.Close()

            //client.Closed()
            transport.Closed()
            c.removeClient(transport.Cid)
            break
        }
    }
}

func (c *ClientPool) boardcastHandler(boardcastChan <-chan *DataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        //Log("boardcastHandler: chan Waiting for input")
        dp := <-boardcastChan
        //buf := c.datagram.pack(dp)

        sendCid, ok := dp.Other.(int)
        if !ok {
            sendCid = 0
        }

        for Cid, c := range c.Clients.All() {
            if sendCid == Cid {
                continue
            }
            c.GetTransport().Outgoing <- dp
        }
        //Log("boardcastHandler: Handle end!")
    }
}

//send boardcast message data for other object
func (c *ClientPool) SendBoardcast(transport *Transport, data []byte) {
    //dp := &DataPacket{Type: DATAPACKET_TYPE_BOARDCAST, Data: data, Other: transport.Cid}
    //c.boardcastChan <- dp
}

//send message data for other object
func (c *ClientPool) SendDP(transport *Transport, dataType byte, data []byte) {
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
