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

type LGClientMap struct {
    maplock *sync.RWMutex

    maps map[int]LGIClient
    mapsByName map[string]int
}

func LGNewClientMap() *LGClientMap {
    return &LGClientMap{
        maplock: new(sync.RWMutex),
        maps: make(map[int]LGIClient),
        mapsByName: make(map[string]int),
    }
}

func (tm *LGClientMap) Add(cid int,name string, client LGIClient) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    tm.maps[cid] = client
    if len(name) > 0 {
        tm.mapsByName[name] = cid
    }
}

func (tm *LGClientMap) Remove(cid int) {
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

func (tm *LGClientMap) RemoveByName(name string) {
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

func (tm *LGClientMap) Get(cid int) LGIClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    c, ok := tm.maps[cid]
    if ok {
        return c
    }
    return nil
}

func (tm *LGClientMap) GetByName(name string) LGIClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    cid, ok := tm.mapsByName[name]
    if ok {
        return tm.maps[cid]
    }
    return nil
}

func (tm *LGClientMap) All() map[int]LGIClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return tm.maps
}

func (tm *LGClientMap) Len() int {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return len(tm.maps)
}

