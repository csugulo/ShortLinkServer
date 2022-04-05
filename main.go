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

type Args struct {
	configPath string
	domain     string
}

func parseArgs() Args {
	parser := argparse.NewParser("ShortLinkServer", "short link web service")

	domainPtr := parser.String("d", "domain", &argparse.Options{Required: false, Default: "localhost", Help: "domain"})
	configPathPtr := parser.String("c", "config", &argparse.Options{Required: true, Help: "config path"})
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
	return Args{
		configPath: *configPathPtr,
		domain:     *domainPtr,
	}
}

type App struct {
	httpServer *gin.Engine
}

func main() {
	args := parseArgs()

	config.InitConf(args.configPath, args.domain)
	db.InitRocksDB(config.Conf.GetString("rocksdb.path"))
	db.InitSqliteDB(config.Conf.GetString("sqlite.path"))

	controllers.InitServer()
	controllers.Server.Run(fmt.Sprintf("0.0.0.0:%v", config.Conf.GetString("http.port")))
}
