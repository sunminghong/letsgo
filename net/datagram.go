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
import (
    "encoding/binary"
    "github.com/sunminghong/letsgo/helper"
)

const (
    mask1 = byte(0x59)
    mask2 = byte(0x7a)
)

type Datagram struct {
    endian int
    Endianer helper.ItoB
}

func NewDatagram(endian int ) *Datagram{
    dg := &Datagram{}

    dg.SetEndian(endian)
    return dg
}

func (d *Datagram) GetEndian() int {
    return d.endian
}

func (d *Datagram) Clone(endian int) IDatagram {
    dg := &Datagram{}

    dg.SetEndian(endian)
    return dg
}

func (d *Datagram) SetEndian(endian int) {
    d.endian = endian
    if endian == helper.BigEndian {
        d.Endianer = binary.BigEndian
    } else {
        d.Endianer = binary.LittleEndian
    }
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


//flag1(byte)+flag2(byte)+datatype(byte)+data(datasize(int32)+body)+fromcid(int16)
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
        //Log("pos:",pos)

        //拆包
        if c.DPSize > 0 {
            if ilen-pos < c.DPSize {
                //如果缓存去数据长度不够就退出接着等后续数据
                return
            }
        } else {
            //Log("ilen,pos:",ilen,pos)
            if ilen-pos < 7 {
                return
            }

            heads,_ := cs.Read(7)
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

                if dataType == DATAPACKET_TYPE_GENERAL{
                    dpSize = int(_dpSize)
                } else {
                    dpSize = int(_dpSize) + 2
                }

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

        data,size := cs.Read(dpSize)
        if size > 0 {
            dp := &DataPacket{Type:dataType}

            if dataType != DATAPACKET_TYPE_GENERAL {
                dp.FromCid = int(d.Endianer.Uint16(data[dpSize-2:]))
                dp.Data = data[:dpSize-2]
            } else {
                dp.Data = data
            }

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
func (d *Datagram) Pack__(dp *DataPacket) []byte {
    ilen := len(dp.Data)
    if (dp.Type != DATAPACKET_TYPE_GENERAL) {
        ilen += 2
    }
    buf := make([]byte, ilen+7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)

    d.Endianer.PutUint32(buf[3:], uint32(ilen-2))

    d.encrypt(buf)

    copy(buf[7:], dp.Data)

    if (dp.Type != DATAPACKET_TYPE_GENERAL) {
        d.Endianer.PutUint16(buf[5+ilen:], uint16(dp.FromCid))
    }
    return buf
}


//对数据进行封包
func (d *Datagram) PackWrite(write WriteFunc,dp *DataPacket) []byte {
    buf := make([]byte,7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)

    ilen := len(dp.Data)
    d.Endianer.PutUint32(buf[3:], uint32(ilen))

    d.encrypt(buf)

    write(buf)
    write(dp.Data)

    if (dp.Type == DATAPACKET_TYPE_DELAY) {
        cid := make([]byte,2)
        d.Endianer.PutUint16(cid, uint16(dp.FromCid))
        write(cid)
    }

    return buf
}

