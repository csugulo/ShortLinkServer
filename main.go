package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/csugulo/ShortLinkWeb/config"
	"github.com/csugulo/ShortLinkWeb/controllers"
	"github.com/csugulo/ShortLinkWeb/db"
	"github.com/gin-gonic/gin"
)

type Args struct {
	configPath string
}

func parseArgs() Args {
	parser := argparse.NewParser("ShortLinkWeb", "short link web service")
	configPathPtr := parser.String("c", "config", &argparse.Options{Required: true, Help: "config path"})
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
	return Args{
		configPath: *configPathPtr,
	}
}

type App struct {
	httpServer *gin.Engine
}

func main() {
	args := parseArgs()

	config.InitConf(args.configPath)
	db.InitRocksDB(config.Conf.GetString("rocksdb.path"))
	db.InitSqliteDB(config.Conf.GetString("sqlite.path"))

	controllers.InitServer()
	controllers.Server.Run(fmt.Sprintf("0.0.0.0:%v", config.Conf.GetString("http.port")))
}
