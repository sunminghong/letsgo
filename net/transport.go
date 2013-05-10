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

    Stream *RWStream
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
    c.Stream.Reset()
}

// Write appends the contents of p to the []byte.  The return
// value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with
// ErrTooLarge.
func (c *Transport) BuffAppend(p []byte) (n int) {
    return c.Stream.Write(p)
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
        Stream:       NewRWStream(make([]byte,1024),true)
    }

    c.InitBuff()

    return c
}
