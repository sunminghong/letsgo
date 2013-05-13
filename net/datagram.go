/*=============================================================================
#     FileName: datagram.go
#         Desc: Datagram pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-06 14:35:37
#      History:
=============================================================================*/
package net

//type DatagramInterface interface {
//    Fetch(data []byte) int,[]byte
//    Pack(int,[]byte) []byte
//}

const (
    mask1 = byte(0x59)
    mask2 = byte(0x7a)

    DATAPACKET_TYPE_GENERAL = 0
    DATAPACKET_TYPE_DELAY = 1
    DATAPACKET_TYPE_BOARDCAST = 3
)

type Datagram struct {
    Stream RWStream     //can use enddian
}

func NewDatagram(endian int) *Datagram{
    dg := &Datagram{}
    dg.Stream = *NewRWStream([]byte{1},endian)

    return dg
}

func (d *Datagram) encrypt(plan []byte){
    for i,_ := range plan {
        plan[i] ^= 0x37
    }
}

func (d *Datagram) decrypt(plan []byte){
    for i,_ := range plan {
        plan[i] ^= 0x37
    }
}

//对数据进行拆包
func (d *Datagram) Fetch(c *Transport) (n int, dps []*DataPacket) {
    dps = []*DataPacket{}

    cs := c.Stream
    ilen := cs.Len()
    if ilen == 0 {
        return
    }

    var dpSize int

    var dataType,m1,m2 byte
    for {
        pos := cs.GetPos()
        Log("pos:",pos)

        //拆包
        if c.DPSize > 0 {
            if ilen-pos < c.DPSize {
                //如果缓存去数据长度不够就退出接着等后续数据
                return
            }
        } else {
            Log("ilen,pos:",ilen,pos)
            if ilen-pos < 7 {
                return
            }

            _,heads := cs.Read(7)
            d.decrypt(heads)

            cs.SetPos(-7)
            m1,_ = cs.ReadByte()
            m2,_ = cs.ReadByte()
            if m1==mask1 && m2==mask2 {
                dataType,_ = cs.ReadByte()
                _dpSize,err := cs.ReadUint32()
                if err != nil {
                    c.InitBuff()
                    c.DPSize = 0
                    c.DataType = 0
                    return 0,nil
                }

                dpSize = int(_dpSize)
                pos = cs.GetPos()
                if ilen - pos < dpSize {
                    c.DPSize = dpSize
                    c.DataType = dataType

                    return
                }

            } else {
                //如果错位则将缓存数据抛弃
                c.InitBuff()
                return
            }
        }

        size,data := cs.Read(dpSize)
        if size > 0 {
            dp := &DataPacket{Type:dataType, Data:data}
            dps = append(dps,dp)
            n += 1
        }

        c.DPSize = 0
        c.DataType = 0

        //send to channel for consume
        //c.ProcessDP(dp)

        if ilen - cs.GetPos() > 7 {
            continue
        } else {
            c.InitBuff()
            return
        }
    }
    return
}

//对数据进行封包
func (d *Datagram) Pack(dp *DataPacket) []byte {
    ilen := len(dp.Data)
    buf := make([]byte, ilen+7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)

    d.Stream.Endianer.PutUint32(buf[3:], uint32(ilen))

    d.encrypt(buf)

    copy(buf[7:], dp.Data)
    return buf
}
