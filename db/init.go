package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/tecbot/gorocksdb"
)

var RocksDB *gorocksdb.DB
var SqliteDB *sql.DB

func InitRocksDB(dbPath string) {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	if rocksDB, err := gorocksdb.OpenDb(options, dbPath); err != nil {
		log.Fatalf("can not open rocksdb: %v\n, err: %v", dbPath, err)
	} else {
		RocksDB = rocksDB
	}
	log.Info("rocksdb inited")
}

func InitSqliteDB(dbPath string) {
	if db, err := sql.Open("sqlite3", dbPath); err != nil {
		log.Fatalf("can not open sqlite3: %v\n, err: %v", dbPath, err)
	} else {
		SqliteDB = db
	}
}
