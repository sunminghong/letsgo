/*=============================================================================
#     FileName: gate.go
#         Desc: game gate server
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:09:28
#      History:
=============================================================================*/
package main

import (
    "fmt"
    "time"
    "errors"
    "runtime"
)


type LGInterval struct {
    duration time.Duration
    quit chan bool
    stop chan bool
    callback func(self *LGInterval)

    isRun bool

    timer *time.Ticker

    isExecuting bool
}

func NewLGInterval(duration time.Duration,callbackfn func(self *LGInterval)) *LGInterval {
    return &LGInterval{duration,make(chan bool),make(chan bool),callbackfn,false,nil,false}
}

func (self *LGInterval) Start() error {
    if self.isRun {
        return errors.New("the instance is already running!")
    }

    go func() {
        fmt.Println("is starting...")
        self.isRun = true
        self.timer = time.NewTicker(self.duration)
        for {
            if !self.isRun {
                //<-self.quit
                goto quit
            }

            select {
            //case <-self.quit:
                //fmt.Println("quit")
                //goto quit
            case <-self.timer.C:
                self.isExecuting = true
                self.callback(self)
                self.isExecuting = false
            }
        }

        quit:
        self.isRun = false
        self.stop <- true
    }()

    return nil
}

var istop int =0
func (self *LGInterval) _stop(stoped chan bool,callback ...func()) {
    go func() {
        if self.isRun {
            self.isRun = false
            //self.quit<- true
            fmt.Println("enter _stop()")

            a := istop
            istop++
            //for self.isExecuting {
                //fmt.Println("self.isExecuting is true")
                //time.Sleep(100*time.Millisecond)
            //}
            fmt.Println("istop a:",a)
            <-self.stop
            self.timer.Stop()
            fmt.Println("istop a:",a)

            if len(callback)>0 {
                callback[0]()
            }

            fmt.Println("-----||```------")
        }
        fmt.Println("-----------")
        stoped <- true
    }()
}

func (self *LGInterval) Stop(callback ...func()) {
    stoped := make(chan bool)

    self._stop(stoped,callback...)
    //time.Sleep(100*Millsecond)
}

func (self *LGInterval) Reset(newDuration time.Duration) {
    go func() {
        fmt.Println("reset")

        stoped := make(chan bool)
        self._stop(stoped)
        <-stoped

        self.duration = newDuration
        self.Start()

        fmt.Println("reset2")
    }()
}

var lasttime time.Time =  time.Now()
var i int= 0
func callback(self *LGInterval) {
    fmt.Println(i,time.Now().Sub(lasttime))
    lasttime = time.Now()

    i++
    if i > 20 {
        self.Stop(func() {
            fmt.Println("stop")
            c <- true
        })
    } else if i==10 {
        self.Reset(500 * time.Millisecond)
    }

    time.Sleep(300*time.Millisecond)
}


var c chan bool = make(chan bool)
var duration=100 * time.Millisecond
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    interval := NewLGInterval(duration,callback)
    err := interval.Start()
    if (err!=nil) {
        fmt.Println("err:",err)
    }


    time.Sleep(5 * time.Second)
    //interval.Reset(1000 * time.Millisecond)
    fmt.Println("stop1")
    interval.Stop()

    fmt.Println("//-----------////")
    time.Sleep(5 * time.Second)
    interval.Start()

    fmt.Println("//////")
    fmt.Println("aaa")


    i=0

    time.Sleep(5 * time.Second)
    //interval.Reset(1000 * time.Millisecond)
    fmt.Println("stop3")
    interval.Stop()

    fmt.Println("//-----------////")
    time.Sleep(5 * time.Second)
    interval.Start()

    fmt.Println("//////")
    fmt.Println("aaa")

    <-c
    //fmt.Println(">>>>>>>>>>>>")
}

