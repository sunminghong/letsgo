package net

import (
    "math/rand"
)


//add check code to old id
func LGGenerateID(oldid int) int {
    code := rand.Intn(1024)
    return oldid << 10 | int(code)
}

func LGCombineID(oldid int,code int) int {
    return oldid << 10 | code
}

func LGParseID(id int) (oldid int,checkcode int) {
    //return fromCid >> 10,fromCid & 3ff
    oldid = id >> 10
    checkcode = id - oldid
    return
}
