/*=============================================================================
#     FileName: interval.go
#         Desc: define a interval timer
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-08-30 09:30:16
#      History:
=============================================================================*/
package helper

import (
    "time"
    "errors"
)

type LGInterval struct {
    duration time.Duration
    stop chan bool
    callback func(self *LGInterval)

    IsRun bool
    timer *time.Ticker

    //callback func if is executing
    IsExecuting bool
}

func NewLGInterval(duration time.Duration,callbackfn func(self *LGInterval)) *LGInterval {
    return &LGInterval{duration,make(chan bool),callbackfn,false,nil,false}
}

func (self *LGInterval) Start(newDuration ...time.Duration) error {
    if len(newDuration)>0 {
        self.duration = newDuration[0]
    }

    if self.IsRun {
        return errors.New("the instance is already running!")
    }

    go func() {
        self.IsRun = true
        self.timer = time.NewTicker(self.duration)
        for {
            if !self.IsRun {
                goto quit
            }

            select {
            case <-self.timer.C:
                self.IsExecuting = true
                self.callback(self)
                self.IsExecuting = false
            }
        }

        quit:
        self.IsRun = false
        self.stop <- true
    }()

    return nil
}

func (self *LGInterval) Stop(callback ...func(interval *LGInterval)) {
    go func() {
        if self.IsRun {
            self.IsRun = false

            <-self.stop
            self.timer.Stop()

            if len(callback)>0 {
                callback[0](self)
            }
        }
    }()
}
