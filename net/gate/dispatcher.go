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
    . "github.com/sunminghong/letsgo/helper"
    . "github.com/sunminghong/letsgo/log"
)

type LGDefaultDispatcher struct {
    messageCodemaps map[int]LGSliceInt
    removeGrids map[int]int
}

func LGNewDispatcher() *LGDefaultDispatcher {
    r := &LGDefaultDispatcher{make(map[int]LGSliceInt),make(map[int]int)}
    return r
}

func (r *LGDefaultDispatcher)Init() {
    r.messageCodemaps = make(map[int]LGSliceInt)
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
    LGTrace("messagecodemaps1:",r.messageCodemaps)
}

func (r *LGDefaultDispatcher) Remove(gridID int) {
    LGTrace("removegrids:",r.removeGrids)
    r.removeGrids[gridID] = 1
}

func (r *LGDefaultDispatcher) addDisp(gridID int, code int) {
    dises,ok := r.messageCodemaps[code]
    if ok {
        r.messageCodemaps[code] = append(dises,gridID)
        return
    }

    r.messageCodemaps[code] = LGSliceInt{gridID}
}

func (r *LGDefaultDispatcher) Dispatch(messageCode int) (gridID int,ok bool) {
    gcode := r.groupCode(messageCode)

    gridIDArr,ok := r.messageCodemaps[gcode]
    LGTrace("gridIDArr,gcode:",gridIDArr,gcode)
    if !ok {
        gridIDArr,ok = r.messageCodemaps[0]
        LGTrace("gridIDArr2,gcode:",gridIDArr,gcode)
    }

    if !ok {
        return 0,false
    }

    l := len(gridIDArr)
    i := rand.Intn(l)
    LGTrace("rand:",l,i)
    for i<l {
        gridID = gridIDArr[i]
        _,ok := r.removeGrids[gridID]
        LGTrace("removeGrids: ",r.removeGrids)
        if ok {
            LGTrace("this grid is already down",gridID)

            gridIDArr.RemoveAtIndex(i)
            r.messageCodemaps[gcode] = gridIDArr

            LGTrace("removed ",i,":",gridIDArr)
            l = len(gridIDArr)
            if l == 0 {
                break
            }

            if i >= l {
                i = 0
            }
            continue
        }

        LGTrace(
            "dispatcher Handler func messageCode,messageCode,gridID:",
        messageCode,gcode,gridID,gridIDArr)

        return gridID,true
    }

    return 0,false
}

//将协议编号分组以供Dispatch决策用那个Grid 来处理
func (r *LGDefaultDispatcher) groupCode(messageCode int) int {
    return int(messageCode / 100)
}

