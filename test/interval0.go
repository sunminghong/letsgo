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

func test() {

    c := LGNewRedis("192.168.18.18:6379")
    _, err := c.Connect(0)
    if err != nil {
        t.Errorf("Connect: %v", err)
        return
    }

    a1:= NewC()

    // Set
    err =c.Hmset("Hello", a1)
    if err != nil {
        t.Errorf("Set: %v", err)
        return
    }

}
