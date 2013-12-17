/*=============================================================================
#     FileName: log.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2013-11-19 18:36:07
#      History:
=============================================================================*/

package log

import "github.com/sunminghong/freelog"

//--------------------
// LOG LEVEL
//--------------------

/*
// Log levels to control the logging output.
const (
	LGLevelAll = iota
	LGLevelTrace
	LGLevelDebug
	LGLevelInfo
	LGLevelWarn
	LGLevelError
    LGLevelPanic
	LGLevelFatal
	LGLevelOff
)*/

func LGSetLogger(inifile *string) {
    freelog.CallDepth = 3
    freelog.Start(inifile)
}

func LGTrace(v ...interface{}) {
    freelog.Trace(v...)
}

func LGDebug(v ...interface{}) {
    freelog.Debug(v...)
}

func LGInfo(v ...interface{}) {
    freelog.Info(v...)
}

func LGWarn(v ...interface{}) {
    freelog.Warn(v...)
}

// Error logs a message at error level.
func LGError(v ...interface{}) {
    freelog.Error(v...)
}

// Critical logs a message at critical level.
func LGPanic(v ...interface{}) {
    freelog.Panic(v...)
}

func LGPrintPanicStack() {
    freelog.PrintPanicStack()
}
