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

import (
    "testing"
    "fmt"
)

func Test_ntp(t *testing.T) {
    fmt.Println("\nTest_ntp:")
    ntp := LGNewNtp(5, 100)

    ntp.Push(1,5,6,5)
    ntp.Push(6,11,12,12)
    ntp.Push(14,20,21,21)
    ntp.Push(23,26,27,26)
    if ntp.Push(27,31,31,32) == false {
        fmt.Println("t:" , ntp.TimeError())
        down,up := ntp.Delay(27,31,31,32)
        fmt.Println("up,down:" ,down,up)
    }

    //ntp.Push()
}
