package controller

import (
	"go-bootstrap-demo/constants"
	"net/http"

	"github.com/gin-gonic/gin"
	utils "github.com/hicode0101/go-utils"
)

type DefaultCtrl struct {
}

func (_self *DefaultCtrl) RegisterRouter(router *gin.Engine) {

	router.GET("/", _self.Index)
	router.GET("/ping", _self.Ping)
	router.GET("/health", _self.Health)
	router.GET("/version", _self.Version)
}

// Index 欢迎页
func (_self *DefaultCtrl) Index(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", gin.H{
		"AppName":   "GoBootstrapDemo",
		"Version":   constants.Version,
		"Title":     "欢迎使用 GoBootstrapDemo",
		"SubTitle":  "基于 Gin 的 Go Web 脚手架示例项目",
		"BuildTime": utils.Date.CurrentTime(),
	})
}

func (_self *DefaultCtrl) Ping(c *gin.Context) {

	c.String(200, "pong")

}

func (_self *DefaultCtrl) Health(c *gin.Context) {

	c.String(200, "ok")

}

func (_self *DefaultCtrl) Version(c *gin.Context) {

	c.String(200, constants.Version)

}
