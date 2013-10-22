/*============================================================================= #     FileName: interface.go
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
    LGDATAPACKET_TYPE_GENERAL = 0
    LGDATAPACKET_TYPE_DELAY = 2 + 1
    LGDATAPACKET_TYPE_CLOSE = 4 + 1 //close a player client
    LGDATAPACKET_TYPE_BROADCAST = 6 + 1
    LGDATAPACKET_TYPE_GATE_REGISTER= 8
    LGDATAPACKET_TYPE_GATE_REMOVE= 10 //remove a gate client
    LGDATAPACKET_TYPE_CLOSED = 12 + 1 //a player client closed tell to gridserver
    LGDATAPACKET_TYPE_DELAY_DATAS = 14 + 1
    LGDATAPACKET_TYPE_DELAY_DATAS_COMPRESS = 16 + 1
    LGDATAPACKET_TYPE_DATAS_COMPRESS = 18 //to player client connection
)

//define a struct or class of rec transport connection
type LGDataPacket struct {
    Type  byte
    Data  []byte

    FromCid int
}

type LGWriteFunc func (data []byte) (int,error)

//datagram and datapacket define
type LGIDatagram interface {
    //Encrypt([]byte)
    //Decrypt([]byte)

    Clone(endian int) LGIDatagram
    GetEndian() int
    SetEndian(endian int)
    Fetch(c *LGTransport) (n int, dps []*LGDataPacket)
    Pack(dp *LGDataPacket) []byte
    PackWrite(write LGWriteFunc,dp *LGDataPacket)
}

type LGNewClientFunc func(name string, transport *LGTransport) LGIClient

const (
    LGCLIENT_TYPE_GENERAL = 0
    LGCLIENT_TYPE_GATE = 1
)
type LGIClient interface {

    GetType() int
    SetType(t int)
    GetName() string
    ProcessDPs(dps []*LGDataPacket)
    Close()
    Closed()
    GetTransport() *LGTransport
    SendMessage(fromcid int,msg LGIMessageWriter)
    SendBroadcast(fromcid int,msg LGIMessageWriter)

    /*
       SetStatus(status int)
       // return this client status ,=0 connected =1 disconnect =2 pause
       GetStatus() int
    */
}

type LGIServ interface {
    SetMaxConnections(max int)
    SendBroadcast(dp *LGDataPacket)
}

type LGIServer interface {
    GetName() string
    GetServerid() int

    LGIServ
}

type LGnewTransportFunc func(
    newcid int, conn net.Conn, server LGIServer) *LGTransport

type LGIMessageWriter interface {
    SetCode(code int, ver byte)
    GetCode() int

    preWrite(wind int)
    writeMeta(datatype int)
    WriteUint16(x int, wind int)
    WriteUint32(x int, wind int)

    WriteUint(x int, wind int)
    WriteUints(x ...int)

    WriteInt(x int, wind int)
    WriteInts(x ...int)

    WriteString(x string, wind int)
    //WriteList(list *MessageListWriter, wind int)

	WriteUints(xs ...int)
	WriteInts(xs ...int)

    //对数据进行封包
    ToBytes() []byte
}

///////////////////////////////////////////////////////////////////////////////

type LGIMessageReader interface {
    ReadCode() int
    ReadUint() int
    ReadInt() int
    ReadUint32() int
    ReadUint16() int
    ReadString() string
    //ReadList() *MessageListReader
}

type LGIIDAssign interface {
    GetFree() int
    Free(id int)
}
