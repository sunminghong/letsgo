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
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

const (
    TY_UINT   = 0
    TY_STRING = 1
    TY_INT    = 2
    TY_LIST   = 3
    TY_UINT16 = 4
    TY_UINT32 = 5
)

type LGMessageWriter struct {
    Code int
    Ver  byte

    // data item buff
    buf *LGRWStream

    //meta data buff
    metabuf *LGRWStream

    meta map[int]byte
    //items map[int]interface{}

    //meta data write item current index
    wind   int
    maxInd int

    needWriteMeta bool

    alreadyToBytes bool
    outBuffer []byte
}

func LGNewMessageWriter(endian int) *LGMessageWriter {
    msg := &LGMessageWriter{}

    msg.init(128, endian)
    return msg
}

func (msg *LGMessageWriter) init(bufsize int, endian int) {
    msg.meta = make(map[int]byte)
    msg.buf = LGNewRWStream(bufsize, endian)
    msg.maxInd = 0
    msg.Code = 0
    msg.Ver = 0
    msg.wind = 0
    msg.needWriteMeta = true

    msg.metabuf = LGNewRWStream(30, endian)

    //leave 4 bytes to head(code,ver,metaitemdata)
    //leave 4 bytes to head(list length(uint16),list length(byte),metaitemdataLength(byte))
    msg.metabuf.Write([]byte{0, 0, 0, 0})

    msg.alreadyToBytes = false
    //msg.outBuffer = []byte{}
}

func (msg *LGMessageWriter) SetCode(code int, ver byte) {
    msg.Code = code
    msg.Ver = ver
}

func (msg *LGMessageWriter) GetCode() int {
    return msg.Code
}

func (msg *LGMessageWriter) preWrite(wind int) {
    msg.alreadyToBytes = false
    if wind == 0 {
        return
    }
    if wind < msg.maxInd {
        LGError("messageWriter.preWrite():item write order is wrong!")
        return
    }
    msg.maxInd = wind
}

func (msg *LGMessageWriter) writeMeta(datatype int) {
    if !msg.needWriteMeta {
        return
    }
    msg.metabuf.WriteByte(byte((msg.maxInd << 3) | datatype))
}

func (msg *LGMessageWriter) WriteUint16(x int, wind int) {
    if x < 0 {
        LGError("WriteUint16 only write > 0 integer");return
    }
    msg.preWrite(wind)

    //todo:if x=0 then dont't write
    msg.buf.WriteUint16(uint16(x))
    msg.writeMeta(TY_UINT16)
    msg.wind++
    msg.maxInd++
}

func (msg *LGMessageWriter) WriteUint32(x int, wind int) {
    if x < 0 {
        LGError("WriteUint32 only write > 0 integer");return
    }
    msg.preWrite(wind)

    msg.buf.WriteUint32(uint32(x))
    msg.writeMeta(TY_UINT32)
    msg.wind++
    msg.maxInd++
}

func (msg *LGMessageWriter) WriteUint(x int, wind int) {
    if x < 0 {
        LGError("WriteUint only write > 0 integer");return
    }
    msg.preWrite(wind)

    msg.buf.WriteUint(uint(x))
    msg.writeMeta(TY_UINT)
    msg.wind++
    msg.maxInd++
}

func (msg *LGMessageWriter) WriteUints(xs ...int) {
    for _,x := range xs {
        if x < 0 {
            LGError("WriteUint only write > 0 integer");return
        }
        msg.preWrite(0)

        msg.buf.WriteUint(uint(x))
        msg.writeMeta(TY_UINT)
        msg.wind++
        msg.maxInd++
        //msg.WriteUint(x, 0)
    }
}

func (msg *LGMessageWriter) WriteInt(x int, wind int) {
    msg.preWrite(wind)

    msg.buf.WriteInt(int(x))
    msg.writeMeta(TY_INT)
    msg.wind++
    msg.maxInd++
}

func (msg *LGMessageWriter) WriteInts(xs ...int) {
    for _,x := range xs {
        msg.preWrite(0)

        msg.buf.WriteInt(int(x))
        msg.writeMeta(TY_INT)
        msg.wind++
        msg.maxInd++
    }
}

func (msg *LGMessageWriter) WriteString(x string, wind int) {
    msg.preWrite(wind)

    msg.buf.WriteString(x)
    msg.writeMeta(TY_STRING)
    msg.wind++
    msg.maxInd++
}

func (msg *LGMessageWriter) WriteList(list *LGMessageListWriter, wind int) {
    msg.preWrite(wind)

    msg.buf.Write(list.ToBytes())
    msg.writeMeta(TY_LIST)
    msg.wind++
    msg.maxInd++
}

//write no sign interge
func (msg *LGMessageWriter) WriteU(x ...interface{}) {
    for _, v := range x {
        switch v.(type) {
        case uint:
            vv, _ := v.(uint)
            msg.WriteUint(int(vv), 0)
        case int:
            vv, _ := v.(int)
            if vv < 0 {
                LGError("WriteU only write > 0 integer");return
            }
            msg.WriteUint(int(vv), 0)
        case uint32:
            vv, _ := v.(uint32)
            msg.WriteUint32(int(vv), 0)
        case uint16:
            vv, _ := v.(uint16)
            msg.WriteUint16(int(vv), 0)
        case string:
            vv, _ := v.(string)
            msg.WriteString(vv, 0)
        case *LGMessageListWriter:
            vv, _ := v.(*LGMessageListWriter)
            msg.WriteList(vv, 0)

        }
    }
}

// write sign number
func (msg *LGMessageWriter) Write(x ...interface{}) {
    for _, v := range x {
        switch v.(type) {
        case uint:
            vv, _ := v.(uint)
            msg.WriteInt(int(vv), 0)
        case int:
            vv, _ := v.(int)
            msg.WriteInt(vv, 0)
        case uint32:
            vv, _ := v.(uint32)
            msg.WriteUint32(int(vv), 0)
        case uint16:
            vv, _ := v.(uint16)
            msg.WriteUint16(int(vv), 0)
        case string:
            vv, _ := v.(string)
            msg.WriteString(vv, 0)
        case *LGMessageListWriter:
            vv, _ := v.(*LGMessageListWriter)
            msg.WriteList(vv, 0)

        }
    }
}

//对数据进行封包
func (msg *LGMessageWriter) ToBytes() []byte {
    if msg.alreadyToBytes {
        return msg.outBuffer
    }

    if msg.Code == 0 {
        LGWarn("messagewriter ToBytes() msg.Code == 0")
        return nil
    }

    msg.metabuf.SetPos(0)
    msg.buf.SetPos(0)
    //write heads
    heads, _ := msg.metabuf.Read(4)
    msg.metabuf.Endianer.PutUint16(heads, uint16(msg.Code))
    heads[2] = msg.Ver

    //LGTrace("code:%d,wind:%d",msg.Code, msg.wind)
    heads[3] = byte(msg.wind)
    //LGTrace("code:%d,metabuf:% X", msg.Code,msg.metabuf.Bytes())

    bs := make([]byte,msg.metabuf.Len()+msg.buf.Len())
    copy(bs,msg.metabuf.Bytes())
    copy(bs[msg.metabuf.Len():],msg.buf.Bytes())
    msg.alreadyToBytes = true
    msg.outBuffer = bs

    //LGTrace("code:%d,metabuf:% X", msg.Code,bs)
    return bs
}

/////////////////////////////////////////////////////////////////////////////////

type LGMessageReader struct {
    Code int
    Ver  int

    endian int
    // data item buff
    buf *LGRWStream

    //meta data write item current index
    wind int

    meta map[int]byte
    //items map[int]interface{}

    maxInd  int
    itemnum int
}

func LGNewMessageReader(data []byte, endian int) *LGMessageReader {
    msg := &LGMessageReader{}

    msg.endian = endian
    msg.buf = LGNewRWStream(data, endian)
    buf := msg.buf

    code, _ := buf.ReadUint16()
    ver, _ := buf.ReadByte()

    msg.Code = int(code)
    msg.Ver = int(ver)

    msg.init()

    return msg
}

func (msg *LGMessageReader) init() {

    buf := msg.buf
    _itemnum, _ := buf.ReadByte()
    itemnum := int(_itemnum)
    meta, n := buf.Read(itemnum)
    if n < itemnum {
        LGError("messageReader data init error",n,itemnum,buf.Bytes())
        return
    }

    //LGTrace("init meta:", meta)
    maxind := 0
    msg.meta = make(map[int]byte)

    for i := 0; i < itemnum; i++ {
        m := meta[i]
        ind := int(m >> 3)
        if ind > maxind {
            maxind = ind
        }
        //msg.meta[ind] = (i<<3) |(m & 0x07)
        msg.meta[ind] = (m & 0x07)
    }
    msg.maxInd = maxind
    msg.itemnum = itemnum
    msg.wind = 0
    //LGTrace(msg.meta)
}

func checkConvert(err error) {
    if err != nil {
        LGError("messageReader type convert failed:",err)
    }
}

/*
func (msg *LGMessageReader) preRead() {
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
func (msg *LGMessageReader) ReadUint(wind int) uint {
    if len(msg.items) == 0 {
        msg.preRead()
    }

    v := msg.items[wind]
    a,ok := v.(uint)
    if !ok {
        LGError("type cast failed!");return
    }
    return uint(a),ok
}

func (msg *LGMessageReader) ReadInt(wind int) int {
    v := msg.items[wind]
    a,ok := v.(int)
    if !ok {
        LGError("type cast failed!");return
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

func (msg *LGMessageReader) alignPos(wind int) {
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

func (msg *LGMessageReader) checkRead(datatype int) bool {
    //LGTrace("checkread wind,maxInd", msg.wind, msg.maxInd)
    if msg.wind > msg.maxInd {
        return false
    }

    ty, ok := msg.meta[msg.wind]
    //LGTrace("checkread ty,ok", ty, ok, datatype)
    if !ok {
        msg.wind++
        return false
    }

    /////if (ty & 0x07) != TY_UINT{
    if ty != byte(datatype) {
        LGError("messageReader: item data type that is reader is wrong",msg.wind,datatype)
        return false
    }
    return true
}

func (msg *LGMessageReader) ReadUint() int {
    if msg.checkRead(TY_UINT) != true {
        return 0
    }

    v, err := msg.buf.ReadUint()
    checkConvert(err)
    msg.wind++
    return int(v)
}

func (msg *LGMessageReader) ReadInt() int {
    if msg.checkRead(TY_INT) != true {
        return 0
    }

    v, err := msg.buf.ReadInt()
    checkConvert(err)
    msg.wind++
    return v
}

func (msg *LGMessageReader) ReadUint32() int {
    if msg.checkRead(TY_UINT32) != true {
        return 0
    }

    v, err := msg.buf.ReadUint32()
    checkConvert(err)
    msg.wind++
    return int(v)
}

func (msg *LGMessageReader) ReadUint16() int {
    if msg.checkRead(TY_UINT16) != true {
        return 0
    }

    v, err := msg.buf.ReadUint16()
    checkConvert(err)
    msg.wind++
    return int(v)
}

func (msg *LGMessageReader) ReadString() string {
    if msg.checkRead(TY_STRING) != true {
        return ""
    }

    v, err := msg.buf.ReadString()
    checkConvert(err)
    msg.wind++
    return v
}

func (msg *LGMessageReader) ReadList() *LGMessageListReader {
    if msg.checkRead(TY_LIST) != true {
        return nil
    }

    list := LGNewMessageListReader(msg.buf)

    msg.wind++
    return list
}

func (msg *LGMessageReader) ReadCode() int {
    return msg.Code
}

