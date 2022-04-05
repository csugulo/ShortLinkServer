package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/csugulo/ShortLinkServer/config"
	"github.com/csugulo/ShortLinkServer/controllers"
	"github.com/csugulo/ShortLinkServer/db"
	"github.com/gin-gonic/gin"
)

func parseArgs() {
	parser := argparse.NewParser("ShortLinkServer", "short link web service")

	domain := parser.String("d", "domain", &argparse.Options{Required: false, Default: "localhost", Help: "domain name"})
	port := parser.Int("p", "port", &argparse.Options{Required: false, Default: 8080, Help: "http port"})
	sqliteDBPath := parser.String("s", "sqlite", &argparse.Options{Required: false, Default: "sqlite.db", Help: "sqlite db path"})
	rocksDBPath := parser.String("r", "rocksdb", &argparse.Options{Required: false, Default: "rocksdb", Help: "rocksdb path"})
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
	config.Conf = config.Config{
		Domian:       *domain,
		Port:         *port,
		SqliteDBPath: *sqliteDBPath,
		RocksDBPath:  *rocksDBPath,
	}
}

type App struct {
	httpServer *gin.Engine
}

func main() {
	parseArgs()
	db.InitRocksDB(config.Conf.RocksDBPath)
	db.InitSqliteDB(config.Conf.SqliteDBPath)

	controllers.InitServer()
	controllers.Server.Run(fmt.Sprintf("0.0.0.0:%v", config.Conf.Port))
}
