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
    itemnum int

    meta []byte
}
func (list *MessageListWriter) Init() {
    Log("messagelistwriter Init by called")
    list.init(768)
    //list.MessageWriter.init(768)
}

func (msg *MessageListWriter) WriteEndTag() {
    if msg.itemnum ==0 {
        msg.itemnum = msg.maxInd
        msg.needWriteMeta = false
    } else if msg.maxInd != msg.itemnum {
        panic("write list item num is wrong")
    }

    msg.length ++

    msg.wind = 0
    msg.maxInd = 0
}

//对数据进行封包
func (msg *MessageListWriter) ToBytes() []byte {
    msg.metabuf.SetPos(0)
    msg.buf.SetPos(0)
    //write heads
    _,heads := msg.metabuf.Read(4)

    //write list bytes length
    msg.metabuf.Endianer.PutUint16(heads,
        uint16(msg.buf.Len()+msg.metabuf.Len() - 2))
    heads[2] = byte(msg.length)

    Log("wind:",msg.wind)
    heads[3] = byte(msg.wind)
    Log("metabuf",msg.metabuf.Bytes())

    msg.metabuf.Write(msg.buf.Bytes())

    Log("metabuf",msg.metabuf.Bytes())
    return msg.metabuf.Bytes()
}

type MessageListReader struct {
    *MessageReader

    //list length
    Length int

    //list byte length
    ByteLength int
}

//对数据进行拆包
func (msg *MessageListReader) Init(data []byte) {
    msg.buf = NewRWStream(data,BigEndian)
    buf := msg.buf

    byteLength,_:= buf.ReadUint16()
    length,_ := buf.ReadByte()

    msg.ByteLength = int(byteLength)
    msg.Length = int(length)

    msg.init()
}


func (msg *MessageListReader) ReadStartTag(data []byte) {
    if msg.wind == 0 {
        return
    }

    //对齐列表项，如果列表数据项比读取的多，读下一个列表的数据是需要先将指针对齐
    for i:=msg.wind;i<msg.maxInd;i++ {
        ty,ok := msg.meta[i]
        Log("checkread ty,ok",ty,ok)
        if !ok {
            continue
        }

        switch ty{
        case TY_UINT:
            msg.buf.ReadUint()
        case TY_INT:
            msg.buf.ReadInt()
        case TY_UINT16:
            msg.buf.ReadUint16()
        case TY_UINT32:
            msg.buf.ReadUint32()
        case TY_STRING:
            msg.buf.ReadString()
        case TY_LIST:
            ll,_ := msg.buf.ReadUint16()
            msg.buf.Read(int(ll))
        }
    }
    msg.wind = 0
}

func (msg *MessageListReader) ReadEndTag(data []byte) {
    msg.wind = 0

    //对齐列表项，如果列表数据项比读取的多，读下一个列表的数据是需要先将指针对齐
}

