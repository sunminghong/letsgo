/*=============================================================================
#     FileName: idassign_test.go
#         Desc: class with unix's process id alloc
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-05 14:54:19
#      History:
=============================================================================*/
package helper

import (
    "fmt"
    "testing"
    "math/rand"
//    "time"
)

var maxid = 1 << 16
var ida *IDAssign = NewIDAssign(maxid)


func log(v ...interface{}) {
    fmt.Println(v)
}

func logs(v ...interface{}) {
    fmt.Print(v[0])
}

func TestFree(t *testing.T) {
    log("colMask,lineMask:",colMask,lineMask)

    ida.setBit(30,1)
    if ida.test(30) != 1 {
        t.Error("setBit(30,1),but test(30)!=1")
    }
    ida.setBit(30,0)
    if ida.test(30) != 0 {
        t.Error("setBit(30,0),but test(30)!=0")
    }

    fmt.Println("test(32) = ",ida.test(32))

    ida.setBit(32,0)
    if ida.test(32) != 0 {
        t.Error("setBit(32,0),but test(32)!=0")
    }
}

func TestGetFreeID(t *testing.T) {
    ida.Init()
    counts := make(map[int]int)
    for i:=0;i<maxid;i++ {
        ida.GetFree()
    }


    f :=func(sid ,mid int) {
        for i:=sid;i<mid;i+=1 {
            //log("offset,value:",i+1,ida.test(i+1))
            //_id := rand.Intn(maxid)
            _id := i+1
            ida.Free(_id)
            //ida.free_(i+1)
            //log("changed:",i+1,ida.test(i+1))
            //if ida.test(_id) != 0 {
            //    t.Error("free() is error")
            //}
        }
    }
    go f(0,int(maxid /2) +1)
    go f(int(maxid /2),maxid )


    for i:=0;i<maxid;i++ {
        _id := ida.GetFree()

        v,ok := counts[_id]
        if ok {
            counts[_id] = v+1
        } else {
            counts[_id] = 1
        }
        if _id == 0 {
           //log("no freeid:",_id,i)
        } else {
            fmt.Print(_id," ")
        }

        if rand.Intn(100) < 50 {
            _id := rand.Intn(maxid)
            ida.Free(_id)
        }
    }

    for i:=1;i< maxid;i++ {
        v,ok := counts[i]
        if !ok {
            log("%D count is not exists",i)

            continue
        }
        if v!=1 {
            log("%D count is %d",i,v)
        }
    }
}
