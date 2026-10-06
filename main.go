package main

import (
	"flag"
	"fmt"
	"go-bootstrap-demo/constants"
	"go-bootstrap-demo/services"
	"os"

	logger "github.com/hicode0101/go-logger"
	"github.com/kardianos/service"
)

var (
	appName     = ""
	serviceName = ""
	AppCmdS     = new(AppCmd)
)

func init() {
	appName = services.ConfigS.AppName
	serviceName = services.ConfigS.ServiceName

	logger.Info("Ready to start", serviceName)

	//Service 接口的主要方法 Run、Start、Stop、Restart、Install、Uninstall

	//配置服务的显示信息
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: appName,
		Description: services.ConfigS.AppName,
	}

	//执行程序的路径，如果不设置，则为当前程序
	//serviceConfig.Executable = "/opt/StartUp.sh"
	//注册到服务中，启动时需要带的参数

	workPath, _ := os.Getwd()
	serviceConfig.WorkingDirectory = workPath

	serviceConfig.Arguments = []string{"run"}

	var err error
	AppCmdS.Serv, err = service.New(AppCmdS, serviceConfig)
	if err != nil {
		logger.Error("", err)
	}

}

func main() {

	flag.Parse()

	fmt.Println("-----------------------")
	fmt.Println(logger.BrightGreen.Format("GoBootstrapDemo		" + constants.Version))
	fmt.Println(logger.BrightBlack.Format("[Go Bootstrap Demo]"))
	fmt.Println("-----------------------")

	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	logger.Info("GoBootstrapDemo starting ...")

	//支持的参数有 run , start ,  stop ,  restart ,  install ,  uninstall
	//fmt.Println(flag.Args())

	action := "run"
	args := flag.Args()
	if len(args) >= 1 {
		action = args[0]
	}

	logger.Debug("arg：", action)

	//识别执行参数
	err := AppCmdS.CmdHandler(action)
	if err != nil {
		panic(err)
		return
	}

}
