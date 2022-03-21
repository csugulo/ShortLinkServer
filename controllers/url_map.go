package controllers

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"time"

	"github.com/csugulo/ShortLinkWeb/config"
	"github.com/csugulo/ShortLinkWeb/db"
	"github.com/csugulo/ShortLinkWeb/db/dals"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func createOrFindExistedUrlMap(url string) (urlID string, err error) {
	h := fnv.New32a()
	h.Write([]byte(url))
	digest := h.Sum32()
	for {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, digest)
		urlID = base64.StdEncoding.EncodeToString(bs)
		var existedUrl string
		if existedUrl, _, err = dals.GetUrl(db.RocksDB, urlID); err != nil {
			return
		}
		if existedUrl == url || existedUrl == "" {
			dals.SetUrlMap(db.RocksDB, urlID, url, time.Now().AddDate(1, 0, 0))
			return
		}
		digest = digest + 1
	}
}

func CreateCommon(originUrl string, c *gin.Context) {
	var err error
	var urlID string
	if _, e := url.ParseRequestURI(originUrl); e != nil {
		err = fmt.Errorf("invalid url: %v", originUrl)
	} else {
		urlID, err = createOrFindExistedUrlMap(originUrl)
	}
	if err == nil {
		log.Infof("url: %v get urlID: %v", originUrl, urlID)
		shortURL := fmt.Sprintf("%v:%v/%v", config.Conf.GetString("http.domain"), config.Conf.GetInt("http.port"), urlID)
		c.JSON(http.StatusOK, gin.H{
			"short_url": shortURL,
		})
	} else {
		log.Errorf("url: %v failed to get urlID, err: %v", originUrl, err)
		c.JSON(http.StatusOK, gin.H{
			"err_message": err,
		})
	}
}

func CreateGet(c *gin.Context) {
	CreateCommon(c.Param("url"), c)
}

type CreatePostPayload struct {
	Url string `json:"url"`
}

func CreatePost(c *gin.Context) {
	var payload CreatePostPayload
	CreateCommon(payload.Url, c)
}

func Redirect(c *gin.Context) {
	urlID := c.Param("urlID")
	url, expiredTime, err := dals.GetUrl(db.RocksDB, urlID)
	if err == nil && expiredTime.After(time.Now()) {
		log.Infof("redirect urlID: %v, to %v", urlID, url)
		c.Redirect(http.StatusMovedPermanently, url)
	} else {
		log.Errorf("urlID: %v is expired", urlID)
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}
