package net

import (
    "time"
    "math/rand"
)


//add check code to old id
//oldid max value = 2097151 = 0x1fffff
func LGGenerateID(oldid int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
    code := r.Intn(1024)
    return oldid | int(code) << 21
}

func LGCombineID(oldid int,code int) int {
    return oldid | code << 21
}

func LGParseID(id int) (oldid int,checkcode int) {
    //return fromCid >> 10,fromCid & 3ff
    oldid = id & 0x1fffff  //(1 << 21 -1)
    checkcode = id >> 21
    return
}
