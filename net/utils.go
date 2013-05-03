package net

import (
    "fmt"
    "bytes"
    "encoding/binary"
)

func Int32ToBytes(i int32) []byte {
    var buf = make([]byte, 8)
    binary.BigEndian.PutUint32(buf, uint32(i))
    return buf
}


func BytesToInt32(buf []byte) int32 {
    return int32(binary.BigEndian.Uint32(buf))
}


//add buff cap
func BytesGrow(buff []byte,addlen int) int {
    m := len(buff)
    if m + addlen > cap(buff) {
        var b_ []byte
        // not enough space anywhere
        b_ = make([]byte,m+addlen)
        copy(b_, buff)
        buff = b_
    }
    return m
}

// Write appends the contents of p to the []byte.  The return
// value n is the length of p; err is always nil.
// If the buffer becomes too large, Write will panic with
// ErrTooLarge.
func BytesAppend(buff []byte,p []byte) (n int, err error) {
    m := Grow(buff,len(p))
    return copy(buff[m:], p)
}

