package config

type Config struct {
	Domian       string
	Port         int
	SqliteDBPath string
	RocksDBPath  string
}

var Conf Config
