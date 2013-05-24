/*=============================================================================
#     FileName: interface.go
#         Desc: server base
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-14 09:25:09
#      History:
=============================================================================*/
package net

import (
    "net"
)

//define a struct or class of rec transport connection
type DataPacket struct {
    Type  byte
    Code  uint16
    Data  []byte
    Other interface{}
}

//datagram and datapacket define
type IDatagram interface {
    //Encrypt([]byte)
    //Decrypt([]byte)

    GetEndian() int
    Fetch(c *Transport) (n int, dps []*DataPacket)
    Pack(dp *DataPacket) []byte
}

//define client
type NewClientFunc func(name string, transport *Transport) IClient

type IClient interface {
    GetName() string
    ProcessDPs(dps []*DataPacket)
    Close()
    Closed()
    GetTransport() *Transport
    SendMessage(msg IMessageWriter)
    SendBoardcast(msg IMessageWriter)

    /*
       SetStatus(status int)
       // return this client status ,=0 connected =1 disconnect =2 pause
       GetStatus() int
    */
}

type IServer interface {
    SetMaxConnections(max int)

    SendDP(t *Transport, dataType byte, data []byte)

    SendBoardcast(t *Transport, data []byte)
}

type newTransportFunc func(
    newcid int, conn net.Conn, server IServer) *Transport

type IRouter interface {
    Init()
    //Add(client IClient,protocols []int)
    Add(cid int, protocols string)
    Handler(dp DataPacket) (cid int, ok bool)
    ParseProtos(messageCode int) int
}

type IMessageWriter interface {
    SetCode(code int, ver byte)

    preWrite(wind int)
    writeMeta(datatype int)
    WriteUint16(x int, wind int)
    WriteUint32(x int, wind int)

    WriteUint(x int, wind int)
    WriteInt(x int, wind int)

    WriteString(x string, wind int)
    //WriteList(list *MessageListWriter, wind int)

    //对数据进行封包
    ToBytes() []byte
}

/////////////////////////////////////////////////////////////////////////////////

type IMessageReader interface {
    ReadUint() int
    ReadInt() int
    ReadUint32() int
    ReadUint16() int
    ReadString() string
    //ReadList() *MessageListReader
}

