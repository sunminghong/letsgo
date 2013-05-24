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
    DataType byte
    DPSize  int

    Server IServer
    Conn   net.Conn
}


//define method what Close transport's connection for struct Transport
func (c *Transport) Close() {
    c.Quit <- true
    c.Conn.Close()
}

//define method what Close transport's connection for struct Transport
func (c *Transport) Closed() {
    //
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

func (c *Transport) SendDP(dataType byte, data []byte) {
    if data == nil {
        return 
    }
    c.Server.SendDP(c, dataType, data)
}

func (c *Transport) SendBoardcast(data []byte) {
    if data == nil {
        return 
    }
    c.Server.SendBoardcast(c, data)
}

// new Transport object
func NewTransport(newcid int, conn net.Conn, server IServer) *Transport {
    c := &Transport{
        Cid:      newcid,
        Conn:     conn,
        Server:   server,
        Outgoing: make(chan *DataPacket, 10),
        Quit:     make(chan bool),
        Stream:   NewRWStream(1024,BigEndian),
    }

    c.InitBuff()

    return c
}
