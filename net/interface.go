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

//datagram and datapacket define
type IDatagram interface {
    //Encrypt([]byte)
    //Decrypt([]byte)

    Fetch(c *Transport) (n int, dps []*DataPacket)
    Pack(dp *DataPacket) []byte
}

//define proto
type NewProtocolFunc func(name string,transport *Transport) IProtocol

type IProtocol interface {
    ProcessDPs(dps []*DataPacket)
    Close()
    Closed()
    GetTransport() *Transport

    /*
    SetStatus(status int)
    // return this protocol status ,=0 connected =1 disconnect =2 pause
    GetStatus() int
    */
}

type IServer interface {
    SendDP(t *Transport,dataType byte, data []byte)

    SendBoardcast(t *Transport,data []byte)
}


type newTransportFunc func(newcid int, conn net.Conn, server IServer) *Transport

