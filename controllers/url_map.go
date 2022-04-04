package controllers

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strings"

	"github.com/csugulo/ShortLinkServer/config"
	"github.com/csugulo/ShortLinkServer/consts"
	"github.com/csugulo/ShortLinkServer/db"
	"github.com/csugulo/ShortLinkServer/db/dals"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func stringDigest(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

func digest2str(digest uint32) string {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, digest)
	encoded := base64.StdEncoding.EncodeToString(bs)
	encoded = strings.ReplaceAll(strings.ReplaceAll(encoded, "=", ""), "/", "-")
	return encoded
}

func upsertUrlMap(url string) (urlID string, existed bool, err error) {
	digest := stringDigest(url)
	for {
		urlID = digest2str(digest)
		var existedUrl string
		if existedUrl, err = dals.GetUrl(db.RocksDB, urlID); err != nil {
			return
		}
		fmt.Println("GetUrl", urlID, existedUrl, err)

		if existedUrl == "" {
			go dals.SetUrlMap(db.RocksDB, urlID, url)
			return
		} else if existedUrl == url {
			existed = true
			return
		} else {
			digest = digest + 1
		}
	}
}

func CreateCommon(originUrl string, c *gin.Context) {
	var err error
	var urlID string
	if _, e := url.ParseRequestURI(originUrl); e != nil {
		err = fmt.Errorf("invalid url: %v", originUrl)
		go dals.AddLog(consts.Create, originUrl, urlID, consts.Failed, err.Error())
	} else {
		var existed bool
		urlID, existed, err = upsertUrlMap(originUrl)
		if !existed {
			go dals.AddLog(consts.Create, originUrl, urlID, consts.Success, "")
		}
	}
	if err == nil {
		log.Infof("url: %v get urlID: %v", originUrl, urlID)
		shortURL := fmt.Sprintf("%v:%v/%v", config.Conf.GetString("http.host"), config.Conf.GetInt("http.port"), urlID)
		c.JSON(http.StatusOK, gin.H{
			"short_url": shortURL,
		})
	} else {
		log.Errorf("url: %v failed to get urlID, err: %v", originUrl, err)
		c.JSON(http.StatusOK, gin.H{
			"err_message": err.Error(),
		})
	}
}

func CreateGet(c *gin.Context) {
	CreateCommon(c.Param("url"), c)
}

func CreatePost(c *gin.Context) {
	type Payload struct {
		Url string `json:"url"`
	}
	var payload Payload
	c.BindJSON(&payload)
	CreateCommon(payload.Url, c)
}

func Redirect(c *gin.Context) {
	urlID := c.Param("urlID")
	url, err := dals.GetUrl(db.RocksDB, urlID)
	if err == nil {
		log.Infof("redirect urlID: %v, to %v", urlID, url)
		c.Redirect(http.StatusMovedPermanently, url)
		go dals.AddLog(consts.Redirect, url, urlID, consts.Success, "")
	} else {
		log.Errorf("redirect urlID: %v failed", urlID)
		c.Redirect(http.StatusMovedPermanently, "/")
		go dals.AddLog(consts.Redirect, url, urlID, consts.Failed, err.Error())
	}
}

func Statistics(c *gin.Context) {
	statistics, err := dals.Statistics()
	if err != nil {
		log.Errorf("failed to get statistics, err: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"err_message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, statistics)
	}
}

func Echo(c *gin.Context) {
}
