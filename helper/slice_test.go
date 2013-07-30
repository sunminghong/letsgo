/*=============================================================================
#     FileName: slice.go
#         Desc: class with unix's process id alloc
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-25 17:30:03
#      History:
=============================================================================*/
package helper

import "testing"


func TestRemoveValue(t *testing.T) {
    a1 := LGSliceInt{0,1,2,3,4,5,6,7,8}
    a1 = a1[:]

    a1.RemoveValue(3)
    b1 := []int{0,1,2,4,5,6,7,8}
    if !a1.Eq(b1) {
        t.Error("is not eq:",b1,a1)
    }

    a1 = LGSliceInt{0,1,2,3,4,5,6,5,7,8}
    a1 = a1[:]

    a1.RemoveValue(5)
    b1 = []int{0,1,2,3,4,6,7,8}
    if !a1.Eq(b1) {
        t.Error("is not eq:",b1,a1)
    }

}


func TestRemoveAtIndex(t *testing.T) {
    a1 := LGSliceInt{0,1,2,3,4,5,6,7,8}
    a1 = a1[:]

    a1.RemoveAtIndex(0)
    b1 := []int{1,2,3,4,5,6,7,8}
    b1 = b1[:]
    if !a1.Eq(b1) {
        t.Error("is not eq:",b1,a1)
    }

    a1 = LGSliceInt{0,1,2,3,4,5,6,7,8}
    a1.RemoveAtIndex(3)
    b1 = []int{0,1,2,4,5,6,7,8}
    if !a1.Eq(b1) {
        t.Error("is not eq:",b1,a1)
    }

    a1 = LGSliceInt{0,1,2,3,4,5,6,7,8}
    a1.RemoveAtIndex(8)
    b1 = []int{0,1,2,3,4,5,6,7}
    if !a1.Eq(b1) {
        t.Error("is not eq:",b1,a1)
    }

    a1 = LGSliceInt{0}
    a1.RemoveAtIndex(0)
    b1 = []int{}
    if !a1.Eq(b1) {
        t.Error("is not eq:",b1,a1)
    }
}
