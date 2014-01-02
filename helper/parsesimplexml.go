/*=============================================================================
#     FileName: parsesimplexml.go
#       Author: sunminghong, allen.fantasy@gmail.com, http://weibo.com/5d13
#         Team: http://1201.us
#   LastChange: 2014-01-01 20:06:00
#      History:
=============================================================================*/

package helper

func ParseSimpleXml(xmlBytes []byte) map[string][]byte {
    //t := time.Unix(0, m.Timestamp)
    //datetime := strftime(*logtimeFormat, t)
   var lastToken,lastTokenClose string
    var lastData []byte

    var valBuf [2048]byte
    tokenBuf := [50]byte{}
    flag := 0
    j := 0

    data := make(map[string][]byte)
    for _,b := range xmlBytes {
        switch flag {
        case 0: //read tokenBegin start
            if b == '<' {
                flag = 1
                j = 0
                continue
            }
        case 1:
            if b == '>' {
                //read tokenBegin end
                lastToken = string(tokenBuf[:j])
                flag = 2
                j = 0
            } else {
                //read tokenBegin
                tokenBuf[j] = b
                j ++
            }
            continue
        case 2:
            if b == '<' {
                //read tokenClose start
                lastData = append(lastData,valBuf[:j]...)

                //reader </xxx>
                flag = 3
                j = 0
            } else {
                // read value
                valBuf[j] = b
                j ++
            }
            continue
        case 3:
            if b == '/' {
                //read tokenClose start 2
                flag = 4
                continue
            } else {
                //if it is not tokenEnd ,then restore data to value
                valBuf[j] = '<'
                j++
                valBuf[j] = b
                j ++
                lastData = append(lastData,valBuf[:j]...)
                j = 0
                flag = 2
            }
        case 4:
            if b == '>' {
                //read tokenClose end
                lastTokenClose = string(tokenBuf[:j])
                if lastTokenClose == lastToken {
                    val := make([]byte,len(lastData))
                    copy(val[0:],lastData)
                    data[lastToken] = val
                    lastData = []byte{}
                    flag = 0
                    j = 0
                    continue
                }

                //if beginToken != endToken ,then restore tokenClose to valBuf
                oldj := j
                j = 0
                valBuf[j] = '<'
                j ++
                valBuf[j] = '/'
                j ++

                for i:=0;i< oldj;i ++ {
                    valBuf[j] = tokenBuf[i]
                    j ++
                }
                valBuf[j] = '>'
                j++
                lastData = append(lastData,valBuf[:j]...)
                j = 0

                flag = 2
            } else {
                tokenBuf[j] = b
                j ++
            }
        }

    }

    return data
}
