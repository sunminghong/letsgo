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
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
    "strconv"
    "flag"

    "./protos"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

var endian = LGBigEndian

// clientsender(): read from stdin and send it via network
func clientsender(cid *int,client *LGClientPool) {
    reader := bufio.NewReader(os.Stdin)
    for {
        if (*cid)==0 {
            fmt.Print("you no connect anyone server,please input conn cmd,\n")
        }
        fmt.Print("you> ")
        input, _ := reader.ReadBytes('\n')
        cmd := string(input[:len(input)-1])

        var text string

        if cmd[0] == '/' {
            cmds := strings.Split(cmd," ")
            switch cmds[0]{
            case "/","/conn":
                ///conn s1 :12001 0
                var name,addr string
                var endian int

                if cmds[0] == "/" {
                    name= "s1"
                    addr=":12000"
                    endian=1
                } else {

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
                }

                if len(cmds)>3 {
                    endian,_= strconv.Atoi(cmds[3])
                }

                LGDebug("connect to server use endian:",endian)
                datagram := LGNewDatagram(endian)
                go client.Start(name,addr,datagram)


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

                text = string(input)
            case "/setlog":
                if lv,err := strconv.Atoi(cmds[1]); err == nil {
                    LGSetLevel(lv)
                }
            case "/change":
                name := cmds[1]
                change(cid,client,name)

            case "/quit\n":
                text = "/quit"

            default:
                text = string(input[:len(input)-1])
            }
        } else {
            text = string(input[:len(input)-1])
        }

        c := client.Clients.Get(*cid)
        msg := protos.NewMessageWriter(c)
        msg.SetCode(1011,0)
        msg.WriteString(text,0)

        LGTrace("has %v clients",client.Clients.Len())
        c.SendMessage(0,msg)
    }
}

func change(cid *int,client *LGClientPool,name string,) {
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

var (
    loglevel = flag.Int("loglevel",0,"log level")
)

func main() {
    flag.Parse()

    LGSetLevel(*loglevel)

    datagram := LGNewDatagram(protos.Endian)

    cid := 0
    client := LGNewClientPool(protos.NewClient, datagram)
    go clientsender(&cid,client)

    //client.Start("", 4444)

    quit := make(chan bool)
    <- quit
    
    //running :=1
    //for running==1 {
    //    time.Sleep(3*time.Second)
    //}
}

