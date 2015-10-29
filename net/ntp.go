/*=============================================================================
#     FileName: ntp.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2014-09-26 14:22:31
#      History:
=============================================================================*/
package net
//import (
    //"fmt"
//)

type LGNtp1 struct {
    d1List []int
    d2List []int
    d3List []int

    needNumber int
    maxDelay int

    timeError int
}

func LGNewNtp1(needNumber, maxDelay int) *LGNtp1 {
    return &LGNtp1{
        d1List: []int{},
        d2List: []int{},
        d3List: []int{},

        needNumber: needNumber,
        maxDelay: maxDelay,
    }
}

func (this *LGNtp1) Init() {
    this.d1List = []int{}
    this.d2List = []int{}
    this.d3List = []int{}
} //init

//return ok=false, don't need more data
func (this *LGNtp1) Push(t1,t2,t3,t4 int) (ok bool) {
    if t2-t1 + t4-t3 > this.maxDelay {
        return true
    }

    this.d1List = append(this.d1List, t2-t1)
    this.d2List = append(this.d2List, t4-t3)
    this.d3List = append(this.d3List, t4-t1)

    //fmt.Println("t:", (t2-t1 - t4+t3) / 2 )
    if len(this.d1List) >= this.needNumber {
        //已经达到采样上限

        this._timeError()
        return false
    }

    return true
} //push

func (this *LGNtp1) _timeError() {
    sumd1 := 0
    sumd2 := 0
    n1 := len(this.d1List)
    n2 := len(this.d2List)

    for _,_d1 := range this.d1List {
        sumd1 += _d1
    }

    for _,_d2 := range this.d2List {
        sumd2 += _d2
    }

    //fmt.Println("sumd1:%d, sumd2:%d", sumd1, sumd2, n1,n2)
    this.timeError = (int)((n1 * sumd2 - n2 * sumd1) / (2 * n1 * n2))
}

func (this *LGNtp1) TimeError() int {
    return this.timeError
} //timeError

func (this *LGNtp1) Delay(t1,t2,t3,t4 int) (downDelay, upDelay, delay int) {
    downDelay = t2 - t1 - this.timeError
    upDelay = t4 - t3 + this.timeError

    sumd3 := 0
    for _,_d3 := range this.d3List {
        sumd3 += _d3
    }
    delay = int(sumd3 / len(this.d3List))

    return
}


/////////////////////////////////////////////////
type LGNtp struct {
    tList []float64
    d3List []int

    needNumber int
    maxDelay int

    timeError int
}

func LGNewNtp(needNumber, maxDelay int) *LGNtp {
    return &LGNtp{
        tList: []float64{},
        d3List: []int{},
        needNumber: needNumber,
        maxDelay: maxDelay,
    }
}

func (this *LGNtp) Init() {
    this.tList = []float64{}
    this.d3List = []int{}
} //init

func (this *LGNtp) Push(t1,t2,t3,t4 int) (ok bool) {
    if t2-t1 + t4-t3 > this.maxDelay {
        return true
    }
    t := ((float64)(t2-t1 - t4 +t3))/2
    this.tList = append(this.tList, t)
    this.d3List = append(this.d3List, t4-t1)

    //fmt.Println("t:",t )
    if len(this.tList) >= this.needNumber {
        //已经达到采样上限

        this._timeError()
        return false
    }

    return true
}

func (this *LGNtp) _timeError() {
    sumd1 := 0.0
    n1 := (float64)(len(this.tList))

    for _,_t:= range this.tList {
        sumd1 += _t
    }

    //fmt.Printf("sumd1:%f, n1:%f \n", sumd1, n1)
    this.timeError = (int)(sumd1 / n1)
}

func (this *LGNtp) TimeError() int {
    return this.timeError
} //timeError

func (this *LGNtp) Delay(t1,t2,t3,t4 int) (downDelay, upDelay,delay int) {
    downDelay = t2 - t1 - this.timeError
    upDelay = t4 - t3 + this.timeError

    sumd3 := 0
    for _,_d3 := range this.d3List {
        sumd3 += _d3
    }
    delay = int(sumd3 / len(this.d3List))
    return
}

