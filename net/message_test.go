/*=============================================================================
#     FileName: message_test.go
#         Desc: Message pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-15 18:18:22
#      History:
=============================================================================*/
package net

import (
    "testing"
)

func Test_MessageWrite(t *testing.T) {
    msgw := NewMessageWriter()

    a1 := 989887834
    a2 := 243
    a3 := 3298374
    a4 := -432423423
    a5 := uint32(23)
    a6 := uint16(32234)
    a7 := "aasalfjnsaknhfaksdfashdr8o324rskjdfh8oq734tjkdfq9ytfhasdbhuewrq364tqfgeawgiruhsb njafeuaaa"

    b1 :=uint(32342334)
    b2 :=uint(42323499)
    b3 :="bsdbbbb"

    //msgw.Write(a1)
    msgw.Write(a1,a2,a3,a4,a5,a6,a7)
    msgw.WriteUint(b1,9)
    msgw.WriteU(b2,b3)
    //Log("messageWrite",msgw.ToBytes(1,1))

    data := msgw.ToBytes(1,1)

    msg := NewMessageReader(data)

    v1 := msg.ReadInt() 
    if v1!= a1 {
        t.Error("item a1 ReadInt is wrong:",v1,a1)
    }

    v2 := msg.ReadInt() 
    if v2!= a2 {
        t.Error("item a2 ReadInt is wrong:",v2,a2)
    }

    v3 := msg.ReadInt() 
    if v3!= a3 {
        t.Error("item a3 ReadInt is wrong:",v3,a3)
    }

    v4 := msg.ReadInt() 
    if v4!= a4 {
        t.Error("item a4 ReadInt is wrong:",v4,a4)
    }

    v5 := msg.ReadUint32() 
    if v5!= a5 {
        t.Error("item a5 ReadInt is wrong:",v5,a5)
    }

    v6 := msg.ReadUint16() 
    if v6!= a6 {
        t.Error("item a6 ReadInt is wrong:",v6,a6)
    }

    v7 := msg.ReadString() 
    if v7!= a7 {
        t.Error("item a7 ReadInt is wrong:",v7,a7)
    }
    _ = msg.ReadUint() 
    _ = msg.ReadUint() 

    vv1 := msg.ReadUint() 
    if vv1!= b1 {
        t.Error("item a1 ReadInt is wrong:",vv1,b1)
    }

    vv2 := msg.ReadUint() 
    if vv2!= b2 {
        t.Error("item a1 ReadInt is wrong:",vv2,b2)
    }

    vv3 := msg.ReadString() 
    if vv3!= b3 {
        t.Error("item a1 ReadInt is wrong:",vv3,b3)
    }


}
/*
type Message struct {
    Code uint16
    Ver byte

    // data item buff
    buf *RWStream

    //meta data buff
    metabuf *RWStream

    //meta data write item current index
    wind int

    meta map[int]byte
    items map[int]interface{}

    maxItem int
}
*/
