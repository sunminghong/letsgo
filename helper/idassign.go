/*=============================================================================
#     FileName: idassign.go
#         Desc: class with unix's process id alloc
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-05 10:46:10
#      History:
=============================================================================*/
package helper

import (
    "unsafe"
    "time"
//    "fmt"
)

//width: unsafe.Sizeof([]byte(nil)),
var (
    MaxID int = 32768
    bitsPerByte int = 8
    bytesPerUint int = int(unsafe.Sizeof(uint(1)))
    colMask int = int(bitsPerByte * bytesPerUint -1)
    //lineMask uint = uint(bytesPerUint + 1)
    lineMask uint = 6
)

// 分配一个唯一的ID，如clientID，pid
type LGIDAssign struct {
    maxid int
    bitsPerPageMask int

    //上次分配的ID
    lastid int

    //map页，用于指示某个ID是否分配
    page []uint

    free int

    //是否第一遍扫描
    first bool

    idChan chan int
    freeChan chan int
}

func LGNewIDAssign(maxid ...int) *LGIDAssign {
    _maxid := MaxID
    if len(maxid)>0 {
        _maxid = maxid[0]
    }

    ia := &LGIDAssign{}
    ia.maxid = _maxid

    ia.idChan = ia.getFreeChan()
    ia.freeChan = make(chan int)

    go func() {
        for {
            offset := <-ia.freeChan
            ia.free_(offset)
        }
    }()

    go func() {
        for {
            offset := <- ia.freeChan
            ia.free_(offset)
        }
    }()


    ia.Init()

    return ia
}

func (ia *LGIDAssign) Init() {
    ia.lastid = 0
    ia.free = ia.maxid
    ia.bitsPerPageMask = ia.maxid -1
    ia.first = true
    ia.page = make([]uint,int((ia.maxid)/bytesPerUint)+1)[:]
}

//分配一个ID，如果没有可分配的ID 了，就返回0
func (ia *LGIDAssign) GetFree() int {
    
    select {
    case _id := <- ia.idChan:
        return _id
    case <- time.After(1 * time.Second):
        return 0
    }

    return 0
}

//释放一个id
func (ia *LGIDAssign) Free(id int) {
    ia.freeChan <- id
}

func (ia *LGIDAssign) getFreeChan() (chan int) {
    out := make(chan int)
    go func() {
        for {
            var _id int
            _id = ia.getFree() 
            for _id ==0 {
                time.Sleep(200*time.Millisecond)
                _id = ia.getFree() 
            }

            out <- _id
        }
    }()
    return out
}


func (ia *LGIDAssign) free_(offset int) {
    if ia.test(offset) == 0 {
        return
    }
    ia.setBit(offset,0)
    ia.free ++
}

//设置offset位值，0或1'''
func (ia *LGIDAssign) setBit(offset,value int) {
    bit_off := uint(offset & colMask)
    int_off := offset >> lineMask

    //fmt.Println("offset,int,bit,value:",offset,int_off,bit_off,value,ia.page[int_off])
    if value == 1 {
        ia.page[int_off] |= (1 << bit_off)
    } else {
        ia.page[int_off] &= (^(1 << bit_off))
    }
    //fmt.Println("changed:",(ia.page[int_off] & (1 << bit_off))==0)
}

func (ia *LGIDAssign) test(offset int) int {
    bit_off := uint(offset & colMask)
    int_off := offset >> lineMask

    if (ia.page[int_off] & (1 << bit_off)) == 0 {
        return 0
    }
    return 1
}


//扫描map，返回一个为0的位'''
func (ia *LGIDAssign) findFree(offset int) int {
    size := ia.maxid
    page := ia.page
    for offset < size {
        bit_off := uint(offset & colMask)
        int_off := offset >> lineMask

        if (page[int_off] & (1<<bit_off)) != 0 {
            offset += 1
            continue
        }
        return offset
    }

    return -1
}

func (ia *LGIDAssign) getFree() int {
    if ia.free == 0 {
        return 0
    }

    if ia.first {
        ia.lastid ++
        if ia.lastid <= ia.maxid {
            ia.setBit(ia.lastid,1)
            ia.free --
            return ia.lastid
        }
        ia.lastid = 0
        ia.first = false
    }

    //fmt.Println("lastid:",ia.lastid,ia.bitsPerPageMask)
    lid := ia.lastid + 1
    offset := lid // & ia.bitsPerPageMask
    offset = ia.findFree(offset)
    if offset ==-1 {
        offset = ia.findFree(0)
        if offset == -1 {
            return 0
        }
    }

    ia.setBit(offset,1)
    ia.free --
    ia.lastid = offset

    return offset
}
