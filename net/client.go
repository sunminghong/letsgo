package net

import (
    "net"
)

type Client struct {
    //client connection 唯一 id
    Cid  int

    //需要输出的数据(protocolcode+body) 的channel
    Outgoing chan *DataPacket

    Quit chan bool

    Buff     []byte
    DataType int
    MsgSize  int

    Server *Server
    Conn   net.Conn
}

type newClientFunc func(newclientid int, conn net.Conn, server *Server) *Client

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
    if c.Cid == other.Cid {
        return true
    }
    return false
}

func (c *Client) InitBuff() {
    //c.Buff = make([]byte,0)
    c.Buff = []byte{}
}


//add buff cap
func (c *Client)buffGrow(addlen int) int{
    m := len(c.Buff)
    if m + addlen > cap(c.Buff) {
        var b_ []byte
        // not enough space anywhere
        b_ = make([]byte,m+addlen)
        copy(b_, c.Buff)
        Log("b_",b_)
        c.Buff = b_
    }
    return m
}

// Write appends the contents of p to the []byte.  The return
// value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with
// ErrTooLarge.
func (c *Client) BuffAppend(p []byte) (n int) {
    Log("len(buff)=",len(c.Buff),"len(p)=",len(p))
    m:= c.buffGrow(len(p))
    Log("buff",c.Buff)
    a := copy((c.Buff)[m:], p)

    Log("buff2",a,c.Buff)
    return a
}


func (c *Client) SendMsg(dataType int, data []byte) {
    c.Server.SendMsg(c, dataType, data)
}

func (c *Client) SendBoardcast(data []byte) {
    c.Server.SendBoardcast(c, data)
}

////////////////////////////////////////////////////////////////
//process fetch message
// default to echo return,impent
func (c *Client) ProcessMsg(msg *DataPacket) {
    //echo
    c.Outgoing <- msg
}



// new Client object
func NewClient(newclientid int, conn net.Conn, server *Server) *Client {
    c := &Client{
        Cid:      newclientid,
        Conn:     conn,
        Server:   server,
        Outgoing: make(chan *DataPacket, 10),
        Quit:     make(chan bool),
    }

    c.InitBuff()

    return c
}
