package main

import (
//        "reflect"
        "fmt"
       )

type tfunc func(a int) int

func a(a int) int {
    return a
}

func b(a int) int {
    return a+a
}

func main(){
funcs := make(map[string]tfunc)
           funcs["a"] = a
           funcs["b"] = b

           fmt.Println(funcs["a"])

           fmt.Println(funcs["a"](3))
           fmt.Println(funcs["b"](7))
}
