package controllers

import (
	"net/http"

	"github.com/csugulo/ShortLinkWeb/config"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"domain": config.Conf.GetString("http.domain"),
		"port":   config.Conf.GetInt("http.port"),
	})
}
