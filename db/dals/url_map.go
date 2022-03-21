package dals

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tecbot/gorocksdb"
)

var readOptions *gorocksdb.ReadOptions
var writeOptions *gorocksdb.WriteOptions

func SetUrlMap(db *gorocksdb.DB, urlID, url string, expireTime time.Time) (err error) {
	return db.Put(writeOptions, []byte(urlID), []byte(fmt.Sprintf("%v %v", url, expireTime.Unix())))
}

func GetUrl(db *gorocksdb.DB, linkID string) (url string, expireTime time.Time, err error) {
	var slice *gorocksdb.Slice
	if slice, err = db.Get(readOptions, []byte(linkID)); err != nil {
		return
	}
	if slice.Size() == 0 {
		return
	}
	tuple := strings.Split(string(slice.Data()), " ")
	if len(tuple) != 2 {
		err = fmt.Errorf("unkown value format, value: %v", string(slice.Data()))
		return
	}
	var unixTime int64
	if unixTime, err = strconv.ParseInt(tuple[1], 10, 64); err != nil {
		return
	}
	return tuple[0], time.Unix(unixTime, 0), nil
}

func init() {
	readOptions = gorocksdb.NewDefaultReadOptions()
	writeOptions = gorocksdb.NewDefaultWriteOptions()
}
