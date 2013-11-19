package net

import (
    "math/rand"
    "time"
)

//add check code to old id
//oldid max value = 2097151 = 0x1fffff
func LGGenerateID(oldid int) int {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    //code := r.Intn(1024)
    //return oldid | int(code)<<21
    code := r.Intn(99999)
    return code + oldid * 100000
}

func LGCombineID2(oldid int, code int) int {
    //return oldid | code<<21
    return code + oldid * 100000
}

func LGParseID(id int) (oldid int, checkcode int) {
    //return fromCid >> 10,fromCid & 3ff
    //oldid = id & 0x1fffff //(1 << 21 -1)
    //checkcode = id >> 21

    oldid = int(id / 100000)
    checkcode = id % 100000
    return
}
