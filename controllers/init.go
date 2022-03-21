package controllers

import "github.com/gin-gonic/gin"

var Server *gin.Engine

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	Server = gin.New()
	Server.LoadHTMLFiles("pages/index.html")

	Server.GET("/:urlID", Redirect)
	Server.GET("/create/:url", CreateGet)
	Server.POST("/create", CreatePost)
	Server.GET("/", Index)
}
