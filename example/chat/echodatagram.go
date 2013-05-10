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

type EchoDatagram struct {

}


//对数据进行拆包
func (d *EchoDatagram) Fetch(c *Transport) (n int,msgs []*DataPacket) {
    msgs = []*DataPacket{}

    ilen := len(c.Buff)
    if ilen == 0 {
        return
    }
    Log("Fetch",c.Buff)
    msg := &DataPacket{Data: c.Buff}
    msgs = append(msgs,msg)
    n += 1

    //send to channel for consume
    c.InitBuff()

    return
}

//对数据进行封包
func (d *EchoDatagram) Pack(dp *DataPacket) []byte {
    return dp.Data
}
