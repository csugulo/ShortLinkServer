package dals

import (
	"github.com/tecbot/gorocksdb"
)

var readOptions *gorocksdb.ReadOptions
var writeOptions *gorocksdb.WriteOptions

func SetUrlMap(db *gorocksdb.DB, urlID, url string) (err error) {
	return db.Put(writeOptions, []byte(urlID), []byte(url))
}

func GetUrl(db *gorocksdb.DB, urlID string) (url string, err error) {
	var slice *gorocksdb.Slice
	if slice, err = db.Get(readOptions, []byte(urlID)); err != nil {
		return
	}
	return string(slice.Data()), nil
}

func init() {
	readOptions = gorocksdb.NewDefaultReadOptions()
	writeOptions = gorocksdb.NewDefaultWriteOptions()
}
