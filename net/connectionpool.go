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
    . "github.com/sunminghong/letsgo/log"
)

type LGConnectionPool struct {
    read_buffer_size int

    newclient LGNewConnectionFunc
    datagram LGIDatagram
    /*runloop IRunLoop*/

    host string
    port int

    Connections       *LGConnectionMap
    TransportNum int

    localhost string
    localport int

    Quit    chan bool
    broadcastChan    chan *LGDataPacket

    connaddr chan string
}

func LGNewConnectionPool(newclient LGNewConnectionFunc, datagram LGIDatagram ) *LGConnectionPool {
    cp := &LGConnectionPool{Connections: LGNewConnectionMap()}
    cp.newclient = newclient
    cp.datagram = datagram

    cp.Quit = make(chan bool)
    cp.read_buffer_size = 1024

    //创建一个管道 chan map 需要make creates slices, maps, and channels only
    cp.broadcastChan = make(chan *LGDataPacket,1)
    go cp.broadcastHandler(cp.broadcastChan)

    return cp
}

func (cp *LGConnectionPool) Start(name string,addr string,datagram LGIDatagram) {
    //go func() {
        ////Log("Hello Connection!")

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

        newcid := cp.allocTransportid()

        if datagram == nil {
            datagram = cp.datagram
        }
        transport := LGNewTransport(newcid, connection, cp,datagram)
        client := cp.newclient(name,transport)
        cp.Connections.Add(newcid,name, client)

        //创建go的线程 使用Goroutine
        go cp.transportSender(transport)
        go cp.transportReader(transport, client)


        time.Sleep(2)

        <-transport.Quit
    //}()
    <-cp.Quit
}

func (cp *LGConnectionPool) SetMaxConnections(max int) {

}

func (cp *LGConnectionPool) Close(cid int) {
    if cid == 0 {
        for _, client := range cp.Connections.All(){
            //c.running[cid] = false
            client.GetTransport().Quit <- true
        }
        return
    }

    //c.running[cid] = false
    cp.Connections.Get(cid).GetTransport().Quit <- true
}

func (cp *LGConnectionPool) removeConnection(cid int) {
    cp.Connections.Remove(cid)
}

func (cp *LGConnectionPool) allocTransportid() int {
    cp.TransportNum += 1
    return cp.TransportNum
}

func (cp *LGConnectionPool) transportReader(transport *LGTransport, client LGIConnection) {
    buffer := make([]byte, cp.read_buffer_size)
    for {

        bytesRead, err := transport.Conn.Read(buffer)

        if err != nil {
            client.Closed()
            transport.Closed()
            cp.removeConnection(transport.Cid)
            //Log(err)
            break
        }

        LGTrace("pool transportReader read to buff:", bytesRead)
        transport.BuffAppend(buffer[0:bytesRead])

        LGTrace("tpool transportReader Buff:%d,% X", len(transport.Stream.Bytes()),transport.Stream.Bytes())
        //, transport.Stream.Bytes())
        //LGTrace(transport.Stream.DebugOut())
        n, dps := transport.Fetch()
        LGTrace("fetch message number", n)
        if n > 0 {
            client.ProcessDPs(dps)
        }
    }
    //Log("TransportReader stopped for ", transport.Cid)
}

func (cp *LGConnectionPool) transportSender(transport *LGTransport) {
    for {
        select {
        case dp := <-transport.Outgoing:
            LGTrace("clientpool transportSender:dp.type=%v,dp.data=% X",dp.Type, dp.Data)
            //buf := cp.datagram.Pack(dp)
            //transport.Conn.Write(buf)

            transport.PackWrite(dp)
        case <-transport.Quit:
            //Log("Transport ", transport.Cid, " quitting")
            transport.Conn.Close()

            //client.Closed()
            transport.Closed()
            cp.removeConnection(transport.Cid)
            break
        }
    }
}

func (cp *LGConnectionPool) broadcastHandler(broadcastChan chan *LGDataPacket) {
    for {
        //在go里面没有while do ，for可以无限循环
        //Log("broadcastHandler: chan Waiting for input")
        dp := <-broadcastChan
        //buf := c.datagram.pack(dp)

        sendCid := dp.FromCid
        for Cid, c := range cp.Connections.All() {
            if sendCid == Cid {
                continue
            }
            c.GetTransport().Outgoing <- dp
        }
        //Log("broadcastHandler: Handle end!")
    }
}

//send broadcast message data for other object
func (cp *LGConnectionPool) SendBroadcast(dp *LGDataPacket) {
    cp.broadcastChan <- dp
}

