/*=============================================================================
#     FileName: chatclient.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-12-10 11:34:28
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

    lnet "github.com/sunminghong/letsgo/net"

	"github.com/sbinet/liner"
    "./lib"
)

// LGIConnection  
type Connection struct {
    Transport *lnet.Transport
    Name      string
    Username    *string
}

func LGMakeConnection(name string,transport *lnet.Transport) lnet.LGIConnection {
    username := "someone"
    c := &Connection{transport, name, &username}

    return c
}

func (c *Connection) GetName() string {
    return c.Name
}

func (c *Connection) ProcessDPs(dps []*lnet.LGDataPacket) {
    for _, dp := range dps {
        md := string(dp.Data)

        fmt.Println()
        fmt.Println(md)
        fmt.Print("you> ")
    }
}

//对数据进行拆包
func (c *Connection) GetTransport() *lnet.Transport {
    return c.Transport
}

func (c *Connection) Close() {
    c.Transport.Close()
}

func (c *Connection) Closed() {
    msg := "system: " + (*c.Username) + " is leave!"
    c.Transport.SendBroadcast([]byte(msg))
}

func (c *Connection) SendMessage(msg lnet.LGIMessageWriter) {
    c.Transport.SendDP(0,msg.ToBytes())
}

func (c *Connection) SendBroadcast(msg lnet.LGIMessageWriter) {
    c.Transport.SendBroadcast(msg.ToBytes())
}


func tabCompleter(line string) []string {
	opts := make([]string, 0)

	if strings.HasPrefix(line, "/") {
		filters := []string{
			"/conn ",
			"/change ",
			"/quit",
			"/reg ",
			"/rereg ",
		}

		for _, cmd := range filters {
			if strings.HasPrefix(cmd, line) {
				opts = append(opts, cmd)
			}
		}
	}

	return opts
}

// clientsender(): read from stdin and send it via network
func clientsender(cid *int,client *lnet.ConnectionPool) {
	term := liner.NewLiner()
	fmt.Println("Skynet Interactive Shell")

	term.SetCompleter(tabCompleter)
    for {
        if (*cid)==0 {
            fmt.Print("you no connect anyone server,please input conn cmd,\n")
        }
		input, e := term.Prompt("> ")
		if e != nil {
			break
		}

        //cmd := string(input[:len(input)-1])
        cmd := string(input)
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

                p := client.Connections.GetByName(name)
                if p != nil {
                    fmt.Println(name," is exists !")
                    continue
                }

                go client.Start(name,addr)


                fmt.Print("please input your name:")
                input, _ := reader.ReadBytes('\n')
                input =input[0:len(input)-1]

                for true {
                    b := client.Connections.GetByName(name)
                    if b!=nil{
                        change(cid,client,name)
                        break
                    }
                    time.Sleep(2*1e3)
                }
                client.Connections.Get(*cid).GetTransport().SendDP(0,input)

            case "/change":
                name := cmds[1]
                change(cid,client,name)

            case "/quit\n":
                client.Connections.Get(*cid).GetTransport().SendDP(0, []byte("/quit"))

            default:
                client.Connections.Get(*cid).GetTransport().SendDP(0,input[0:len(input)-1])
            }
        } else {
            client.Connections.Get(*cid).GetTransport().SendDP(0,input[0:len(input)-1])
        }
    }
}

func change(cid *int,client *lnet.ConnectionPool,name string,) {
    b:= client.Connections.GetByName(name)
    if b!=nil{
        _cid := b.GetTransport().Cid
        *cid = _cid
        fmt.Println("current connection change:")
    }

    for c,p:=range client.Connections.All() {
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
    client := lnet.NewConnectionPool(MakeConnection, datagram)
    go clientsender(&cid,client)

    //client.Start("", 4444)

    running :=1
    for running==1 {
        time.Sleep(3*1e3)
    }
}
