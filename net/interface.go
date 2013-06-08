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

const (
    DATAPACKET_TYPE_GENERAL = 0
    DATAPACKET_TYPE_DELAY = 1
    DATAPACKET_TYPE_BOARDCAST = 3
    DATAPACKET_TYPE_GATECONNECT = 3
)

//define a struct or class of rec transport connection
type DataPacket struct {
    Type  byte
    Data  []byte

    FromCid int
}

type WriteFunc func (data []byte) (int,error)

//datagram and datapacket define
type IDatagram interface {
    //Encrypt([]byte)
    //Decrypt([]byte)

    Clone(endian int) IDatagram
    GetEndian() int
    SetEndian(endian int)
    Fetch(c *Transport) (n int, dps []*DataPacket)
    //Pack(dp *DataPacket) []byte
    PackWrite(write WriteFunc,dp *DataPacket) []byte
}


//define client
type ProcessHandleFunc func(
    code int,msg *MessageReader,c IClient,fromCid int)

type NewClientFunc func(name string, transport *Transport) IClient

const (
    CLIENT_TYPE_GENERAL = 0
    CLIENT_TYPE_GATE = 1
)
type IClient interface {

    GetType() int
    SetType(t int)
    GetName() string
    ProcessDPs(dps []*DataPacket)
    Close()
    Closed()
    GetTransport() *Transport
    SendMessage(fromcid int,msg IMessageWriter)
    SendBoardcast(fromcid int,msg IMessageWriter)

    /*
       SetStatus(status int)
       // return this client status ,=0 connected =1 disconnect =2 pause
       GetStatus() int
    */
}

type IServer interface {
    SetMaxConnections(max int)

    //SendDP(t *Transport, dp *DataPacket)

    SendBoardcast(t *Transport, dp *DataPacket)
}

type newTransportFunc func(
    newcid int, conn net.Conn, server IServer) *Transport

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

///////////////////////////////////////////////////////////////////////////////

type IMessageReader interface {
    ReadUint() int
    ReadInt() int
    ReadUint32() int
    ReadUint16() int
    ReadString() string
    //ReadList() *MessageListReader
}

type IIDAssign interface {
    GetFree() int
    Free(id int)
}
