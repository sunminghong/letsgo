/*=============================================================================
#     FileName: message_b_test.go
#         Desc: Message pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-17 12:30:08
#      History:
=============================================================================*/
package net

import (
    "testing"
    "fmt"
    . "github.com/sunminghong/letsgo/helper"
)

func LGBenchmark_MessageWrite(t *testing.B) {
    for i := 0; i < t.N; i++ {
        test()
    }
}

func test() {
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
    //Log("messageWrite",msgw.ToBytes(1,1))

    msgw.SetCode(2,0)
    data := msgw.ToBytes()

    msg := LGNewMessageReader(data,LGBigEndian)

    v1 := msg.ReadInt() 
    if v1!= a1 {
        
    }

    v2 := msg.ReadInt() 
    if v2!= a2 {
        
    }

    v3 := msg.ReadInt() 
    if v3!= a3 {
        
    }

    v4 := msg.ReadInt() 
    if v4!= a4 {
        
    }

    v5 := msg.ReadUint32() 
    if v5!= int(a5) {
        
    }

    v6 := msg.ReadUint16() 
    if v6!= int(a6) {
        
    }

    v7 := msg.ReadString() 
    if v7!= a7 {
        
    }
    _ = msg.ReadUint() 
    _ = msg.ReadUint() 

    vv1 := msg.ReadUint() 
    if vv1!= int(b1) {
        
    }

    vv2 := msg.ReadUint() 
    if vv2!= int(b2) {
        
    }

    vv3 := msg.ReadString() 
    if vv3!= b3 {
        
    }

    fmt.Println("------------------------------------------------------")
    vv4 := msg.ReadList() 
    
    if vv4.Length != 5 {
        
    }
    
    for i:=0;i<5;i++ {
        vv4.ReadStartTag()
        x := vv4.ReadUint()
        if x!= i {
            
        }
        x = vv4.ReadUint()
        if x!= (i+1) {
            
        }
        x = vv4.ReadUint()
        if x!= (i+2) {
            
        }
        x1 := vv4.ReadString()
        if x1!= string(i+3) {
            
        }
        vv4.ReadEndTag()
    }


}
