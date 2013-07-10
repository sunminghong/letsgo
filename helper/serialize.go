/*=============================================================================
#     FileName: memcache.go
#         Desc: client of memcached client
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-07-08 11:51:07
#      History:
=============================================================================*/
package helper


import (
    "bytes"
    "encoding/gob"
	"encoding/json"
)

type LGISerialize interface {
    Serialize(src interface{}) (dst []byte, err error)
    Deserialize(src []byte, dst interface{}) (err error)
}

type LGGobSerialize struct {
}

// serialize encodes a value using gob.
func (self LGGobSerialize) Serialize(src interface{}) (v []byte, err error) {
    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err = enc.Encode(src)
    if err != nil {
        return
    }
    v = buf.Bytes()
    return
}

// deserialize decodes a value using gob.
func (self LGGobSerialize) Deserialize(src []byte, dst interface{}) (err error) {
    dec := gob.NewDecoder(bytes.NewBuffer(src))
    err = dec.Decode(dst)
    return
}

type LGJsonSerialize struct {
}

// serialize encodes a value using gob.
func (self LGJsonSerialize) Serialize(src interface{}) (v []byte, err error) {
    v,err = json.Marshal(src)
    if err != nil {
        return
    }
    return
}

// deserialize decodes a value using gob.
func (self LGJsonSerialize) Deserialize(src []byte, dst interface{}) (err error) {
    err = json.Unmarshal(src,&dst)
    return
}

