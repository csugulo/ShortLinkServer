package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/csugulo/ShortLinkWeb/service"
)

func parseArgs() (configPath string) {
	parser := argparse.NewParser("ShortLinkWeb", "short link web service")
	configPathPtr := parser.String("c", "config", &argparse.Options{Required: true, Help: "config path"})
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}
	return *configPathPtr
}

func main() {
	configPath := parseArgs()
	service := service.NewService(configPath)
	service.Init()
	service.Start()
}
