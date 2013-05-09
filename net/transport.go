package net

import (
    "net"
)

type Transport struct {
    //transport connection 唯一 id
    Cid  int

    //需要输出的数据(protocolcode+body) 的channel
    Outgoing chan *DataPacket

    Quit chan bool

    Buff     []byte
    DataType int
    DPSize  int

    Server *Server
    Conn   net.Conn
}

type newTransportFunc func(newcid int, conn net.Conn, server *Server) *Transport

//define method what Close transport's connection for struct Transport
func (c *Transport) Close() {
    c.Quit <- true
    c.Conn.Close()
}

func (c *Transport) Equal(other *Transport) bool {
    if c.Cid == other.Cid {
        return true
    }
    return false
}

func (c *Transport) InitBuff() {
    //c.Buff = make([]byte,0)
    c.Buff = []byte{}
}


//add buff cap
func (c *Transport)buffGrow(addlen int) int{
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
func (c *Transport) BuffAppend(p []byte) (n int) {
    Log("len(buff)=",len(c.Buff),"len(p)=",len(p))
    m:= c.buffGrow(len(p))
    Log("buff",c.Buff)
    a := copy((c.Buff)[m:], p)

    Log("buff2",a,c.Buff)
    return a
}

func (c *Transport) SendDP(dataType int, data []byte) {
    c.Server.SendDP(c, dataType, data)
}

func (c *Transport) SendBoardcast(data []byte) {
    c.Server.SendBoardcast(c, data)
}

// new Transport object
func NewTransport(newcid int, conn net.Conn, server *Server) *Transport {
    c := &Transport{
        Cid:      newcid,
        Conn:     conn,
        Server:   server,
        Outgoing: make(chan *DataPacket, 10),
        Quit:     make(chan bool),
    }

    c.InitBuff()

    return c
}
