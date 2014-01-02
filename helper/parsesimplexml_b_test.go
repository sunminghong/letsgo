// This is a client that writes out to a file, and optionally rolls the file

package helper

import (
    "testing"
)

func Benchmark_ParseSimpleXml(t *testing.B) {
    s := `[INFO] <asdfasfasfasfd></aasdflaslfdsdf>sdfjasdlfjhasdfkhjasdlfkhasdlfkhaslidfhassdkfa
    hskfhasdkfhasdfh看我忘了我看完了看完我看完我看完我看完我<asdfa
    hskfhasdkfhasdfh看我忘了我看完了看完我看完我看完我看完我<asdfa
    hskfhasdkfhasdfh看我忘了我看完了看完我看完我看完我看完我<asdfa
    hskfhasdkfhasdfh看我忘了我看完了看完我看完我看完我看完我<asdfa
    s>2014-01-01 14</asdfas>:50:34 _pqd: 0.6[INFO] <asdfas>2014-01-01 14</asdfas>:50:34 _pqd: 0.6`
    text := []byte(`<body>`+s+`</body><dir1>YhQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type>`)


    for i := 0; i < t.N; i++ {
        ParseSimpleXml(text)
    }
}

func Benchmark_ParseSimpleXml2(t *testing.B) {
    s := "[INFO] <asdfasfasfasfd></aasdflasjd>flasjkdflaskdfjalskdfjasldfkjaslfdsdf>sdfjasdlfjhasdfkhjasdlfkhasdlfkhaslidfhassdkfahskfhasdkfhasdfh看我忘了我看完了看完我看完我看完我看完我<asdfas>2014-01-01 14</asdfas>:50:34 _pqd: 0.6[INFO] <asdfas>2014-01-01 14</asdfas>:50:34 _pqd: 0.6"
    text := []byte(`<body>`+s+`</body><dir1>YhQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type>`)


    for i := 0; i < t.N; i++ {
        ParseSimpleXml(text)
    }
}

func Benchmark_ParseSimpleXml3(t *testing.B) {
    s := "[INFO] <asdfas>2014-01-01 14</asdfas>:50:34 _pqd: 0.6"
    text := []byte(`<body>`+s+`</body><dir1>YhQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type>`)


    for i := 0; i < t.N; i++ {
        ParseSimpleXml(text)
    }
}

func Benchmark_ParseSimpleXml32(t *testing.B) {
    s := "[INFO] asdfas2014-01-01 14asdfasaldfjasldfjasldf可的杀伤力的房间爱睡懒觉阿斯顿浪费空间按时劳动法看静安寺老地方空间阿斯兰；飞撒了点附近啥的发撒了点付款就爱上打翻了按时缴费啥的了饭卡水电费了奥斯卡京东方啊水力发电50:34 _pqd: 0.6"
    text := []byte(`<body>`+s+`</body><dir1>YhQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type>`)


    for i := 0; i < t.N; i++ {
        ParseSimpleXml(text)
    }
}

func Benchmark_ParseSimpleXml4(t *testing.B) {
    s := "[INFO] 2014-01-01 14:sdlfsaflasdflasjflasfdas50:34 _pqd: 0.6"
    text := []byte(`<body>`+s+`</body><aaaaaa>asdf</aaaaaa><bbbbb>sadfasdfasfdsaf</bbbbb><dir1>YhQ</dir1><dir2>1201.us</dir2><dir3>local18</dir3><type>fightserver</type>`)


    for i := 0; i < t.N; i++ {
        ParseSimpleXml(text)
    }
}
