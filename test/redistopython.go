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
    . "github.com/sunminghong/letsgo/net/cache"
)

type Player struct {
    Id int `orm:"auto;pk"`

    AreaId int

    Name string `orm:"index"`
    Sex int

    Lv int `orm:"default(0)"`
    Money int `orm:"default(0)"`
    Gold int `orm:"default(0)"`

    Regtime int `orm:"auto_now_add"`
    Lasttime int `orm:"auto_now"`

    //don't save to db table
    Onlinetime int `orm:"-"`
}

func main() {

    c := LGNewRedis("192.168.18.18:6379")
    _, err := c.Connect(0)
    if err != nil {
        fmt.Println("Connect: %v", err)
        return
    }



    a1:= make(map[string]string)

    // Set
    err =c.Hgetall("player_1", &a1)
    if err != nil {
        fmt.Println("Set: %v", err)
        return
    }
    fmt.Println("player_1=",a1)

    a1 = make(map[string]string)
    err =c.Hgetall("datadict", &a1)
    if err != nil {
        fmt.Println("Set: %v", err)
        return
    }
    fmt.Println("datadict=",a1)

    b,err :=c.Get("key1")
    if err != nil {
		fmt.Println("Expecting %s, Received %s", err)
    }
    fmt.Println("b=",b)

    b,err =c.Get("key2")
    if err != nil {
		fmt.Println("Expecting %s, Received %s", err)
    }
    fmt.Println("b=",b)

    p := &Player{ 2,1,"abc",1,20,10000,100000,132432143,1243141234,2342 }
    err = c.Hmset("player_2",&p)
    if err != nil {
        fmt.Println("Set: %v", err)
        return
    }

    p = &Player{}
    err =c.Hgetall("player_2", &p)
    if err != nil {
        fmt.Println("Get: %v", err)
        return
    }
    fmt.Println("player_2=",p)

    fmt.Println(p.Name)

    p = &Player{}
    err =c.Hgetall("player_3", &p)
    if err != nil {
        fmt.Println("Get: %v", err)
        return
    }
    fmt.Println("player_3=",p)

    fmt.Println(p.Name)


    //list redis
    
    arr := []string{1,2,3,"ewr"}

    err = c.Lpush("list_1",[]byte(strings.Itoa(1)))
}

