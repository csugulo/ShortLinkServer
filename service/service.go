package service

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
	"time"

	"github.com/csugulo/ShortLinkWeb/dals"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tecbot/gorocksdb"
)

type Service struct {
	httpEngine *gin.Engine
	conf       *viper.Viper
	rocksDB    *gorocksdb.DB
}

func NewService(configPath string) *Service {
	conf := viper.New()
	conf.SetConfigFile(configPath)
	if err := conf.ReadInConfig(); err != nil {
		log.Fatalf("can not read config: %v\n, err: %v", configPath, err)
	}
	return &Service{
		conf: conf,
	}
}

func (service *Service) initLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceQuote:    true,
	})
}

func (service *Service) initRocksDB() {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	if rocksDB, err := gorocksdb.OpenDb(options, service.conf.GetString("rocksdb.path")); err != nil {
		log.Fatalf("can not open rocksdb: %v\n, err: %v", service.conf.GetString("rocksdb.path"), err)
	} else {
		service.rocksDB = rocksDB
	}
	log.Info("rocksdb inited")
}

func (service *Service) createUrlMap(url string) (urlID string, err error) {
	h := fnv.New32a()
	h.Write([]byte(url))
	digest := h.Sum32()
	for {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, digest)
		urlID = base64.StdEncoding.EncodeToString(bs)
		var existedUrl string
		if existedUrl, _, err = dals.GetUrl(service.rocksDB, urlID); err != nil {
			return
		}
		if existedUrl == url || existedUrl == "" {
			dals.SetUrlMap(service.rocksDB, urlID, url, time.Now().AddDate(1, 0, 0))
			return
		}
		digest = digest + 1
	}
}

func (service *Service) initHttpServer() {
	gin.SetMode(gin.ReleaseMode)
	service.httpEngine = gin.New()
	service.httpEngine.LoadHTMLFiles("pages/index.html")
	service.httpEngine.GET(
		"/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"host": service.conf.GetString("http.host"),
				"port": service.conf.GetInt("http.port"),
			})
		},
	)
	service.httpEngine.GET(
		"/:urlID", func(c *gin.Context) {
			urlID := c.Param("urlID")
			url, expiredTime, err := dals.GetUrl(service.rocksDB, urlID)
			log.Infof("get url from rocksdb by urlID, urlID: %v, url: %v, expiredTime: %v, err: %v", urlID, url, expiredTime, err)
			if err == nil && expiredTime.After(time.Now()) {
				if !strings.HasPrefix(url, "http") {
					url = "http://" + url
				}
				c.Redirect(http.StatusMovedPermanently, url)
			} else {
				c.Redirect(http.StatusMovedPermanently, "/")
			}
		},
	)
	service.httpEngine.GET(
		"/create/:url", func(c *gin.Context) {
			url := c.Param("url")
			log.Infof("create urlID, url: %v", url)
			if urlID, err := service.createUrlMap(url); err != nil {
				log.Errorf("failed to create url map, err: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"err": err,
				})
			} else {
				log.Infof("url: %v get urlID: %v", url, urlID)
				shortUrl := fmt.Sprintf("%v/%v", service.conf.GetString("http.host"), urlID)
				if service.conf.GetInt("http.port") != 80 {
					shortUrl = fmt.Sprintf("%v:%v/%v", service.conf.GetString("http.host"), service.conf.GetInt("http.port"), urlID)
				}
				c.JSON(http.StatusOK, gin.H{
					"short_url": shortUrl,
				})
			}

		},
	)
}

func (service *Service) Init() {
	service.initLogger()
	service.initRocksDB()
	service.initHttpServer()
}

func (service *Service) Start() {
	service.httpEngine.Run(fmt.Sprintf("0.0.0.0:%v", service.conf.GetString("http.port")))
}
