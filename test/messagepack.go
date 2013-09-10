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
    //"time"
    //"errors"
    //"runtime"
    "github.com/ugorji/go/codec"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
)

var (
        v interface{} // value to decode/encode into
        mh codec.MsgpackHandle
    )

func main() {
    a := map[string]interface{}{
        "I":23,
        "J":2345,
        "k":"234234",
    }

    test(a)
    b := map[string]interface{}{
        "I":23,
        "J":2345,
        "k":"adb32",
        "M":[]map[string]int {
            { "K":23, "L":2345 },
            { "K":354, "L":2345 },
            { "K":345, "L":2345 },
        },
    }
    test(b)


    json := &LGJsonSerialize{}
    vv,_ := json.Serialize(b)
    fmt.Println("len(json)",len(vv),vv)


    msgw := LGNewMessageWriter(1)
    msgw.WriteU(23,2345,"234234")

    b4 := LGNewMessageListWriter(1)
    b4.WriteStartTag()
    b4.WriteUint(23,0)
    b4.WriteUint(2345,0)
    b4.WriteEndTag()

    b4.WriteStartTag()
    b4.WriteUint(345,0)
    b4.WriteUint(2323,0)
    b4.WriteEndTag()

    b4.WriteStartTag()
    b4.WriteUint(345,0)
    b4.WriteUint(2323,0)
    b4.WriteEndTag()

    msgw.WriteList(b4,0)
    msgw.SetCode(1,1)
    data := msgw.ToBytes()
    fmt.Println("len(data)",len(data),data)
}

func test(a interface{}) {
    fmt.Println(a)
    b := []byte{}
    //enc = codec.NewEncoder(w, &mh)
    enc := codec.NewEncoderBytes(&b, &mh)
    err := enc.Encode(a)
    if err == nil {
        fmt.Println("len(b)=%d,v=%v",len(b),b)
    } else {
        fmt.Println("err=",err)
        return
    }

    var v interface{} // value to decode/encode into
    //dec = codec.NewDecoder(r, &mh)
    dec := codec.NewDecoderBytes(b, &mh)
    err = dec.Decode(&v)

    fmt.Println("v=",v,"\n\n")

}



