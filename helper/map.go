/*=============================================================================
#     FileName: idassign.go
#         Desc: class with unix's process id alloc
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-05 10:46:10
#      History:
=============================================================================*/
package helper

import (
    "sync"
)

type LGMap struct {
    lock *sync.RWMutex
    bm   map[int]interface{}
}

func NewLGMap() *LGMap {
    return &LGMap{
        lock: new(sync.RWMutex),
        bm: make(map[int]interface{}),
    }
}

//Get from maps return the k's value
func (m *LGMap) Get(k int) (interface{},bool) {
    m.lock.RLock()
    defer m.lock.RUnlock()

    val, ok := m.bm[k]
    return val,ok
}

// if the key is already in the map and changes nothing.
func (m *LGMap) Set(k int, v interface{}) {
    m.lock.Lock()
    defer m.lock.Unlock()
    if val, ok := m.bm[k]; !ok {
        m.bm[k] = v
    } else if val != v {
        m.bm[k] = v
    }
}

// Returns true if k is exist in the map.
func (m *LGMap) Check(k int) bool {
    m.lock.RLock()
    defer m.lock.RUnlock()
    if _, ok := m.bm[k]; !ok {
        return false
    }
    return true
}

func (m *LGMap) Delete(k int) {
    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.bm, k)
}

func (m *LGMap) Clear() {
    m.lock.Lock()
    defer m.lock.Unlock()
    m.bm = make(map[int]interface{})
}
