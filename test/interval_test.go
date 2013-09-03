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

var lasttime time.Time =  time.Now()
var i int= 0
func callback(self *LGInterval) {
    fmt.Println(i,time.Now().Sub(lasttime))
    lasttime = time.Now()

    i++
    if i > 20 {
        self.Stop(func(interval *LGInterval) {
            fmt.Println("stop")
            c <- true
        })
    } else if i==10 {
        self.Stop(func(interval *LGInterval) {
            fmt.Println("reset")
            interval.Start(500 * time.Millisecond)
        })
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

