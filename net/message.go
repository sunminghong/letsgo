/*=============================================================================
#     FileName: message.go
#         Desc: MessageWriter pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-15 17:57:39
#      History:
=============================================================================*/
package net

import (
    //"encoding/binary"
)

const (
    TY_UINT = 0
    TY_STRING = 1
    TY_INT = 2
    TY_LIST = 3
    TY_UINT16 = 4
    TY_UINT32 = 5
    TY_INT32 = 6
)

type MessageWriter struct {
    Code uint16
    Ver byte

    // data item buff
    buf *RWStream

    //meta data buff
    metabuf *RWStream

    //meta data write item current index
    wind int

    meta map[int]byte
    //items map[int]interface{}

    maxItem int
}

func NewMessageWriter() *MessageWriter {
    m := &MessageWriter{}
    m.Init()
    return m
}

//对数据进行拆包
func (msg *MessageWriter) Init() {
    msg.meta = make(map[int]byte)
    msg.buf = NewRWStream(128,BigEndian)
    msg.maxItem = 0
    msg.Code = 0
    msg.Ver = 0
    msg.wind = 0

    msg.metabuf = NewRWStream(30,BigEndian)

    //leave 4 bytes to head(code,ver,metaitemdata)
    msg.metabuf.Write([]byte{0,0,0,0})
}

func (msg *MessageWriter) preWrite(wind int) {
    if wind == 0 {
        return
    }
    if wind < msg.maxItem{
        panic("item write order is wrong!")
    }
    msg.maxItem = wind
}
func (msg *MessageWriter) WriteUint16(x uint16,wind int){
    msg.preWrite(wind)

    msg.buf.WriteUint16(uint16(x))
    msg.metabuf.WriteByte(byte((msg.maxItem << 3) | TY_UINT16))
    msg.wind ++
    msg.maxItem ++
}

func (msg *MessageWriter) WriteUint32(x uint32,wind int){
    msg.preWrite(wind)

    msg.buf.WriteUint32(uint32(x))
    msg.metabuf.WriteByte(byte((msg.maxItem << 3) | TY_UINT32))
    msg.wind ++
    msg.maxItem ++
}

func (msg *MessageWriter) WriteUint(x uint,wind int){
    msg.preWrite(wind)

    msg.buf.WriteUint(uint(x))
    msg.metabuf.WriteByte(byte((msg.maxItem << 3) | TY_UINT))
    msg.wind ++
    msg.maxItem ++
}

func (msg *MessageWriter) WriteInt(x int,wind int){
    msg.preWrite(wind)

    msg.buf.WriteInt(int(x))
    msg.metabuf.WriteByte(byte((msg.maxItem << 3) | TY_INT))
    msg.wind ++
    msg.maxItem ++
}

func (msg *MessageWriter) WriteString(x string,wind int){
    msg.preWrite(wind)

    msg.buf.WriteString(x)
    msg.metabuf.WriteByte(byte((msg.maxItem << 3) | TY_STRING))
    msg.wind ++
    msg.maxItem ++
}

func (msg *MessageWriter) WriteList(list *MessageListWriter,wind int) {
    msg.preWrite(wind)

    msg.buf.Write(list.ToBytes(0,0))
    msg.metabuf.WriteByte(byte((msg.maxItem << 3) | TY_LIST))
    msg.wind ++
    msg.maxItem ++
}

//write no sign interge
func (msg *MessageWriter) WriteU(x ...interface{}) {
    for _,v := range x {
        switch v.(type) {
        case uint:
            vv,_:= v.(uint)
            msg.WriteUint(uint(vv),0)
        case int:
            vv,_:= v.(int)
            if vv <0 {
                panic("WriteU only write > 0 integer")
            }
            msg.WriteUint(uint(vv),0)
        case uint32:
            vv,_:= v.(uint32)
            msg.WriteUint32(vv,0)
        case uint16:
            vv,_:= v.(uint16)
            msg.WriteUint16(vv,0)
        case string:
            vv,_:= v.(string)
            msg.WriteString(vv,0)
        case *MessageListWriter:
            vv,_:= v.(*MessageListWriter)
            msg.WriteList(vv,0)

        }
    }
}

// write sign number
func (msg *MessageWriter) Write(x ...interface{}) {
    for _,v := range x {
        switch v.(type) {
        case uint:
            vv,_:= v.(uint)
            msg.WriteInt(int(vv),0)
        case int:
            vv,_:= v.(int)
            msg.WriteInt(vv,0)
        case uint32:
            vv,_:= v.(uint32)
            msg.WriteUint32(vv,0)
        case uint16:
            vv,_:= v.(uint16)
            msg.WriteUint16(vv,0)
        case string:
            vv,_:= v.(string)
            msg.WriteString(vv,0)
        case *MessageListWriter:
            vv,_:= v.(*MessageListWriter)
            msg.WriteList(vv,0)

        }
    }
}

//对数据进行封包
func (msg *MessageWriter) ToBytes(code int,ver byte) []byte {
    msg.metabuf.SetPos(0)
    msg.buf.SetPos(0)
    //write heads
    _,heads := msg.metabuf.Read(4)
    msg.metabuf.Endianer.PutUint16(heads, uint16(code))
    heads[2] = ver

    Log("wind:",msg.wind)
    heads[3] = byte(msg.wind)
    Log("metabuf",msg.metabuf.Bytes())
    msg.metabuf.Write(msg.buf.Bytes())

    Log("metabuf",msg.metabuf.Bytes())
    return msg.metabuf.Bytes()
}

/////////////////////////////////////////////////////////////////////////////////

type MessageReader struct {
    Code uint16
    Ver byte

    // data item buff
    buf *RWStream

    //meta data write item current index
    wind int

    meta map[int]byte
    //items map[int]interface{}

    maxItem int
}

func NewMessageReader(data []byte) *MessageReader{
    m := &MessageReader{}
    m.Init(data)
    return m
}

//对数据进行拆包
func (msg *MessageReader) Init(data []byte) {
    msg.buf = NewRWStream(data,BigEndian)
    buf := msg.buf

    code,_:= buf.ReadUint16()
    ver,_ := buf.ReadByte()

    msg.Code = code
    msg.Ver = ver

    _itemnum,_ := buf.ReadByte()
    itemnum := int(_itemnum)
    n,meta := buf.Read(itemnum)
    if n < itemnum {
        panic("data init error")
    }

    Log("init meta:",meta)
    maxitem := 0
    msg.meta = make(map[int]byte)

    for i:=0;i<itemnum;i++ {
        m := meta[i]
        ind := int(m>>3)
        if ind > maxitem {
            maxitem = ind
        }
        //msg.meta[ind] = (i<<3) |(m & 0x07)
        msg.meta[ind] = (m & 0x07)
    }
    msg.maxItem = maxitem
    msg.wind = 0
    Log(msg.meta)
}


func checkConvert(err error) {
    if err !=nil {
        panic("type cast failed!")
    }
}
/*
func (msg *MessageReader) preRead() {
    buf := msg.buf
    //data item meta data
    itemnum,_ = buf.ReadByte()
    items = make(map[int]interface{})
    msg.meta = make(map[byte]byte)
    for i:=0;i<itemnum;i++ {
        m := meta[i]
        msg.meta[m>>3] = m & 0x07

        switch m & 0x07 {
        case TY_UINT:
            v,err := buf.ReadUint()
            checkConvert(err)
            items[m>>3] = v
        case TY_INT:
            v,err := buf.ReadInt()
            checkConvert(err)
            items[m>>3] = v
        case TY_UINT16:
            v,err := buf.ReadUint16()
            checkConvert(err)
            items[m>>3] = v
        case TY_UINT32:
            v,err := buf.ReadUint32()
            checkConvert(err)
            items[m>>3] = v
        case TY_INT32:
            v,err := buf.ReadInt32()
            checkConvert(err)
            items[m>>3] = v
        case TY_STRING:
            v,err := buf.ReadString()
            checkConvert(err)
            items[m>>3] = v
        case TY_LIST:
            v := &MessageReaderList{}
            v.PreRead(buf)
            items[m>>3] = v
        }
    }
    msg.items = items
}
func (msg *MessageReader) ReadUint(wind int) uint {
    if len(msg.items) == 0 {
        msg.preRead()
    }

    v := msg.items[wind]
    a,ok := v.(uint)
    if !ok {
        panic("type cast failed!")
    }
    return uint(a),ok
}

func (msg *MessageReader) ReadInt(wind int) int {
    v := msg.items[wind]
    a,ok := v.(int)
    if !ok {
        panic("type cast failed!")
    }
    
    return a,ok
}
    m := msg.meta[wind]
    switch m & 0x07 {
    case TY_UINT:
    case TY_INT:
        a,ok := v.(int)
    case TY_UINT16:
        a,ok := v.(uint16)
    case TY_UINT32:
        a,ok := v.(uint32)
    case TY_INT32:
        return v.(int32)
    }

func (msg *MessageReader) alignPos(wind int) {
    for i:=msg.wind;i<wind;i++ {
        m := msg.meta[i]
        switch m & 0x07 {
        case TY_UINT:
            return int(buf.ReadUint())
        case TY_INT:
            return buf.ReadInt()
        case TY_UINT16:
            return int(buf.ReadUint16())
        case TY_UINT32:
            return int(buf.ReadUint32())
        case TY_INT:
            return int(buf.ReadInt32())
        }
    }
}
*/

func (msg *MessageReader) checkRead(datatype int) bool{
    Log("checkread wind,maxItem",msg.wind,msg.maxItem)
    if msg.wind > msg.maxItem {
        return false
    }

    ty,ok := msg.meta[msg.wind]
    Log("checkread ty,ok",ty,ok)
    if !ok {
        msg.wind ++
        return false
    }

    /////if (ty & 0x07) != TY_UINT{
    if ty != byte(datatype){
        panic("item data type that is reader is wrong")
    }
    return true
}

func (msg *MessageReader) ReadUint() uint {
    if msg.checkRead(TY_UINT) != true{
        return 0
    }

    v,err := msg.buf.ReadUint()
    checkConvert(err)
    msg.wind ++
    return v
}

func (msg *MessageReader) ReadInt() int {
    if msg.checkRead(TY_INT) != true{
        return 0
    }

    v,err := msg.buf.ReadInt()
    checkConvert(err)
    msg.wind ++
    return v
}

func (msg *MessageReader) ReadUint32() uint32 {
    if msg.checkRead(TY_UINT32) != true{
        return 0
    }

    v,err := msg.buf.ReadUint32()
    checkConvert(err)
    msg.wind ++
    return uint32(v)
}

func (msg *MessageReader) ReadUint16() uint16 {
    if msg.checkRead(TY_UINT16) != true{
        return 0
    }

    v,err := msg.buf.ReadUint16()
    checkConvert(err)
    msg.wind ++
    return uint16(v)
}

func (msg *MessageReader) ReadString() string{
    if msg.checkRead(TY_STRING) != true{
        return ""
    }

    v,err := msg.buf.ReadString()
    checkConvert(err)
    msg.wind ++
    return v
}

func (msg *MessageReader) ReadList() *MessageListReader{
    if msg.checkRead(TY_LIST) != true{
        return nil
    }

    return &MessageListReader{}
}

