/*=============================================================================
#     FileName: defaultdispatcher.go
#         Desc: default dispatcher
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-06 14:37:57
#      History:
=============================================================================*/
package gate

import (
    "strconv"
    "strings"
    "math/rand"
    . "github.com/sunminghong/letsgo/log"
)

type LGDefaultDispatcher struct {
    messageCodemaps map[int][]int
}

func LGNewDispatcher() *LGDefaultDispatcher {
    r := &LGDefaultDispatcher{make(map[int][]int)}
    return r
}

func (r *LGDefaultDispatcher)Init()  {
    r.messageCodemaps = make(map[int][]int)
}

func (r *LGDefaultDispatcher) Add(gridID int, messageCodes *string) {
    cs := strings.Replace(*messageCodes," ","",-1)
    if len(cs) ==0 {
        r.addDisp(gridID,0)
    }

    codes:= strings.Split(cs,",")
    LGTrace("add disp",codes)
    for _,p_ := range codes {
        p := strings.Trim(p_," ")
        if len(p) == 0 {
            continue
        }
        pmessageCode, err := strconv.Atoi(p)
        if err ==nil {
            r.addDisp(gridID,pmessageCode)
        }
    }
    LGTrace(r.messageCodemaps)
}

func (r *LGDefaultDispatcher) addDisp(gridID int, code int) {
    dises,ok := r.messageCodemaps[code]
    if ok {
        r.messageCodemaps[code] = append(dises,code)
        return
    }

    r.messageCodemaps[code] = []int{gridID}
}


func (r *LGDefaultDispatcher) Dispatch(messageCode int) (gridID int,ok bool) {
    gcode := r.GroupCode(messageCode)

    var gridIDArr []int
    gridIDArr,ok = r.messageCodemaps[gcode]
    if !ok {
        gridIDArr,ok = r.messageCodemaps[0]
    }

    if ok {
        i := rand.Intn(len(gridIDArr))
        gridID = gridIDArr[i]
        LGTrace(
            "dispatcher Handler func messageCode,messageCode,gridID:",
        messageCode,gcode,gridID,gridIDArr)
        return gridID,ok
    }
    return 0,false
}

//将协议编号分组以供Dispatch决策用那个Grid 来处理
func (r *LGDefaultDispatcher) GroupCode(messageCode int) int {
    return int(messageCode / 100)
}

