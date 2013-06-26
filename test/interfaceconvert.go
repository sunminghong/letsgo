/*=============================================================================
#     FileName: gate.go
#         Desc: game gate server
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-09 10:09:28
#      History:
=============================================================================*/
package main

import (
    "fmt"
)

type iA interface {
    print()
}

type A struct {
    i int
    j int
}

type B struct {
    *A
    ii int
    jj int
}

func (a *A) print() {
    fmt.Println("print",a.i)
}

func (a *B) print() {
    fmt.Println("print",a.ii)
}

func (b *B) show() {
    fmt.Println("show")
}

func m1(a *A) {
    a.print()
}

func m2(o iA) {
    bb,ok :=o.(*B)

    if ok {
        bb.show()
        fmt.Println(bb.jj)
    }
}


func main() {
    a := &A{1,2}
    b := &B{A:&A{3,4},ii:5,jj:6}

    m1(a)
    m2(b)
    
}
