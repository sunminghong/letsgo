/*=============================================================================
#     FileName: init.go
#         Desc: client init
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-27 12:07:21
#      History:
=============================================================================*/
package protos

import (
    . "github.com/sunminghong/letsgo/net"
    "github.com/sunminghong/letsgo/helper"
    "github.com/sunminghong/letsgo/log"
)

var Endian int = helper.LGLittleEndian

var Handlers map[int]LGProcessHandleFunc= make(map[int]LGProcessHandleFunc)


func processHandl(code int,msg LGIMessageReader,c LGIClient) {
    h, ok := Handlers[code]
    if ok {
        h(msg,c,0)
    }
}

// LGIClient  
type Client struct {
    *LGBaseClient
    Username *string
}

func NewClient (name string,transport *LGTransport) LGIClient {
    username := "someone"
    //c := &Client{}
    //c.Username = &username
    //c.Name = name
    //c.Transport = transport

    c := &Client{
        &LGBaseClient{Transport:transport,Name:name},
        &username,
    }
    return c
}

//对数据进行拆包
func (c *Client) ProcessDPs(dps []*LGDataPacket) {
    for _, dp := range dps {
        msg := LGNewMessageReader(dp.Data,Endian)
        log.LGTrace("msg.code:",msg.Code,msg.Ver)

        //todo: route don't execute
        processHandl(msg.Code,msg,c)
    }
}

func (c *Client) Closed() {
}

/*
func (c *Client) SendMessage(msg LGIMessageWriter) {
    c.Transport.SendDP(0,msg.ToBytes())
}

func (c *Client) SendBroadcast(msg LGIMessageWriter) {
    c.Transport.SendBroadcast(msg.ToBytes())
}
*/

func NewMessageWriter(c LGIClient) *LGMessageWriter {
    return LGNewMessageWriter(c.GetTransport().Stream.Endian)
}

func init() {
    Handlers[2011] = Process2011
    Handlers[2001] = Process2001
}

