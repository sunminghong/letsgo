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
    "github.com/ugorji/go/codec"
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

    b := []byte
    //enc = codec.NewEncoder(w, &mh)
    enc = codec.NewEncoderBytes(&b, &mh)
    err = enc.Encode(a)
    if err == nil {
        fmt.Println("len(b)=%d,v=%v",len(b),b)
    } else {
        fmt.Println("err=",err)
        return
    }

    var v interface{} // value to decode/encode into
    //dec = codec.NewDecoder(r, &mh)
    dec = codec.NewDecoderBytes(b, &mh)
    err = dec.Decode(&v)

    fmt.Println("v=%v",v)

}
