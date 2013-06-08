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

type ClientMap struct {
    maplock *sync.RWMutex

    maps map[int]IClient
    mapsByName map[string]int
}

func NewClientMap() *ClientMap {
    return &ClientMap{
        maplock: new(sync.RWMutex),
        maps: make(map[int]IClient),
        mapsByName: make(map[string]int),
    }
}

func (tm *ClientMap) Add(cid int,name string, client IClient) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    tm.maps[cid] = client
    if len(name) > 0 {
        tm.mapsByName[name] = cid
    }
}

func (tm *ClientMap) Remove(cid int) {
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

func (tm *ClientMap) RemoveByName(name string) {
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

func (tm *ClientMap) Get(cid int) IClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    c, ok := tm.maps[cid]
    if ok {
        return c
    }
    return nil
}

func (tm *ClientMap) GetByName(name string) IClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    cid, ok := tm.mapsByName[name]
    if ok {
        return tm.maps[cid]
    }
    return nil
}

func (tm *ClientMap) All() map[int]IClient {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return tm.maps
}

func (tm *ClientMap) Len() int {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    return len(tm.maps)
}

