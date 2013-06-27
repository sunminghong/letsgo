/*=============================================================================
#     FileName: gate.go
#         Desc: game grid server
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:09:28
#      History:
=============================================================================*/
package grids

import (
    //"reflect"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/net/grid"
    p "./protos"
)

/*
*
*
		a := "sssssss"
			args := []reflect.Value{reflect.ValueOf(a)}
			c := reflect.ValueOf(test).Call(args)
			fmt.Println(c)

            CanInterface()

            func (t *commonType) NumMethod() int
返回该类型拥有的方法数量

            func (t *commonType) Implements(u Type) bool
判断类型是否实现u这个接口.注意u必须不能为nil,且是一个接口

		type B struct {
				c string
				b byte
				a int
			}

			func (b B) test() {

			}

			func main() {
				b := B{}
				fmt.Println(reflect.TypeOf(b).Method(0).Name)  //test
			}


*/
type Client struct {
    *LGGridClient
}

func (c *Client) Closed() {
    LGTrace("a grid client closed")

    //msg := "system: " + (*c.Username) + " is leave!"
    //mw := lnet.NewMessageWriter(c.Transport.Stream.Endian)
    //mw.SetCode(2011,0)
    //mw.WriteString(msg,0)

    //c.Transport.SendBroadcast(mw.ToBytes())
}

func NewClient(name string,transport *LGTransport) LGIClient {
    cg := &LGGridClient{
        &LGBaseClient{Transport:transport,Name:name},
        ProccessHandle,
    }

    c := &Client{ cg }

    return c
}


type ProccessFunc func(msg LGIMessageReader,c LGIClient,fromCid int,session *p.Session)

var Handlers map[int]ProccessFunc= make(map[int]ProccessFunc)

func ProccessHandle(msg LGIMessageReader,c LGIClient,fromCid int) {
    LGTrace("message is request")
    code := msg.ReadCode()
    h, ok := Handlers[code]
    if ok {
        uid := p.Uidmap.GetUid(fromCid,c.GetTransport().Cid)
        if uid == 0 {
            h(msg,c,fromCid,nil)
        } else {
            sess:=p.GetSession(uid)
            h(msg,c,fromCid,sess)
        }
    }
}


func init() {
    //uidmap = NewLGUidMap()
    Handlers[1011] = p.Process1011
    //Handlers[2011] = p.Process2011
}

