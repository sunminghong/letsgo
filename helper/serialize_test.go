/*=============================================================================
#     FileName: memcache_struct_test.go
#         Desc: memcache test struct set/get
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-10 10:41:38
#      History:
=============================================================================*/

package helper

import (
    //"os/exec"
    "testing"
    "fmt"
)

type A struct {
    I1 int
    Str1 string
    I2 int
    I3 int
    I4 int
    I5 int
    I6 int
    I7 int
    I8 int
    I9 int
    I10 int
    I11 int
    I12 int
    I13 int
    I14 int
    I15 int

    B1 *B
    Barr []*B
    Bmap2 map[string]B
    Bmap1 map[string]*B
}

type B struct {
    I2 int
    Str2 string
}

func (a1 *A) Eq(a2 *A) bool {
    if a1.I1 != a2.I1 {
        return false
    }
    if a1.Str1 != a2.Str1 {
        return false
    }
    for i,B1 := range a1.Barr {
        b2 := a2.Barr[i]
        if !B1.Eq(b2) {
            return false
        }
    }

    for i,B1 := range a1.Bmap1 {
        b2 := a2.Bmap1[i]
        if !B1.Eq(b2) {
            return false
        }
    }

    for i,B1 := range a1.Bmap2 {
        b2 := a2.Bmap2[i]
        if !B1.Eq(&b2) {
            return false
        }
    }
    return true
}

func (B1 *B) Eq(b2 *B) bool {
    if B1.I2 != b2.I2 {
        return false
    }
    if B1.Str2 != b2.Str2 {
        return false
    }
    return true
}

func NewA1A2() (a1,a2 *A) {
    //Bmap1 := make(map[int]*B)
    Bmap1 := map[string]*B {
        "11":&B{2,"def"},
        "22":&B{22,"defdef"},
        "33":&B{222,"defdefdef"},
    }

    Bmap2 := map[string]B {
        "aa":B{3,"ghi"},
        "bb":B{33,"ghighi"},
        "cc":B{333,"ghighighi"},
    }

    a1 = &A{
        I1:1,
        I2:12,
        I3:13,
        I4:1234,
        I5:1342,
        I6:123423,
        I7:123,
        I8:2341,
        I9:452234,
        I10:690,
        I11:2342,
        I12:56878,
        I13:213134,
        I14:2352,
        I15:234,

        Str1:"abc",
        Barr:[]*B{&B{2,"def"},&B{22,"defdef"},&B{222,"defdefdef"}},
        Bmap1:Bmap1,
        Bmap2:Bmap2,
    }

    a2 = &A{
        I1:11,
        Str1:"abcabc",
        Barr:[]*B{&B{22,"def"},&B{2222,"defdef"},&B{222222,"defdefdef"}},
        Bmap1:Bmap1,
        Bmap2:Bmap2,
    }
    return a1,a2
}

func NewSerialize() LGISerialize {
    return LGGobSerialize{}
}


func Test_gob_serialize(t *testing.T) {
    //fmt.Println("////////////////////////test serialize //////////////////////////////")
    gs := &LGGobSerialize{}

    v1 := 1
    v2 := -1

    v,err := gs.Serialize(v1)
    err = gs.Deserialize(v,&v2)
    //fmt.Println("v1,v,v2,err",v1,v,v2,err)

    if err != nil || v1 != v2 {
        t.Errorf("v1 != v2",v1,v2)
    }

    v1 = 0
    v2 = -2
    v,err = gs.Serialize(v1)
    err = gs.Deserialize(v,&v2)
    //fmt.Println("v1,v,v2,err",v1,v,v2,err)

    if err != nil || v1 != v2 {
        t.Errorf("v1 != v2",v1,v2)
    }

    a1,_ := NewA1A2()

    a2 := &A{}
    v,err = gs.Serialize(a1)
    fmt.Println("gob serialize length :",len(v))
    err = gs.Deserialize(v,&a2)
    //fmt.Println("a1,v,a2,err",a1,v,a2,err)

    if err != nil || !a1.Eq(a2) {
        t.Errorf("a1 != a2",a1,a2)
    }
}


func Test_json_serialize(t *testing.T) {
    //fmt.Println("////////////////////////test serialize //////////////////////////////")
    gs := &LGJsonSerialize{}

    v1 := 1
    v2 := -1

    v,err := gs.Serialize(v1)
    err = gs.Deserialize(v,&v2)
    //fmt.Println("v1,v,v2,err",v1,v,v2,err)

    if err != nil || v1 != v2 {
        t.Errorf("v1 != v2",v1,v2)
    }

    v1 = 0
    v2 = -2
    v,err = gs.Serialize(v1)
    err = gs.Deserialize(v,&v2)
    //fmt.Println("v1,v,v2,err",v1,v,v2,err)

    if err != nil || v1 != v2 {
        t.Errorf("v1 != v2",v1,v2)
    }

    a1,_ := NewA1A2()

    a2 := &A{}
    v,err = gs.Serialize(*a1)
    fmt.Println("json serialize length :",len(v))
    err = gs.Deserialize(v,&a2)
    //fmt.Println("a1,v,a2,err",a1,v,a2,err)

    if err != nil {
        t.Errorf("deserialize(A) is error",err)
    }

    if err != nil || !a1.Eq(a2) {
        t.Errorf("a1 != a2",a1,a2)
    }
}
