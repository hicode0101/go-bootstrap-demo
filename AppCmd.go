package main

import (
	"go-bootstrap-demo/services"
	"sync"

	"github.com/hicode0101/go-logger"
	"github.com/hicode0101/go-utils"
	"github.com/kardianos/service"
)

type AppCmd struct {
	Serv service.Service
}

func (_self *AppCmd) CmdHandler(action string) error {
	logger.Debug("ServiceName", serviceName)
	logger.Debug("CmdHandler", action)
	var err error

	switch action {
	case "run":
		err = _self.Serv.Run()
		if err != nil {
			logger.Error("CmdHandler run error", err)
		}
	default:
		//要用管理员权限运行命令，才能成功注册为服务
		//支持的参数有 "start", "stop", "restart", "install", "uninstall"
		err = service.Control(_self.Serv, action)
		if err != nil {
			logger.Error("service.Control error", err)
		} else {
			logger.Info("CmdHandler success!", "action", action)
			if action == "install" {
				logger.Info("Installed path：", utils.File.GetCurrentExe())
			}
		}
	}

	return err
}

func (_self *AppCmd) Start(s service.Service) error {
	//启动服务时触发，建议不要有阻塞代码，执行真实的启动逻辑时要用异步执行
	logger.Debug("CmdS.Start", "serviceName", serviceName)
	go _self.run()
	return nil
}
func (_self *AppCmd) run() {
	//启动服务执行的逻辑，需要是阻塞的
	logger.Debug("CmdS.run", "serviceName", serviceName)

	//添加阻塞逻辑
	wg := sync.WaitGroup{}
	wg.Add(1)

	logger.Info("App started!")

	logger.Debug("WebServerAddr：", services.ConfigS.WebServerAddr)
	webServer := &WebServer{BindAddr: services.ConfigS.WebServerAddr}
	go webServer.StartWebServer()

	wg.Wait()
}

func (_self *AppCmd) Stop(s service.Service) error {
	//停止服务时触发，不要有阻塞代码
	logger.Debug("CmdS.Stop", "serviceName", serviceName)
	return nil
}
