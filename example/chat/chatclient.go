/*=============================================================================
#     FileName: chatclient.go
#         Desc: chat client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-13 17:48:50
#      History:
=============================================================================*/
package main

import (
    "./lib"
    "bufio"
    "fmt"
    lnet "github.com/sunminghong/letsgo/net"
    "os"
    "strings"
    "time"
    "strconv"
)

// LGIClient  
type Client struct {
    Transport *lnet.Transport
    Name      string
    Username    *string
}

func LGMakeClient(name string,transport *lnet.Transport) lnet.LGIClient {
    username := "someone"
    c := &Client{transport, name, &username}

    return c
}

func (c *Client) GetName() string {
    return c.Name
}

func (c *Client) ProcessDPs(dps []*lnet.LGDataPacket) {
    for _, dp := range dps {
        md := string(dp.Data)

        fmt.Println()
        fmt.Println(md)
        fmt.Print("you> ")
    }
}

//对数据进行拆包
func (c *Client) GetTransport() *lnet.Transport {
    return c.Transport
}

func (c *Client) Close() {
    c.Transport.Close()
}

func (c *Client) Closed() {
    msg := "system: " + (*c.Username) + " is leave!"
    c.Transport.SendBroadcast([]byte(msg))
}

func (c *Client) SendMessage(msg lnet.LGIMessageWriter) {
    c.Transport.SendDP(0,msg.ToBytes())
}

func (c *Client) SendBroadcast(msg lnet.LGIMessageWriter) {
    c.Transport.SendBroadcast(msg.ToBytes())
}

// clientsender(): read from stdin and send it via network
func clientsender(cid *int,client *lnet.ClientPool) {
    reader := bufio.NewReader(os.Stdin)
    for {
        if (*cid)==0 {
            fmt.Print("you no connect anyone server,please input conn cmd,\n")
        }
        fmt.Print("you> ")
        input, _ := reader.ReadBytes('\n')
        cmd := string(input[:len(input)-1])
        if cmd[0] == '/' {
            cmds := strings.Split(cmd," ")
            switch cmds[0]{
            case "/conn":
                var name,addr string
                if len(cmds)>2 {
                    name = cmds[1]
                    addr = cmds[2]
                }else {
                    name = "c_" + strconv.Itoa(*cid)
                    addr = cmds[1]
                }

                p := client.Clients.GetByName(name)
                if p != nil {
                    fmt.Println(name," is exists !")
                    continue
                }

                go client.Start(name,addr)


                fmt.Print("please input your name:")
                input, _ := reader.ReadBytes('\n')
                input =input[0:len(input)-1]

                for true {
                    b := client.Clients.GetByName(name)
                    if b!=nil{
                        change(cid,client,name)
                        break
                    }
                    time.Sleep(2*1e3)
                }
                client.Clients.Get(*cid).GetTransport().SendDP(0,input)

            case "/change":
                name := cmds[1]
                change(cid,client,name)

            case "/quit\n":
                client.Clients.Get(*cid).GetTransport().SendDP(0, []byte("/quit"))

            default:
                client.Clients.Get(*cid).GetTransport().SendDP(0,input[0:len(input)-1])
            }
        } else {
            client.Clients.Get(*cid).GetTransport().SendDP(0,input[0:len(input)-1])
        }
    }
}

func change(cid *int,client *lnet.ClientPool,name string,) {
    b:= client.Clients.GetByName(name)
    if b!=nil{
        _cid := b.GetTransport().Cid
        *cid = _cid
        fmt.Println("current connection change:")
    }

    for c,p:=range client.Clients.All() {
        if p.GetName() != name {
            fmt.Println(" ",c,p.GetName())
        } else {
            fmt.Println("*",c,p.GetName())
        }
    }
}

func main() {
    datagram := &lib.Datagram{}

    cid := 0
    client := lnet.NewClientPool(MakeClient, datagram)
    go clientsender(&cid,client)

    //client.Start("", 4444)

    running :=1
    for running==1 {
        time.Sleep(3*1e3)
    }
}
