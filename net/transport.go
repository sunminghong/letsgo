package net

import (
    "net"
    . "github.com/sunminghong/letsgo/helper"
)

type LGTransport struct {
    //transport connection 唯一 id
    Cid  int

    //需要输出的数据包的channel
    outgoing chan *LGDataPacket

    //需要输出的数据流 的channel
    outgoingBytes chan []byte

    Quit chan bool

    Stream *LGRWStream
    DataType byte
    DPSize  int

    datagram LGIDatagram
    Server LGIServer
    Conn   net.Conn
}


//define method what Close transport's connection for struct Transport
func (c *LGTransport) Close() {
    c.Quit <- true
    c.Conn.Close()
}

//define method what Close transport's connection for struct Transport
func (c *LGTransport) Closed() {
    //
}

func (c *LGTransport) Equal(other *LGTransport) bool {
    if c.Cid == other.Cid {
        return true
    }
    return false
}

func (c *LGTransport) InitBuff() {
    c.Stream.Reset()
}

// Write appends the contents of p to the []byte.  The return
// value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with
// ErrTooLarge.
func (c *LGTransport) BuffAppend(p []byte) (n int) {
    return c.Stream.Write(p)
}

func (c *LGTransport) Fetch() (n int, dps []*LGDataPacket) {
    return c.datagram.Fetch(c)
}

func (c *LGTransport) PackWrite(dp *LGDataPacket) {
    c.datagram.PackWrite(c.Conn.Write,dp)
}

func (c *LGTransport) SendDP(dp *LGDataPacket) {
    c.outgoing <- dp
}

func (c *LGTransport) SendBroadcast(dp *LGDataPacket) {
    c.Server.SendBroadcast(c, dp)
}

// new Transport object
func LGNewTransport(newcid int, conn net.Conn, server LGIServer,datagram LGIDatagram) *LGTransport {
    c := &LGTransport{
        Cid:      newcid,
        Conn:     conn,
        datagram: datagram,
        Server:   server,
        outgoing: make(chan *LGDataPacket, 1),
        outgoingBytes: make(chan []byte),
        Quit:     make(chan bool),
        Stream:   LGNewRWStream(1024,datagram.GetEndian()),
    }

    c.InitBuff()

    return c
}
