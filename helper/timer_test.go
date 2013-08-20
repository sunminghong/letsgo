/*=============================================================================
#     FileName: timer_test.go
#         Desc: class with unix's process id alloc
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-08-15 11:24:50
#      History:
=============================================================================*/
package helper

import (
    "time"
    "testing"
    "fmt"
)

func Testtime(t *testing.T) {
    t2 := time.Now().Unix()
    fmt.Println("now timestamp is %d",t2)

    t3 := time.Now()

    t4 := time.Date(2013,8,1,0,0,0,0,time.Local)

    t5 := int(t3.Sub(t4).Seconds())
    t6 := LGNetTimestamp(t3)

    if t5 != t6 {
        t.Error("nettimestamp is error:",t5,t6)
    }

    d := "2013-07-13 03:20:33"
    t7,err := LGStrttime(d)
    if err !=nil {
        t.Error("LGStrttime is error:",err)
    }

    t8 := LGStrftime(t7)
    if t8 != d {
        t.Error("LGStrftime is not equ LGStrttime",t7,t8)
    }

    //t9 := LGStrftime(t2)
    //if t9 !=  {
    //    t.Error("LGStrftime is error:",t9,t2)
    //}

    d1 := "2013-07-13"
    t10 := LGTodayUnix(t7)
    t11,err := LGStrttime(d1,"2013-02-01")
    if err !=nil {
        t.Error("LGStrttime is error:",err)
    }

    if t10 != int(t11.Unix()) {
        t.Error("LGTodayUnix is error:",t10,t11)
    }

    t12 := LGToday(t7)
    if t12 != 20130713 {
        t.Error("LGToday is error:",t12,20130713)
    }

    t12 = LGToday()
    if t12 != 20130815 {
        t.Error("LGToday is error:",t12,20130815)
    }


    t12 = LGYesterday(t7)
    if t12 != 20130712 {
        t.Error("LGToday is error:",t12,20130712)
    }

    t12 = LGYesterday()
    if t12 != 20130814 {
        t.Error("LGToday is error:",t12,20130814)
    }

}
