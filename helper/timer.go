/*=============================================================================
#     FileName: timer.go
#         Desc: class with unix's process id alloc
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-08-15 10:08:27
#      History:
=============================================================================*/
package helper

import (
    "time"
)

var BaseTimestamp int = 1375286400

func LGNetTimestamp(times ...interface{}) int {
    //减去 2013-8-1

    var unix int
    if len(times) == 0 {
        unix = int(time.Now().Unix())
        return unix - BaseTimestamp
    }

    ttime := times[0]
    if t, ok := ttime.(time.Time); ok {
        unix = int(t.Unix())
    } else if t, ok := ttime.(int); ok {
        unix = t
    } else if t, ok := ttime.(int64); ok {
        unix = int(t)
    }

    return unix - BaseTimestamp
}

func LGStrttime(strtime string, format ...string) (time.Time, error) {
    f := "2006-01-02 15:04:05"
    if len(format) > 0 {
        f = format[0]
    }
    return time.ParseInLocation(f, strtime, time.Local)
}

func LGStrftime(times ...interface{}) string {
    f := "2006-01-02 15:04:05"
    if len(times) == 0 {
        return time.Now().Local().Format(f)
    }

    ttime := times[0]

    if len(times) > 1 {
        f,_ = times[1].(string)
    }

    if t, ok := ttime.(time.Time); ok {
        return t.Format(f)
    } else if t, ok := ttime.(int); ok {
        return time.Unix(int64(t), 0).Format(f)
    } else if t, ok := ttime.(int64); ok {
        return time.Unix(t, 0).Format(f)
    }

    panic("LGStrftime's param is not time.Time or int、int64")
}

// return xxxx-xx-xx 00:00:00 的unix时间戳
func LGTodayUnix(times ...interface{}) int {
    var tt time.Time
    if len(times) == 0 {
        tt = time.Now()
    } else {
        ttime := times[0]
        if t, ok := ttime.(time.Time); ok {
            tt = t
        } else if t, ok := ttime.(int); ok {
            tt = time.Unix(int64(t), 0)
        } else if t, ok := ttime.(int64); ok {
            tt = time.Unix(t, 0)
        }
    }

    tt = tt.Local()
    tt = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.Local)

    return int(tt.Unix())
}

// return yyyymmdd 格式
func LGToday(times ...interface{}) int {
    var tt time.Time
    if len(times) == 0 {
        tt = time.Now()
    } else {
        ttime := times[0]
        if t, ok := ttime.(time.Time); ok {
            tt = t
        } else if t, ok := ttime.(int); ok {
            tt = time.Unix(int64(t), 0)
        } else if t, ok := ttime.(int64); ok {
            tt = time.Unix(t, 0)
        }
    }

    tt = tt.Local()
    return tt.Year()*10000 + int(tt.Month())*100 + tt.Day()
}

func LGYesterday(times ...interface{}) int {
    var tt time.Time
    if len(times) == 0 {
        tt = time.Now()
    } else {
        ttime := times[0]
        if t, ok := ttime.(time.Time); ok {
            tt = t
        } else if t, ok := ttime.(int); ok {
            tt = time.Unix(int64(t), 0)
        } else if t, ok := ttime.(int64); ok {
            tt = time.Unix(t, 0)
        }
    }

    tt = tt.AddDate(0, 0, -1)
    tt = tt.Local()
    mm := int(tt.Month())
    return tt.Year()*10000 + mm*100 + tt.Day()

}

// 计算两个时间的天数差值
func diffday(time1 int, time2 int) int {

    tt1 := time.Unix(int64(time1), 0)
    tt2 := time.Unix(int64(time2), 0)

    return tt1.YearDay() - tt2.YearDay()
}
