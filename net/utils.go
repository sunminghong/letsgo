package net

import (
    //"fmt"
    //"bytes"
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

