/*=============================================================================
#     FileName: byteutils.go
#         Desc: byte int convert helper or utils
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-09 17:23:02
#      History:
=============================================================================*/
package net

import (
    //"fmt"
    //"bytes"
    "encoding/binary"
)

type RWStream struct {
    buf []byte
    startpos int
    off
}


//add buff cap
func (c *RWStream)grow(addlen int) int{
    if pos + addlen > cap(rw.buf) {
        var b_ []byte
        b_ = make([]byte,pos+addlen)
        copy(b_, rw.buf)
        c.buf = b_
    }
    return m
}

// Write appends the contents of p to the []byte.  The return
// value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with
// ErrTooLarge.
func (c *Transport) bufAppend(p []byte) (n int) {
    Log("len(buff)=",len(c.buf),"len(p)=",len(p))
    m:= c.buffGrow(len(p))
    Log("buff",c.buf)
    a := copy((c.buf)[m:], p)

    Log("buff2",a,c.buf)
    return a
}

func ( rw *RWStream) write(p []byte) {

}

func Int32ToBytes(i int32) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint32(buf, uint32(i))
    return buf
}

func BytesToInt32(buf []byte) int32 {
    return int32(binary.BigEndian.Uint32(buf))
}


