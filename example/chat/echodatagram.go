/*=============================================================================
#     FileName: echodatagram.go
#         Desc: echo text server Datagram pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-07 18:15:55
#      History:
=============================================================================*/
package main

import (
    lnet "github.com/sunminghong/letsgo/net"
)

type EchoDatagram struct {

}


//对数据进行拆包
func (d *EchoDatagram) Fetch(c *lnet.Transport) (n int,msgs []*lnet.DataPacket) {
    msgs = []*lnet.DataPacket{}

    ilen := c.Stream.Len()
    if ilen == 0 {
        return
    }
    lnet.Log("Fetch",c.Stream.Bytes())
    msg := &lnet.DataPacket{Data: c.Stream.Bytes()}
    msgs = append(msgs,msg)
    n += 1

    //send to channel for consume
    c.InitBuff()

    return
}

//对数据进行封包
func (d *EchoDatagram) Pack(dp *lnet.DataPacket) []byte {
    return dp.Data
}
