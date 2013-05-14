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
    "strconv"
    "time"
)

// IProtocol  
type Client struct {
    Transport *lnet.Transport
    Name      *string
}

func MakeClient(transport *lnet.Transport) lnet.IProtocol {
    name := "someone"
    c := &Client{transport, &name}


    return c
}

func (c *Client) ProcessDPs(dps []*lnet.DataPacket) {
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
    msg := "system: " + (*c.Name) + " is leave!"
    c.Transport.SendBoardcast([]byte(msg))
}

// clientsender(): read from stdin and send it via network
func clientsender(cid *int,client *lnet.Client) {
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
                addr := cmds[1]
                oldnum := client.TransportNum

                go client.Start(addr,0)

                fmt.Print("please input your name:")
                input, _ := reader.ReadBytes('\n')
                input =input[0:len(input)-1]
                name2 := string(input)
                c.Name = &name2

                for true {
                    if client.TransportNum > oldnum  
                    && client.Protos.Get(client.LastCid) !=nil{
                        _cid = client.LastCid
                        cid = &_cid
                        break
                    }
                    time.Sleep(2*1e3)
                }
                c.GetTransport().SendDP(0,input)

            case "/change":
                fmt.Println(cmds[1])
                _cid,err := strconv.Atoi(cmds[1])
                if err !=nil {
                    fmt.Println("command format is wrong!")
                    continue
                }
                cid = &_cid

                for c,_:=range client.Protos.All() {
                    fmt.Println(c)
                }

            case "/quit\n":
                client.Protos.Get(*cid).GetTransport().SendDP(0, []byte("/quit"))
            default:
                client.Protos.Get(*cid).GetTransport().SendDP(0,input[0:len(input)-1])
            }
        } else {
            client.Protos.Get(*cid).GetTransport().SendDP(0,input[0:len(input)-1])
        }
    }
}

func main() {
    datagram := &lib.Datagram{}

    cid := 0
    client := lnet.NewClient(MakeClient, datagram)
    go clientsender(&cid,client)

    //client.Start("", 4444)

    running :=1
    for running==1 {
        time.Sleep(3*1e3)
    }
}
