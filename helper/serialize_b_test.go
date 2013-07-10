/*=============================================================================
#     FileName: serizlize_b_test.go
#         Desc: memcache test struct set/get
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-10 11:48:55
#      History:
=============================================================================*/

package helper

import (
    "testing"
    "fmt"
)

func Benchmark_SerializeGob(t *testing.B) {
    //fmt.Println("////////////////////////test serialize b //////////////////////////////")
    for i := 0; i < t.N; i++ {
        test()
    }
    bytelength := test()
    fmt.Println("bytelengthGob=",bytelength)
}

func test() int {
    bl := 0
    gs := &LGGobSerialize{}

    v1 := 23422
    v2 := -1

    v,_:= gs.Serialize(v1)
    //bl += len(v)
    _ = gs.Deserialize(v,&v2)

    v1 = 0
    v2 = -2
    v,_ = gs.Serialize(v1)
    //bl += len(v)
    _ = gs.Deserialize(v,&v2)

    a1,_ := NewA1A2()

    a2 := &A{}
    v,_ = gs.Serialize(a1)
    bl += len(v)
    _ = gs.Deserialize(v,&a2)

    return bl
}

func Benchmark_Serialize_Json(t *testing.B) {
    //fmt.Println("////////////////////////test serialize b //////////////////////////////")
    for i := 0; i < t.N; i++ {
        test2()
    }

    bytelength := test2()
    fmt.Println("bytelengthJson=",bytelength)
}

func test2() int {
    bl := 0
    gs := &LGJsonSerialize{}

    v1 := 23422
    v2 := -1

    v,_:= gs.Serialize(v1)
    //bl += len(v)
    _ = gs.Deserialize(v,&v2)

    v1 = 0
    v2 = -2
    v,_ = gs.Serialize(v1)
    //bl += len(v)
    _ = gs.Deserialize(v,&v2)

    a1,_ := NewA1A2()

    a2 := &A{}
    v,_ = gs.Serialize(a1)
    bl += len(v)
    _ = gs.Deserialize(v,&a2)

    return bl
}

