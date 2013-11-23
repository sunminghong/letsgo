/*=============================================================================
#     FileName: log.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-11-19 18:36:07
#      History:
=============================================================================*/

package log

import (
    "log"
    "os"
    "strings"
)

//--------------------
// LOG LEVEL
//--------------------

// Log levels to control the logging output.
const (
    LGLevelLGTrace = iota
    LGLevelDebug
    LGLevelInfo
    LGLevelWarning
    LGLevelError
    LGLevelCritical
)

// logLGLevel controls the global log level used by the logger.
var level = LGLevelLGTrace

// LogLGLevel returns the global log level and can be used in
// own implementations of the logger interface.
func LGGetLevel() int {
    return level
}

// SetLogLGLevel sets the global log level used by the simple
// logger.
func LGSetLevel(l int) {
    level = l
}
// logger references the used application logger.
var LetsLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

// SetLogger sets a new logger.
func LGSetLogger(l *log.Logger) {
    LetsLogger = l
}

func sout(t string,v ...interface{}) {
    s,ok := v[0].(string)
    if !ok || strings.Index(s,"%") == -1 {
        LetsLogger.Printf("["+ t +"] %v\n", v)
    } else {
        LetsLogger.Printf("["+ t +"] " + s + "\n", v[1:]...)
    }
}

// LGTrace logs a message at trace level.
func LGTrace(v ...interface{}) {
    if level <= LGLevelLGTrace {
        sout("T",v...)
    }
}

// Debug logs a message at debug level.
func LGDebug(v ...interface{}) {
    if level <= LGLevelDebug {
        sout("D",v...)
    }
}

// Info logs a message at info level.
func LGInfo(v ...interface{}) {
    if level <= LGLevelInfo {
        sout("I",v...)
    }
}

// Warning logs a message at warning level.
func LGWarn(v ...interface{}) {
    if level <= LGLevelWarning {
        sout("W",v...)
    }
}

// Error logs a message at error level.
func LGError(v ...interface{}) {
    if level <= LGLevelError {
        sout("E",v...)
    }
}

// Critical logs a message at critical level.
func LGCritical(v ...interface{}) {
    if level <= LGLevelCritical {
        sout("C",v...)
    }
}

/*
package main

import (
       "runtime"
       "fmt"
)

func main() {
        funcName, file, line, ok := runtime.Caller(0)
        if ok {
            fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
            fmt.Printf("file: %s    line=%d\n", file, line)
        }
}
*/
