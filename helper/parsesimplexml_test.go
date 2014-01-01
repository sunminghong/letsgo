// This is a client that writes out to a file, and optionally rolls the file

package helper

import (
	"fmt"
    "testing"
)

func Test_ParseSimpleXml(t *testing.T) {
    text := []byte(`<dir1>YHQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type><body>[INFO] 2014-01-01 14:50:34 based  tornado  framework!</body>`)

    xmlData := ParseSimpleXml(text)

    for xmlData,val := range xmlData {
        fmt.Println(xmlData,string(val))
    }

    if string(xmlData["dir1"]) != "YHQ" {
        t.Errorf("xmlData dir1's value not eq 'YHQ'")
    }

    if string(xmlData["dir2"]) != "1201.us" {
        t.Errorf("xmlData dir2's value not eq '1201.us'")
    }

    if string(xmlData["dir3"]) != "local18" {
        t.Errorf("xmlData dir3's value not eq 'local18'")
    }

    if string(xmlData["type"]) != "fightserver" {
        t.Errorf("xmlData type's value not eq 'fightserver'",xmlData["type"])
    }

    s := "[INFO] 2014-01-01 14:50:34 based  tornado  framework!"
    if string(xmlData["body"]) != s {
        t.Errorf("xmlData body's value not eq '...'",xmlData["type"])
    }


    text = []byte(`<dir1>YhQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type><body>[INFO] <asdfas>2014-01-01 14</asdfas>:50:34 _pqd: 0.6</body>`)

    s = "[INFO] <asdfas>2014-01-01 14</asdfas>:50:34 _pqd: 0.6"

    xmlData = ParseSimpleXml(text)

    for xmlData,val := range xmlData {
        fmt.Println(xmlData,string(val))
    }

    if string(xmlData["dir1"]) != "YhQ" {
        t.Error("xmlData dir1's value not eq 'YHQ'")
    }

    if string(xmlData["dir2"]) != "1201.us" {
        t.Error("xmlData dir2's value not eq '1201.us'")
    }

    if string(xmlData["dir3"]) != "local18" {
        t.Error("xmlData dir3's value not eq 'local18'")
    }

    if string(xmlData["type"]) != "fightserver" {
        t.Error("xmlData type's value not eq 'fightserver'",string(xmlData["type"]))
    }

    if string(xmlData["body"]) != s {
        t.Error("xmlData body's value not eq '...'","\nbody "+string(xmlData["body"]))
    }


    //fmt.Println(xmlData)
}
