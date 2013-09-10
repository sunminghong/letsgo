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
    "fmt"
    . "github.com/sunminghong/letsgo/helper"
)

func LGTest_MessageWrite(t *testing.T) {
    msgw := LGNewMessageWriter(LGBigEndian)
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

    c1 := 234
    c2 := 23425
    c3 :=1
    c4 := 4352

    //编码/写数据
    b4 := LGNewMessageListWriter(LGBigEndian)
    for i:=0;i<5;i++ {
        b4.WriteStartTag()

        b4.WriteUint(i,0)
        b4.WriteUint(i+1,0)
        b4.WriteUint(i+2,0)
        b4.WriteString(string(i+3),0)

        b4.WriteEndTag()
    }

    //msgw.Write(a1)
    msgw.Write(a1,a2,a3,a4,a5,a6,a7)
    msgw.WriteUint(int(b1),9)
    msgw.WriteU(b2)
    msgw.WriteU(b3)
    msgw.WriteList(b4,0)

    msgw.WriteUints(c1,c2,c3,c4)
    //fmt.Println("messageWrite",msgw.ToBytes(1,1))

    msgw.SetCode(1,1)
    data := msgw.ToBytes()

    //解码/读数据
    msg := LGNewMessageReader(data,LGBigEndian)

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
    if v5!= int(a5) {
        t.Error("item a5 ReadInt is wrong:",v5,a5)
    }

    v6 := msg.ReadUint16() 
    if v6!= int(a6) {
        t.Error("item a6 ReadInt is wrong:",v6,a6)
    }

    v7 := msg.ReadString() 
    if v7!= a7 {
        t.Error("item a7 ReadInt is wrong:",v7,a7)
    }
    _ = msg.ReadUint() 
    _ = msg.ReadUint() 

    vv1 := msg.ReadUint() 
    if vv1!= int(b1) {
        t.Error("item a1 ReadInt is wrong:",vv1,b1)
    }

    vv2 := msg.ReadUint() 
    if vv2!= int(b2) {
        t.Error("item a1 ReadInt is wrong:",vv2,b2)
    }

    vv3 := msg.ReadString() 
    if vv3!= b3 {
        t.Error("item a1 ReadInt is wrong:",vv3,b3)
    }

    fmt.Println("------------------------------------------------------")
    vv4 := msg.ReadList() 
    
    if vv4.Length != 5 {
        t.Error("item list Readlist length is wrong:",vv4.Length,5)
    }
    
    for i:=0;i<5;i++ {
        vv4.ReadStartTag()
        x := vv4.ReadUint()
        if x!= i {
            t.Error("item list(",i,",1) is wrong:",x,i)
        }
        x = vv4.ReadUint()
        if x!= (i+1) {
            t.Error("item list(",i,",2) is wrong:",x,i+1)
        }
        x = vv4.ReadUint()
        if x!= (i+2) {
            t.Error("item list(",i,",3) is wrong:",x,i+2)
        }
        x1 := vv4.ReadString()
        if x1!= string(i+3) {
            t.Error("item list(",i,",4) is wrong:",x1,i+3)
        }
        vv4.ReadEndTag()
    }


    if vv1!= int(c1) {
        t.Error("item a1 ReadInt is wrong:",vv1,c1)
    }

    vv2 = msg.ReadUint()
    if vv2!= int(c2) {
        t.Error("item a1 ReadInt is wrong:",vv2,c2)
    }

    vvv3 := msg.ReadUint()
    if vvv3!= c3 {
        t.Error("item a1 ReadInt is wrong:",vvv3,c3)
    }
    vvv4 := msg.ReadUint()
    if vvv4!= c4 {
        t.Error("item a1 ReadInt is wrong:",vvv4,c4)
    }

}
/*
type LGMessage struct {
    Code uint16
    Ver byte

    // data item buff
    buf *LGRWStream

    //meta data buff
    metabuf *LGRWStream

    //meta data write item current index
    wind int

    meta map[int]byte
    items map[int]interface{}

    maxItem int
}
*/
