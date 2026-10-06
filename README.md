# go-bootstrap-demo

参照 HiDnsServer 的项目结构与代码框架创建的 Go Web 脚手架示例项目，实现了一个简单的 Web 欢迎界面。

## 项目结构

```
go-bootstrap-demo/
├── main.go                  # 程序入口：版本 banner + 启动 Web 服务
├── WebServer.go             # Gin Web 服务器：模板加载、静态资源、路由注册
├── config.json              # 应用配置（AppName / WebServerAddr / RunMode）
├── constants/
│   └── Version.go           # 版本常量
├── controller/
│   └── DefaultCtrl.go       # 默认控制器：欢迎页 + /ping /health /version
├── logger/
│   ├── ConsoleColor.go      # 控制台颜色输出
│   └── applogs.go           # zap + lumberjack 日志（文件 + 控制台双写）
├── models/
│   └── AppConfig.go         # 配置结构体
├── services/
│   ├── Configs.go           # 配置文件加载
│   └── InitServices.go      # services 包 init() 全局初始化
├── templates/
│   └── default/
│       └── index.html       # 欢迎页模板（自定义定界符 {..{ }.}）
├── static/
│   └── css/style.css        # 欢迎页样式
└── logs/                    # 运行日志输出目录（自动生成）
```

## 快速开始

```bash
go mod tidy
go run .
```

浏览器访问 <http://127.0.0.1:8080> 查看欢迎页。

## 接口列表

| 路由      | 说明                 |
| --------- | -------------------- |
| `/`       | 欢迎页（HTML 模板）  |
| `/ping`   | 存活检测，返回 pong  |
| `/health` | 健康检测，返回 ok    |
| `/version`| 返回当前版本号       |

## 配置说明

编辑 `config.json`：

```json
{
    "AppName": "GoBootstrapDemo",
    "ServiceName": "gobootdemo",
    "WebServerAddr": ":8080",
    "RunMode": "debug"
}
```
