/*=============================================================================
#     FileName: datagram_test.go
#         Desc: Datagram pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-09 16:07:54
#      History:
=============================================================================*/
package net

import (
    //"encoding/binary"
    "testing"
    "bytes"
    . "github.com/sunminghong/letsgo/helper"
)

/*

//datagram and datapacket define
type LGIDatagram interface {
    Fetch(c *LGTransport) (n int, dps []*LGDataPacket)
    Pack(dp *LGDataPacket) []byte
}
*/

func LGTest_Pack(t *testing.T) {
    datagram := LGNewDatagram(LGBigEndian)

    data := []byte("1234567890")
    dp := &LGDataPacket{Type:1,Data:data}
    data2 := datagram.Pack(dp)

    _data := [17]byte{0x59 ^ 0x37,0x7a ^ 0x37,1 ^ 0x37,0 ^ 0x37,0 ^ 0x37,0 ^ 0x37,10 ^ 0x37}
    copy(_data[7:],data)

    if !bytes.Equal(_data[:],data2) {
        t.Error("pack return is equal:",data2)
    }

}

func LGTest_Fetch(t *testing.T) {
    /*
  username, _ := getAdmin(1)
    if (username != "admin") {
         t.Error("getAdmin get data error")
    }
    */

    datagram := LGNewDatagram(LGBigEndian)
    trans := LGNewTransport(1,nil,nil,datagram)

    for ii:=0;ii<3;ii++ {
        if ii ==2 {
            trans.InitBuff()
        }

    buff := []byte{0x59 ^ 0x37,0x7a ^ 0x37,1 ^ 0x37,0 ^ 0x37,0 ^ 0x37,0 ^ 0x37,10 ^ 0x37}
    trans.BuffAppend(buff)

    data0 := []byte("1234567890")
    trans.BuffAppend(data0)

    data := trans.Stream.Bytes()
    trans.BuffAppend(data)
    trans.BuffAppend(data)
    trans.BuffAppend(data)
    trans.BuffAppend(data)

    n,dps := datagram.Fetch(trans)
    if n != 5 || len(dps)!= 5 {
        t.Error("fetch dps len is error:",n,len(dps))
    }

    dp := dps[0]
    if dp.Type != 1 {
        t.Error("fetch dps data is error")
    }

    if !bytes.Equal(dp.Data,data0) {
        t.Error("fetch dps data is error" + string(dp.Data))
    }

    dp = dps[4]
    if dp.Type != 1 {
        t.Error("fetch dps data is error",4)
    }

    if !bytes.Equal(dp.Data,data0) {
        t.Error("fetch dps data is error" ,4, string(dp.Data))
    }
}
}

