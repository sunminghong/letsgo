package net

import (
    "fmt"
    "net"
    "container/list"
    "bytes"
)

type Client interface {
    //client connection 唯一 id
    CID    int
    Conn        net.Conn


    //需要输出的数据(protocolcode+body) 的channel
    Outgoing    chan DataPacket

    ProcessMsg(msg *dataPacket)

    InitBuff()

    Close()


    Buff        []byte
    DataType    uint
    MsgSize     uint


    Quit        chan bool

    server      Server
    Conn        net.Conn
}

//define client dict/map set
var clientMap map[int] *Client


func (c *Client) Write(dataType uint, data []byte) {
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
}

func (c *Client) Equal(other *Client) bool {
    if c.CID == other.CID{
        return true
    }
}

//process fetch message
func (c *Client) ProcessMsg(msg *dataPacket) {
    c.Incoming <- msg
}

func(c *Client) InitBuff(){
    c.Buff = [1024]byte{}
}


type newClientFunc func(clientid int, conn net.Conn, server Server) Client

// new Client object
func NewClient(newclientid int,connection *net.Conn,server *Server) Client {
    c := &Client{
        CID: newclientid,
        Conn: connection,
        server: server,
        Incoming: make(chan dataPacket,10),
        Outgoing: make(chan dataPacket,10),
        Quit:make(chan bool)
    }

    c.InitBuff()

    return c
}


