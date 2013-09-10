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
package tmp

import (
    "fmt"
    //"time"
    //"errors"
    //"runtime"
    "testing"
    "strings"
    "github.com/ugorji/go/codec"
    . "github.com/sunminghong/letsgo/net"
    . "github.com/sunminghong/letsgo/helper"
)

var (
    mh codec.MsgpackHandle
)
/*
var a map[string]interface{}= map[string]interface{}{
    "I":23,
    "J":2345,
    "U":2345,
    "k":2345,
    "x":2345,
    "y":2345,
    "o":2345,
    "s":"234234",
}

*/
var b map[string]interface{} = map[string]interface{}{
    "I":23,
    "J":2345,
    "jj":2345,
    "NN":2345,
    "l":2345,
    "k":"adb32",
    "M":[]map[string]int {
        { "K":23, "L":2345 },
        { "K":354, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
        { "K":345, "L":2345 },
    },
}

type A struct {
    I,J,U,k,x,y,o int
    name string
}

var a *A
var av []byte
var bv []byte
var bvjson []byte
var avjson []byte

var json *LGJsonSerialize=&LGJsonSerialize{}

func init(){

    a = &A{23,2345,2345,2345,2345,2345,2345,"234234"}

    av = testencode(a)
    bv = testencode(b)
    avjson,_ = json.Serialize(a)
    bvjson,_ = json.Serialize(b)
    bvwrite := testencodewrite()

    fmt.Println("len(av)=",len(av))
    fmt.Println("len(avjson)=",len(avjson))

    s := string(avjson)
    s2 :=strings.Replace(s,"\"","",-1)
    fmt.Println("len(avjson2)=",len(s2))

    fmt.Println("len(bv)=",len(bv))
    fmt.Println("len(bvjson)=",len(bvjson))

    s3:=strings.Replace(string(bvjson),"\"","",-1)
    fmt.Println("len(bvjson2)=",len(s3))

    fmt.Println("len(bvwrie)=",len(bvwrite))
}

func Benchmark_smallSimple_msgpack_encode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        testencode(a)
    }
}

func Benchmark_smallSimple_msgpack_decode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        testdecode(av)
    }
}

func Benchmark_big_msgpack_encode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        testencode(b)
    }
}

func Benchmark_big_msgpack_decode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        testdecode(bv)
    }
}

func Benchmark_small_json_encode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        json.Serialize(a)
    }
}

func Benchmark_small_json_decode(t *testing.B) {
    var v interface{}
    for i := 0; i < t.N; i++ {
        json.Deserialize(avjson,&v)
    }
}

func Benchmark_big_json_encode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        json.Serialize(b)
    }
}

func Benchmark_big_json_decode(t *testing.B) {
    var v interface{}
    for i := 0; i < t.N; i++ {
        json.Deserialize(bvjson,&v)
    }
}

func Benchmark_message_encode(t *testing.B) {
    for i := 0; i < t.N; i++ {
        testencodewrite()
    }
}

func testencodewrite() []byte {
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

        b4.WriteStartTag()
        b4.WriteUint(345,0)
        b4.WriteUint(2323,0)
        b4.WriteEndTag()

        b4.WriteStartTag()
        b4.WriteUint(345,0)
        b4.WriteUint(2323,0)
        b4.WriteEndTag()

        b4.WriteStartTag()
        b4.WriteUint(345,0)
        b4.WriteUint(2323,0)
        b4.WriteEndTag()

        b4.WriteStartTag()
        b4.WriteUint(345,0)
        b4.WriteUint(2323,0)
        b4.WriteEndTag()

        b4.WriteStartTag()
        b4.WriteUint(345,0)
        b4.WriteUint(2323,0)
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
        return msgw.ToBytes()
}

func testencode(a interface{}) []byte {
    b := []byte{}
    //enc = codec.NewEncoder(w, &mh)
    enc := codec.NewEncoderBytes(&b, &mh)
    enc.Encode(a)

    return b

    //if err == nil {
    //    fmt.Println("len(b)=%d,v=%v",len(b),b)
    //} else {
    //    fmt.Println("err=",err)
    //    return
    //}
}

func testdecode(b []byte) {
    var v map[string]interface{} // value to decode/encode into
    //dec = codec.NewDecoder(r, &mh)
    dec := codec.NewDecoderBytes(b, &mh)
    _ = dec.Decode(&v)

    //fmt.Println("v=",v,"\n\n")
}



