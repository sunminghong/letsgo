/*=============================================================================
#     FileName: log.go
#         Desc: logger
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-23 18:01:57
#      History:
=============================================================================*/
package log

import (
    "log"
    "os"
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
func LGLevel() int {
    return level
}

// SetLogLGLevel sets the global log level used by the simple
// logger.
func LGSetLGLevel(l int) {
    level = l
}

// logger references the used application logger.
var LetsLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// SetLogger sets a new logger.
func LGSetLogger(l *log.Logger) {
    LetsLogger = l
}

// LGTrace logs a message at trace level.
func LGTrace(v ...interface{}) {
    if level <= LGLevelLGTrace {
        LetsLogger.Printf("[T] %v\n", v)
    }
}

// Debug logs a message at debug level.
func LGDebug(v ...interface{}) {
    if level <= LGLevelDebug {
        LetsLogger.Printf("[D] %v\n", v)
    }
}

// Info logs a message at info level.
func LGInfo(v ...interface{}) {
    if level <= LGLevelInfo {
        LetsLogger.Printf("[I] %v\n", v)
    }
}

// Warning logs a message at warning level.
func LGWarn(v ...interface{}) {
    if level <= LGLevelWarning {
        LetsLogger.Printf("[W] %v\n", v)
    }
}

// Error logs a message at error level.
func LGError(v ...interface{}) {
    if level <= LGLevelError {
        LetsLogger.Printf("[E] %v\n", v)
    }
}

// Critical logs a message at critical level.
func LGCritical(v ...interface{}) {
    if level <= LGLevelCritical {
        LetsLogger.Printf("[C] %v\n", v)
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
