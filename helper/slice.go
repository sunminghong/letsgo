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

type LGSliceInt []int

func (self *LGSliceInt) RemoveAtIndex(index int) (ok bool) {
    ok = true

    l := len(*self)
    if index > l-1 {
        return false
    }

    if index == l-1 {
        if index == 0 {
            *self = []int{}
            return
        }

        *self = (*self)[:index]
        return
    }

    if index == 0 {
        *self = (*self)[index + 1:]
        return
    }

    *self = append((*self)[:index] , (*self)[index+1:]...)
    return
}

func (self LGSliceInt) Eq(slice []int) bool {
    if len(self) != len(slice) {
        return false
    }
    for i,v := range slice {
        if self[i] != v {
            return false
        }
    }

    return true

}

