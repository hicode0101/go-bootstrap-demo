package main

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"go-bootstrap-demo/controller"

	logger "github.com/hicode0101/go-logger"

	"github.com/gin-gonic/gin"
)

type WebServer struct {
	BindAddr string
}

func (_self *WebServer) StartWebServer() {

	router := gin.Default()

	// 全局中间件：给所有响应添加 X-Server 头
	router.Use(func(c *gin.Context) {
		c.Header("X-Server", "GoBootstrapDemo Web")
		c.Next()
	})

	router.Delims("{.{", "}.}")
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": _self.formatAsDate,
		"GetMapVal":    _self.GetMapVal,
	})
	router.LoadHTMLGlob("./templates/**/*")

	//静态资源目录
	router.Static("/static", "./static")

	_self.RegisteRouter(router)
	logger.Debug("---------------------------------")

	logger.Info("Web server starting on ", _self.BindAddr)

	//这里会阻塞
	err := router.Run(_self.BindAddr)
	if err != nil {
		logger.Error("Web server start error", err)
	}
}

func (_self *WebServer) RegisteRouter(router *gin.Engine) {
	new(controller.DefaultCtrl).RegisterRouter(router)
}

func (_self *WebServer) formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func (_self *WebServer) GetMapVal(mapData map[string]interface{}, key string) interface{} {
	if val, ok := mapData[key]; ok {
		//存在
		return val
	}
	return nil
}

func (_self *WebServer) ShowDebugMode() gin.HandlerFunc {
	return func(context *gin.Context) {

		logger.Debug("Request Path：", context.FullPath(), " Method：", context.Request.Method, " Ip：", context.ClientIP())
		logger.Debug("Request Query：")
		for k, v := range context.Request.URL.Query() {
			logger.Debug(k, " = ", v)
		}

		logger.Debug("Request Context：")
		for k, v := range context.Keys {
			logger.Debug(k, " = ", v)
		}

		//Request.Form 默认放到了 Request.Body
		logger.Debug("Request Form：")
		//从body加载到Form对象
		context.Request.ParseForm()
		for k, v := range context.Request.PostForm {
			logger.Debug(k, " = ", v)
		}

		//这里会对上传文件和post body,  post form有影响
		//这里读取了，后面的逻辑就取不了值了
		body := context.Request.Body
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			logger.Warn("Request Body empty.")
		} else {
			logger.Debug("Request Body：", string(bodyBytes))
		}

		context.Next()

	}
}
