/*=============================================================================
#     FileName: gateserver.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-11-22 10:49:10
#      History:
=============================================================================*/

package gate

import (
    "net"
    "strconv"
    "time"
    //"bufio"
    //"os"
    "fmt"
    "github.com/sbinet/liner"
    iniconfig "github.com/sunminghong/iniconfig"
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
    . "github.com/sunminghong/letsgo/net"
    "strings"
)

//Dispatcher
type LGIDispatcher interface {
    Init()
    //Add(client Connection,protocols []int)
    Add(gridID int, messageCodes *string)
    Remove(gridID int)

    Dispatch(messageCode int) (gridID int, ok bool)
}

const (
    CONNECTION_STATE_FREE        int = iota
    CONNECTION_STATE_CONNECTTING
    CONNECTION_STATE_AUTORECONNECTION
)

type gridConf struct {
    name         string
    host         string
    messageCodes string
    endian       int
    state        int
    datagram     LGIDatagram
}

type LGGateServer struct {
    *LGServer

    Grids *LGConnectionPool

    Dispatcher LGIDispatcher

    //makeclient NewGateConnectionFunc
    gridConfs map[string]*gridConf

    interval *LGInterval
    autoReconnectDuration time.Duration
}

func (gs *LGGateServer) InitFromConfig(
    configfile string,
    newPlayerConnection LGNewConnectionFunc, datagram LGIDatagram,
    newGridConnection LGNewConnectionFunc, dispatcher LGIDispatcher) {

    c, err := iniconfig.ReadConfigFile(configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    section := "Default"

    logconf, err := c.GetString(section,"logConfigFile")
    if err != nil {
        logconf = ""
    }
    LGSetLogger(&logconf)

    //start grid service
    name, err := c.GetString(section, "name")
    if err != nil {
        LGError(err.Error())
        return
    }

    host, err := c.GetString(section, "host")
    if err != nil {
        LGError(err.Error())
        return
    }

    serverid, err := c.GetInt(section, "serverid")
    if err != nil {
        LGError(err.Error())
        return
    }

    maxConnections, err := c.GetInt(section, "maxConnections")
    if err != nil {
        maxConnections = 1000
    }

    endian, err := c.GetInt(section, "endian")
    if err == nil {
        datagram.SetEndian(endian)
    } else {
        datagram.SetEndian(LGLittleEndian)
    }

    autoDuration, err := c.GetInt(section, "autoReconnectDuration")
    if err != nil {
        autoDuration = 5
    }
    autoReconnectDuration := time.Duration(autoDuration) * time.Second

    gs.Init(name, serverid, host, maxConnections,
        newPlayerConnection, datagram, newGridConnection, dispatcher,autoReconnectDuration)

}

func (gs *LGGateServer) Init(
    name string, gateid int, host string, maxConnections int,
    newPlayerConnection LGNewConnectionFunc, datagram LGIDatagram,
    newGridConnection LGNewConnectionFunc, dispatcher LGIDispatcher,autoReconnectionDuration time.Duration) {

    gs.LGServer = LGNewServer(name, gateid, host, maxConnections, newPlayerConnection, datagram)

    gs.gridConfs = make(map[string]*gridConf)
    gs.Grids = LGNewConnectionPool(newGridConnection, datagram)

    gs.Dispatcher = LGNewDispatcher()

    gs.SetParent(gs)

    gs.autoReconnectDuration = autoReconnectionDuration
}

func (gs *LGGateServer) NewTransport(
    newcid int, conn net.Conn) *LGTransport {

    LGTrace("gs's newtransport is run")
    return LGNewTransport(newcid, conn, gs, gs.Datagram)
}

/*
func (gs *LGGateServer) Start(gridsconfigfile *string) {
    //parse config ini file
    gs.connectGrids(gridsconfigfile)
    gs.LGServer.Start()
}
*/

func (gs *LGGateServer) AddAutoReConnect(name string) {
    if v,ok := gs.gridConfs[name]; ok {
        v.state = CONNECTION_STATE_AUTORECONNECTION

        LGTrace("add autoReconn 2",name)
    }

    if gs.interval == nil {
        gs.interval = NewLGInterval(gs.autoReconnectDuration, gs.autoReConnectGrid)
        gs.interval.Start(gs.autoReconnectDuration)
    } else if !gs.interval.IsRun {
        gs.interval.Start(gs.autoReconnectDuration)
    }

}

func (gs *LGGateServer) autoReConnectGrid(interval *LGInterval) {
    cou := 0
    for name, v := range gs.gridConfs {
        LGTrace("autoReConnectGrid ps is name:", name, v.state)

        //c := gs.Grids.Connections.GetByName(name)
        //if c != nil {
            //LGError("auto reconnection error:this gridserver is connected:%s",name)
            //v.state = CONNECTION_STATE_CONNECTTING
            //continue
        //}

        if v.state != CONNECTION_STATE_AUTORECONNECTION {
            continue
        }
        LGTrace("auto reconnect to ", name, v.state)

        cou ++
        if gs.ConnectGrid(name, v.host, &v.messageCodes, v.datagram) {
            cou --
        }
    }
    if cou == 0 {
        interval.Stop()
    }
}

func (gs *LGGateServer) ReConnectGrids() {
    for name, v := range gs.gridConfs {
        LGTrace("ReConnectGrids ps is name:", name, v.state)

        c := gs.Grids.Connections.GetByName(name)
        if c != nil {
            continue
        }

        if v.state == CONNECTION_STATE_CONNECTTING {
            continue
        }

        gs.ConnectGrid(name, v.host, &v.messageCodes, v.datagram)
    }
}

func (gs *LGGateServer) ConnectGrids(configfile *string) {
    c, err := iniconfig.ReadConfigFile(*configfile)
    if err != nil {
        LGError(err.Error())
        return
    }

    //make some connection to game server
    for i := 1; i < 50; i++ {
        section := "GridServer" + strconv.Itoa(i)
        if !c.HasSection(section) {
            continue
        }

        enabled, err := c.GetBool(section, "enabled")
        if err == nil && !enabled {
            continue
        }

        gname, err := c.GetString(section, "name")
        if err != nil {
            //if err.Reason == iniconfig.SectionNotFound {
            //    break
            //} else {
            LGError(err.Error())
            //    continue
            //}
            break
        }

        host, err := c.GetString(section, "host")
        if err != nil {
            continue
        }

        gCodes, err := c.GetString(section, "process")
        if err != nil {
            gCodes = ""
        }

        endian, err := c.GetInt(section, "endian")

        gs.gridConfs[gname] = &gridConf{gname, host, gCodes, endian, CONNECTION_STATE_FREE, nil}

        if err == nil {
            da := gs.Datagram.Clone(endian)
            gs.gridConfs[gname].datagram = da
            gs.ConnectGrid(gname, host, &gCodes, da)
        } else {
            gs.ConnectGrid(gname, host, &gCodes, nil)
        }
    }
}

func (gs *LGGateServer) ConnectGrid(
    name string, host string, messageCodes *string, datagram LGIDatagram) bool {

    LGInfo("connect to grid:", name)

    pool := gs.Grids
    //todo: maybe add finish chan and remove time.Sleep
    go pool.Start(name, host, datagram)
    time.Sleep(2 * time.Second)

    LGTrace("clientpool:", pool.Connections.All())
    //if Pool don't find it ,then that is no success!
    c := pool.Connections.GetByName(name)
    if c == nil {
        LGError(host + " can't connect")
        return false
    }

    gs.gridConfs[name].state = CONNECTION_STATE_CONNECTTING
    //add dispatche
    gridID := c.GetTransport().Cid
    gs.Dispatcher.Add(gridID, messageCodes)

    LGInfo("be connected to grid ", name)
    return true
}

func tabCompleter(line string) []string {
    opts := make([]string, 0)

    if strings.HasPrefix(line, "/") {
        filters := []string{
            "/sendtoall ",
            "/setmax ",
            "/reconn",
            "/start",
            "/stop",
            "/restart",
            "/exit",
            "/quit",
        }

        for _, cmd := range filters {
            if strings.HasPrefix(cmd, line) {
                opts = append(opts, cmd)
            }
        }
    }

    return opts
}
func (gs *LGGateServer) StartConsole(quit chan bool) {
    term := liner.NewLiner()
    fmt.Println("gate server console")
    defer term.Close()

    term.SetCompleter(tabCompleter)
    //reader := bufio.NewReader(os.Stdin)
    for {
        input, e := term.Prompt(gs.Name + "> ")
        if e != nil {
            break
        }
        //input, _ := reader.ReadBytes('\n')
        //cmd := string(input[:len(input)-1])
        cmd := string(input)

        cmds := strings.Split(cmd, " ")
        switch cmds[0] {
        case "/sendtoall":
            ///conn s1 :12001 0
            if len(cmds) > 1 {
                msg := strings.Join(cmds[1:], " ")

                mw := LGNewMessageWriter(gs.Datagram.GetEndian())
                mw.SetCode(2011, 0)
                mw.WriteString(msg, 0)

                dp := &LGDataPacket{
                    Type:    LGDATAPACKET_TYPE_BROADCAST,
                    Data:    mw.ToBytes(),
                    FromCid: 0,
                }

                gs.SendBroadcast(dp)
            }

        case "/setmax":
            if len(cmds) > 1 {
                max, err := strconv.Atoi(cmds[1])
                if err != nil {
                    fmt.Println("setmax is error:", err)
                    continue
                }
                gs.SetMaxConnections(max)
            } else {
                fmt.Println("please input number of max connections")
            }
        case "/restart":
            gs.Stop()
            gs.Start()
        case "/start":
            gs.Start()
        case "/stop":
            gs.Stop()
        case "/reconn":
            gs.ReConnectGrids()

        case "/exit", "/quit":
            fmt.Println("this gateserver is exit")
            quit <- true
            break
        default:
        }
    }
}
