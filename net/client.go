package net

import (
    "fmt"
    "net"
    "container/list"
)

//define a struct or class of rec client connection
type Client struct {

    //client connection 唯一 id
    Clientid int

    //接受客户端进来的数据包（已经拆分后的，protocolcode+body）的channel
    Incoming chan []byte

    //需要输出的数据(protocolcode+body) 的channel
    Outgoing chan []byte

    Conn net.Conn
    Quit chan bool

    Buff []byte

    datagramFlag = 
}


//define client dict/map set
var clientMap map[int] *Client


//define method what read client's connection data for struct Client
func (c *Client) Read(buffer []byte) (int,bool) {
    bytesread,err := c.Conn.Read(buffer)

    if err !=nil {
        c.Close()
        Log(err)
        return 0,false
    }

    Log("Read ",bytesRead," bytes")
    return bytesRead, true
}

func (c *Client) Write(buffer []byte) {
    bytesread,err := c.Conn.Write(buffer)

    if err !=nil {
        c.Close()
        Log(err)
        return 0,false
    }
}

//define method what Close client's connection for struct Client
func (c *Client) Close() {
    c.Quit <- true
    c.Conn.Close()
    c.RemoveMe()

}

//define method what Remove self for struct Client
func (c *Client) RemoveMe() {
    global clientMap
    a,ok := clientMap[c.Clientid]
    if ok {
        delete(clientMap,c.Clientid)
    }
}

func (c *Client) Equal(other *Client) bool {
    if c.Clientid == other.Clientid {
        return true
    }
}


