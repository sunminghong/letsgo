/*=============================================================================
#     FileName: init.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-12-24 18:31:40
#      History:
=============================================================================*/

package protos

import (
    "os"
    . "github.com/sunminghong/letsgo/net"
    "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

var Endian int = helper.LGLittleEndian

var Handlers map[int]LGProcessHandleFunc= make(map[int]LGProcessHandleFunc)


func processHandl(code int,msg LGIMessageReader,c LGIConnection) {
    h, ok := Handlers[code]
    if ok {
        h(msg,c,0)
    }
}

// LGIConnection  
type Connection struct {
    *LGBaseConnection
    Username *string
}

func NewConnection (name string,transport *LGTransport) LGIConnection {
    username := "someone"
    //c := &Connection{}
    //c.Username = &username
    //c.Name = name
    //c.Transport = transport

    c := &Connection{
        &LGBaseConnection{Transport:transport,Name:name},
        &username,
    }
    return c
}

//对数据进行拆包
func (c *Connection) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        msg := LGNewMessageReader(dp.Data,Endian)
        LGTrace("msg.code:",msg.Code,msg.Ver)

        //todo: route don't execute
        processHandl(msg.Code,msg,c)
    }
}

func (c *Connection) Closed() {
}

/*
func (c *Connection) SendMessage(msg LGIMessageWriter) {
    c.Transport.SendDP(0,msg.ToBytes())
}

func (c *Connection) SendBroadcast(msg LGIMessageWriter) {
    c.Transport.SendBroadcast(msg.ToBytes())
}
*/

func NewMessageWriter(c LGIConnection) *LGMessageWriter {
    return LGNewMessageWriter(c.GetTransport().Stream.Endian)
}

func init() {
    Handlers[201] = Process201
    Handlers[2001] = Process2001
    Handlers[2011] = Process2011
    Handlers[2101] = Process2101

    Handlers[2020] = Process2020
    Handlers[2021] = Process2021
}

func logfightdata(data string) {
    LGTrace(data)

    readerFile := "/Users/Team1201/works/tmp/fightdata.txt"
    fout, err := os.OpenFile(readerFile, os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
    if err != nil {
        LGTrace("Error:", err)
        return
    }

    defer fout.Close()

    fout.WriteString(data)
}
