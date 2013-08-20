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
    . "github.com/sunminghong/letsgo/helper"
//    . "github.com/sunminghong/letsgo/log"
)

const (
    mask1 = byte(0x59)
    mask2 = byte(0x7a)
)

type LGDatagram struct {
    endian int
    Endianer LGItoB
}

func LGNewDatagram(endian int ) *LGDatagram{
    dg := &LGDatagram{}

    dg.SetEndian(endian)
    return dg
}

func (d *LGDatagram) GetEndian() int {
    return d.endian
}

func (d *LGDatagram) Clone(endian int) LGIDatagram {
    dg := &LGDatagram{}

    dg.SetEndian(endian)
    return dg
}

func (d *LGDatagram) SetEndian(endian int) {
    d.endian = endian
    if endian == LGBigEndian {
        d.Endianer = binary.BigEndian
    } else {
        d.Endianer = binary.LittleEndian
    }
}

func (d *LGDatagram) encrypt(plan []byte){
    return
    for i,_ := range plan {
        plan[i] ^= 0x37
    }
}

func (d *LGDatagram) decrypt(plan []byte){
    return
    for i,_ := range plan {
        plan[i] ^= 0x37
    }
}


//flag1(byte)+flag2(byte)+datatype(byte)+data(datasize(int32)+body)+fromcid(int32)
//对数据进行拆包
func (d *LGDatagram) Fetch(c *LGTransport) (n int, dps []*LGDataPacket) {
    dps = []*LGDataPacket{}

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
        //LGTrace("c.dpsize:",c.DPSize)
        if c.DPSize > 0 {
            if ilen-pos < c.DPSize {
                //如果缓存去数据长度不够就退出接着等后续数据
                return
            }
            dpSize = c.DPSize
            dataType = c.DataType
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
            //LGTrace("m1,m2",m1,m2)
            if m1==mask1 && m2==mask2 {
                dataType,_ = cs.ReadByte()
                _dpSize,err := cs.ReadUint32()
                //LGTrace("dataType,dpSize,endian",dataType,_dpSize,cs.Endian)
                if err != nil {
                    c.InitBuff()
                    c.DPSize = 0
                    c.DataType = 0
                    return 0,nil
                }

                dpSize = int(_dpSize)
                if dataType == LGDATAPACKET_TYPE_DELAY || dataType == LGDATAPACKET_TYPE_BROADCAST || dataType == LGDATAPACKET_TYPE_CLOSE || dataType == LGDATAPACKET_TYPE_CLOSED {
                    dpSize += 4
                }

                pos = cs.GetPos()
                //LGTrace("ilen,pos,dpSize",ilen,pos,dpSize)
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
            dp := &LGDataPacket{Type:dataType}

            switch dataType {
            case LGDATAPACKET_TYPE_DELAY,LGDATAPACKET_TYPE_BROADCAST,LGDATAPACKET_TYPE_CLOSE,LGDATAPACKET_TYPE_CLOSED:
                dp.FromCid = int(d.Endianer.Uint32(data[dpSize-4:]))
                dp.Data = data[:dpSize-4]
            default:
                dp.Data = data
            }

            dps = append(dps,dp)
            n += 1
        }

        c.DPSize = 0
        c.DataType = 0

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
func (d *LGDatagram) Pack(dp *LGDataPacket) []byte {
    ilen := len(dp.Data)
    if dp.Type == LGDATAPACKET_TYPE_DELAY || dp.Type == LGDATAPACKET_TYPE_BROADCAST || dp.Type == LGDATAPACKET_TYPE_CLOSE || dp.Type == LGDATAPACKET_TYPE_CLOSED {
        ilen += 4
    }
    buf := make([]byte, ilen+7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)

    d.Endianer.PutUint32(buf[3:], uint32(ilen-4))

    d.encrypt(buf)

    copy(buf[7:], dp.Data)

    if dp.Type == LGDATAPACKET_TYPE_DELAY || dp.Type == LGDATAPACKET_TYPE_BROADCAST || dp.Type == LGDATAPACKET_TYPE_CLOSE || dp.Type == LGDATAPACKET_TYPE_CLOSED {
        d.Endianer.PutUint32(buf[3+ilen:], uint32(dp.FromCid))
    }
    return buf
}

//对数据进行封包
func (d *LGDatagram) PackWrite(write LGWriteFunc,dp *LGDataPacket) {
    buf := make([]byte,7)

    buf[0] = byte(mask1)
    buf[1] = byte(mask2)
    buf[2] = byte(dp.Type)

    ilen := len(dp.Data)
    d.Endianer.PutUint32(buf[3:], uint32(ilen))

    d.encrypt(buf)

    write(buf)
    write(dp.Data)

    if dp.Type == LGDATAPACKET_TYPE_DELAY || dp.Type == LGDATAPACKET_TYPE_BROADCAST || dp.Type == LGDATAPACKET_TYPE_CLOSE || dp.Type == LGDATAPACKET_TYPE_CLOSED {
        cid := make([]byte,4)
        d.Endianer.PutUint32(cid, uint32(dp.FromCid))
        //LGTrace("packwrite():",cid)
        write(cid)
    }
}

