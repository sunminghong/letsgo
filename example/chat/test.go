package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" //引入mysql驱动
	"time"
)

// 处理/upload 逻辑

func test() {
    userid  := 1
    tag     := "tag"
    describ := "describ"
    fname:= "handler.Filename"
    fmt.Println(fname)
    var status int32 = 0
    var ip string    = "xxx.xxx.xx.184"
    ptime := time.Now().Unix()
    db, err := sql.Open("mysql", "root:root@tcp(192.168.18.18:3306)/test") //第一个参数数驱动名
    checkErr(err)
    stmt, err := db.Prepare("INSERT INTO queue set filename=?,userid=?,status=?,ptime=?,ip=?,tag=?,describ=?")
    checkErr(err)
    res, err := stmt.Exec(fname, userid, status, ptime, ip,tag,describ)
    checkErr(err)
    id, err := res.LastInsertId()
    checkErr(err)
    fmt.Println(id)
    db.Close()
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
    test()
}
