/*=============================================================================
#     FileName: messagelist.go
#         Desc: Message pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-15 12:13:58
#      History:
=============================================================================*/
package net

import (
    //"encoding/binary"
)

type MessageListWriter struct {
    *MessageWriter

    length int

    meta []byte
}

type MessageListReader struct {
    *MessageReader

    length int

    meta []byte
}
/*
//对数据进行拆包
func (msg *MessageList) PreRead(data interface{}){
    Log("messagelist is called")
    if len(params) > 0 {
        msg.buf := NewRWStream(data,BigEndian)
        buf = msg.buf

        length,_:= buf.ReadUint()

        itemnum,_ = buf.ReadByte()
        msg.meta = buf.Read(itemnum)
        //meta[b>>3] = (b & 0x07)
    } else {
        msg.buf := NewRWStream(512,BigEndian)
        msg.metabuf := NewRWStream(30,BigEndian)
        msg.Code = 0
        msg.Ver = 0
        msg.wind = 0

        //leave 4 bytes to head(code,ver,metaitemdata)
        msg.metabuf.WriteBytes([]byte{0,0,0,0})
    }
}

func (msg *Message) WriteHead() {
    if msg.metabuf.Len() > 0 {
        return
    }
    
    msg.metabuf := NewRWStream(30,BigEndian)
    msg.metabuf.WriteBytes([]byte{0})
}

//对数据进行封包
func (msg *Message) ToBytes(code int,ver byte) []byte {
    //write heads
    heads := msg.metabuf.Read(4)
    binary.BigEndian.PutUint16(heads, uint16(code))
    heads[2] = ver

    heads[3] = byte(msg.wind)
    msg.metabuf.Write(msg.buf.Bytes())

    fmt.Pringln(msg.metabuf.Bytes())
    return msg.metabuf.Bytes()
}

*/
