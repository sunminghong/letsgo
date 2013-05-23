/*=============================================================================
#     FileName: route.go
#         Desc: route for received request from client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-22 11:40:30
#      History:
=============================================================================*/
package net

import (
    //"encoding/binary"
)
/*
const (
    mask1 = byte(0x59)
    mask2 = byte(0x7a)

    DATAPACKET_TYPE_GENERAL = 0
    DATAPACKET_TYPE_DELAY = 1
    DATAPACKET_TYPE_BOARDCAST = 3
)*/

type RouteMap struct {
    maplock sync.RWMutex

    maps map[int]IClient
    mapsByName map[string]int
}

func (tm *RouteMap) Add(cid int,name string, client IClient) {
    tm.maplock.Lock()
    defer tm.maplock.Unlock()

    tm.maps[cid] = client
    if len(name) > 0 {
        tm.mapsByName[name] = cid
    }
}

func (tm *RouteMap) Remove(cid int) {
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

func (tm *RouteMap) RemoveByName(name string) {
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

func (tm *RouteMap) Get(cid int) IClient {
    c, ok := tm.maps[cid]
    if ok {
        return c
    }
    return nil
}

func (tm *RouteMap) GetByName(name string) IClient {
    cid, ok := tm.mapsByName[name]
    if ok {
        return tm.maps[cid]
    }
    return nil
}

func (tm *RouteMap) All() map[int]IClient {
    return tm.maps
}

func NewRouteMap() *RouteMap { return &RouteMap{maps: make(map[int]IClient),mapsByName: make(map[string]int)} }
