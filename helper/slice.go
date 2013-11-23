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

func (self *LGSliceInt) RemoveValue(val int) (ok bool) {
    n := LGSliceInt{}
    for _,v := range *self {
        if v != val {
            n=append(n,v)
        }
    }
    *self = n
    return true
}

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

type LGISearch interface {
    Len() int
    //if eq return 0
    //if gt return 1
    //if lt return -1
    Compare(index,val int) int
}

//General binary search

func LGSliceSearch(objs LGISearch,key int) int {
    high := objs.Len() -1
    low := 0
    for low <= high {
        mid := (low + high) >> 1
        r := objs.Compare(mid,key)
        switch r {
            case 0:
                return mid
            case -1:
                low = mid + 1
            case 1:
                high = mid - 1
            }
    }

    return -1
}

//func LGSliceDelete(

