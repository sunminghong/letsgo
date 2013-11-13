/*=============================================================================
#     FileName: server.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-03 14:27:57
#      History:
=============================================================================*/
package net

import (
    "sync"
)

type LGConnectionMap struct {
    maplock *sync.RWMutex

    maps map[int]LGIConnection
    mapsByName map[string]int
}

func LGNewConnectionMap() *LGConnectionMap {
    return &LGConnectionMap{
        maplock: new(sync.RWMutex),
        maps: make(map[int]LGIConnection),
        mapsByName: make(map[string]int),
    }
}

func (tm *LGConnectionMap) Add(cid int,name string, client LGIConnection) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    tm.maps[cid] = client
    if len(name) > 0 {
        tm.mapsByName[name] = cid
    }
}

func (tm *LGConnectionMap) Remove(cid int) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    _, ok := tm.maps[cid]
    if ok {
        name := tm.maps[cid].GetName()
        if len(name)>0 {
            _,ok :=tm.mapsByName[name]
            if ok {
                delete(tm.mapsByName,name)
            }
        }
        delete(tm.maps, cid)
    }
}

func (tm *LGConnectionMap) RemoveByName(name string) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    cid,ok :=tm.mapsByName[name]
    if ok {
        _, ok := tm.maps[cid]
        if ok {
            delete(tm.maps, cid)
        }
        delete(tm.mapsByName,name)
    }
}

func (tm *LGConnectionMap) Get(cid int) LGIConnection {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    c, ok := tm.maps[cid]
    if ok {
        return c
    }
    return nil
}

func (tm *LGConnectionMap) GetByName(name string) LGIConnection {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    cid, ok := tm.mapsByName[name]
    if ok {
        return tm.maps[cid]
    }
    return nil
}

func (tm *LGConnectionMap) All() map[int]LGIConnection {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return tm.maps
}

func (tm *LGConnectionMap) Len() int {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return len(tm.maps)
}

