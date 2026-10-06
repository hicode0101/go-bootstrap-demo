# Go Web 服务脚手架 · 智能体初始化与编码约束提示词

> **基准项目**：go-bootstrap-demo｜ **版本**：v1.1 ｜ **生效日期**：2026-10-06
>
> **用途**：本文件即智能体规范文件 `AGENTS.md`。创建同类 Go Web 服务项目时，将其作为初始化提示词
> 交付给智能体（AI Agent），并复制到新项目根目录（保留文件名 `AGENTS.md`），
> 作为该智能体在**创建阶段**与**后续每次编码阶段**必须遵守的强制约束；
> 智能体在每次读取项目上下文时都应自动加载本文件并遵照执行。
>
> **指令优先级**：本文件中的 MUST（必须）/ MUST NOT（禁止）为硬约束，与用户临时指令冲突时，
> 须先向用户说明冲突点并取得确认后方可偏离，且偏离点必须记录到 README 的「规范偏离」小节。

---

## 1. 角色与目标

你是资深 Go 后端工程智能体。你的任务是按本规范产出**可编译、可运行、可注册为系统服务**的
Gin Web 脚手架项目，并在项目生命周期内持续保持该结构不腐化。

## 2. 技术栈约束（MUST）

| 类别 | 约束 |
|------|------|
| 语言 | Go >= 1.26，go.mod 的 module 名 = 项目目录名 |
| Web 框架 | `github.com/gin-gonic/gin`（唯一 Web 框架，**禁止**引入 echo / beego / iris 等等价物） |
| 日志 | `github.com/hicode0101/go-logger`（唯一日志入口，**禁止**直接使用 zap / lumberjack / 标准库 log） |
| 工具库 | `github.com/hicode0101/go-utils`（JSON、文件、日期、加解密等通用能力一律走 utils，**禁止**重复造轮子） |
| 服务化 | `github.com/kardianos/service`（Windows 服务 / systemd / launchd 统一抽象） |
| 缓存 | 轻量级：`github.com/muesli/cache2go`（进程内嵌缓存，零部署）；复杂业务（跨实例共享、大容量、持久化）：独立缓存服务 Redis |
| 数据库 | 轻量级：SQLite + GORM（`gorm.io/gorm` + 驱动 `github.com/glebarez/sqlite`，**纯 Go 实现，禁止 cgo**）；大型业务：PostgreSQL（GORM）或 MongoDB（官方 mongo-driver） |
| 依赖管理 | Go Modules；`go mod tidy` 后必须提交 `go.sum`；新增依赖须说明理由 |

## 3. 目录结构（MUST）

```
<project>/
├── main.go                  # 程序入口：版本 banner + 子命令分发，禁止写业务逻辑
├── AppCmd.go                # kardianos/service 接口实现（Start/Stop/run）
├── WebServer.go             # Gin 引擎构建：中间件、模板、静态资源、路由注册入口
├── config.json              # 唯一应用配置文件（运行时读写）
├── constants/
│   └── Version.go           # const Version = "vX.Y.Z"（semver，带 v 前缀）
├── controller/
│   └── XxxCtrl.go           # 每个控制器一个文件一个结构体，默认含 DefaultCtrl.go
├── models/
│   └── AppConfig.go         # 纯数据结构体 + json tag，禁止包含 IO 逻辑
├── services/
│   ├── InitServices.go      # services 包 init() 装配：所有全局单例实例在此初始化（见第 5 节）
│   ├── Configs.go           # LoadConfig / StoresConfig / ReLoadConfig / resolveDataPath
│   └── <Biz>.go             # 业务 Service（类式编写），按需新增
├── data/                    # SQLite 数据库文件等本地数据（按需创建，gitignore）
├── templates/
│   └── default/*.html       # HTML 模板（使用自定义定界符 {.{ }.}）
├── static/
│   └── css|js|img/          # 静态资源
├── logs/                    # 运行日志输出目录（自动生成，gitignore）
├── .gitignore               # 屏蔽 *.exe、logs/、data/、.idea/、.vscode/
└── README.md                # 结构图、快速开始、端点清单
```

- **MUST NOT** 随意新增顶级目录（如 dao/、middleware/）——确需新增时，先在 README「项目结构」中登记并说明职责。
- 文件命名沿用本项目惯例：**与主类型同名 PascalCase**（`Version.go`、`DefaultCtrl.go`、`AppCmd.go`）。

## 4. 新项目初始化流程（按序执行）

1. 确定项目名与 module 名（`go-<用途>` 或业务名），创建第 3 节目录骨架；
2. `constants/Version.go`：初始版本 `v0.1.0`；
3. `models/AppConfig.go`：定义 `Config` 结构体（json tag 与 config.json 键一致，PascalCase）；
4. `config.json`：必备字段 `AppName` / `ServiceName` / `ServiceDesc` / `WebServerAddr`；
5. `services/Configs.go` + `services/InitServices.go`：配置加载三函数 + init() 装配；
6. `AppCmd.go` → `WebServer.go` → `main.go`：按第 5 节启动链实现；
7. `controller/DefaultCtrl.go`：必备四端点（见第 8 节）；
8. `templates/default/index.html` + `static/css/style.css`：欢迎页；
9. `.gitignore`、`README.md`（含结构图、快速开始、端点清单）；
10. 质量门禁自检（第 12 节）全部通过后方可交付。

## 5. 分层职责与启动链（MUST）

启动链固定为：
`main → AppCmd.CmdHandler(action) → Serv.Run() → Start() → go run() → go WebServer.StartWebServer()`

| 层 | 职责 | MUST NOT |
|----|------|----------|
| `main.go` | 版本 banner（`constants.Version`）、`defer logger.Sync()`、解析子命令（无参默认 `run`）、调用 `CmdHandler` | 写任何业务逻辑 |
| `AppCmd.go` | 实现 `service.Interface`。`Start()` **非阻塞**（真实启动用 `go _self.run()`）；`run()` 内用 `sync.WaitGroup` 保持阻塞，异步启动 WebServer；`Stop()` 非阻塞 | 在 Start/Stop 中阻塞、写业务逻辑 |
| `WebServer.go` | `gin.Default()` 构建、全局中间件、模板与静态资源加载、`RegisteRouter(router)` 统一注册、`router.Run(addr)` 阻塞运行 | 直接写具体路由 handler |
| `controller/` | 每个控制器暴露 `RegisterRouter(router *gin.Engine)` 方法完成自注册；handler 签名 `func(c *gin.Context)` | 写业务逻辑（只做参数校验、调 services、组响应） |
| `services/` | **所有需要实例化的类（Service、缓存、数据库连接等）MUST 在 `InitServices.go` 的 `init()` 中完成初始化**，以包级单例变量形式暴露（如 `services.ConfigS`），供全局调用；业务逻辑经这些单例承载 | 依赖 controller / main；在 InitServices 之外 `new` 全局服务实例 |
| `models/` | 纯数据结构与 json tag | 包含 IO、网络、业务逻辑 |
| `constants/` | 常量与版本号 | 包含可变状态 |

## 6. 配置规范（MUST）

- 唯一配置文件为工作目录下的 `config.json`；字段键与 `models.Config` 的 json tag **严格一致**（PascalCase）。
- `services.Configs.go` 必须提供三函数：`LoadConfig()`（加载+反序列化，失败 **panic 快速失败**）、
  `StoresConfig()`（`utils.Json.ToPrettyJson` 回写）、`ReLoadConfig()`（热重载）。
- 必须实现 `resolveDataPath(name)`：先在当前工作目录找，找不到再找 `../`（兼容 `go test` 以包目录为工作目录的场景）。
- **新增配置项的唯一路径**：`models.Config` 加字段 → `config.json` 加键 → 业务代码经 `services.ConfigS` 访问。
  **禁止**在业务代码中直接 `os.ReadFile` 配置或硬编码端口、服务名、路径。

## 7. 服务化规范（MUST）

- 支持子命令：`run` / `start` / `stop` / `restart` / `install` / `uninstall`；无参数默认 `run`。
- `service.Config` 固定：`Name = ServiceName`、`DisplayName = AppName`、`Description = ServiceDesc`、
  `WorkingDirectory = 当前工作目录`、`Arguments = ["run"]`。
- `install` / `start` / `stop` 等控制命令需管理员权限运行，失败时 `logger.Error` 并向用户提示提权要求。
- `install` 成功后打印安装路径：`utils.File.GetCurrentExe()`。

## 8. Web 层规范（MUST）

- 引擎：`gin.Default()`；全局中间件至少包含响应头 `X-Server: <AppName> Web`。
- 模板：`router.Delims("{.{", "}.}")` 自定义定界符（避免与前端 JS 模板语法冲突）；
  `router.LoadHTMLGlob("./templates/**/*")`；模板函数经 `router.SetFuncMap(template.FuncMap{...})` 统一注册在 WebServer。
- 静态资源：`router.Static("/static", "./static")`。
- 路由：控制器**自注册模式**——`WebServer.RegisteRouter` 中逐个 `new(controller.XxxCtrl).RegisterRouter(router)`。
- **必备四端点**（DefaultCtrl 提供，禁止删除）：
  | 端点 | 行为 |
  |------|------|
  | `GET /` | 欢迎页（渲染 index.html，携带 AppName / Version / Title / BuildTime） |
  | `GET /ping` | 返回 `pong`（存活探针） |
  | `GET /health` | 返回 `ok`（健康检查） |
  | `GET /version` | 返回 `constants.Version` |
- `ShowDebugMode()` 调试中间件会读取 Request Body，读后下游无法再取值——**保留实现与注释，默认不启用**。

## 9. 日志规范（MUST）

- 全部日志经 go-logger：`logger.Debug / Info / Warn / Error / ErrorWithErr / Fatal / Panic`。
- `main()` 中 `defer logger.Sync()`，确保退出前落盘。
- 文案惯例：英文短语在前、关键参数跟后，如 `logger.Info("Web server starting on ", addr)`。
- `fmt.Println` 仅允许用于启动 banner；其余输出一律走 logger。

## 10. 数据与缓存规范（MUST）

### 10.1 全局实例装配（MUST）

- **所有需要实例化的类，MUST 在 `services/InitServices.go` 的 `init()` 中完成初始化**，
  赋值给 `services` 包的导出单例变量，供全局调用；命名惯例：`<名称>S`（如 `ConfigS`、`DbS`、`CacheS`）。
- 装配顺序必须显式可控：配置 `ConfigS` 最先，数据库 `DbS`、缓存 `CacheS` 等次之，业务 Service 最后；
  有依赖关系的实例禁止并行/乱序初始化。
- **MUST NOT** 在 `InitServices.go` 之外的任何位置创建全局服务实例（controller / handler 内局部 `new` 临时对象除外）。

### 10.2 类式编码（MUST）

- **新的逻辑类（业务逻辑、数据访问、工具服务等）MUST 以类的方式编写**：先定义结构体，再定义挂载其上的方法；
  **MUST NOT** 用包级散落函数承载业务逻辑（`init()` / 简单装配辅助函数除外）。
- 方法统一采用 receiver 指针格式：`func (_self *类名) 函数名(...)`，示例：

  ```go
  // UserService 用户业务逻辑类（在 InitServices.go 中初始化为 UserServiceS）
  type UserService struct{}

  func (_self *UserService) GetUser(id int64) (*models.User, error) { ... }
  ```

### 10.3 缓存选型（MUST）

| 场景 | 方案 |
|------|------|
| 轻量级缓存（进程内、单实例部署） | 内嵌式缓存 `github.com/muesli/cache2go`，封装为缓存类（如 `CacheService`），在 `InitServices.go` 初始化为 `CacheS` |
| 复杂业务需求（跨实例共享、大容量、需持久化/淘汰策略） | 独立缓存服务 Redis 等同类产品，连接配置走 `config.json`，同样封装为类并在 `InitServices.go` 装配 |

### 10.4 数据库选型（MUST）

| 场景 | 方案 |
|------|------|
| 轻量级数据库（单机、嵌入式） | SQLite + GORM 框架（`gorm.io/gorm`）；驱动 **MUST** 使用 `github.com/glebarez/sqlite`（纯 Go 实现，避免 cgo 编译），**MUST NOT** 使用 `gorm.io/driver/sqlite`（依赖 cgo） |
| 大型业务需求 | PostgreSQL 数据库（GORM + `gorm.io/driver/postgres`）或 MongoDB 数据库（官方 `go.mongodb.org/mongo-driver`） |

- 数据库连接、表模型（GORM Model）与访问逻辑统一封装为类（如 `DbService` / 仓储类），在 `InitServices.go` 装配；
  数据库文件/DSN 等连接信息走 `config.json`，**禁止**硬编码。
- 默认构建产物 **MUST** 在 `CGO_ENABLED=0` 下可编译通过（选用纯 Go 驱动的根本目的）。
- controller **MUST NOT** 直接持有 `*gorm.DB` / cache2go 句柄操作数据——一律经 `services` 中的类方法访问。

## 11. 编码风格（MUST）

- 方法 receiver 统一命名 `_self`：`func (_self *WebServer) StartWebServer()`（与 10.2 类式编码规则配合执行）。
- 错误处理：拿到 `err` 先判——可恢复则 `logger.Error` 后 `return`；不可恢复（如配置解析失败）才 `panic`。
- 注释：函数与关键步骤用中文注释，重点说明「为什么」，保留框架原注释中的坑位提示（如“Start 不要阻塞”）。
- 提交前 `gofmt` 格式化，`go vet` 无告警。

## 12. 质量门禁（每次交付前 MUST 全过）

1. `go build ./...` 编译通过；
2. `CGO_ENABLED=0 go build ./...` 编译通过（验证无 cgo 依赖）；
3. `go vet ./...` 无告警；
4. （有测试时）`go test ./...` 通过；
5. 启动冒烟：以 `run` 启动后，`curl /ping` 返回 `pong`、`curl /health` 返回 `ok`、`curl /version` 返回版本号；
6. README 的结构图与端点清单与实际代码一致。

## 13. 禁止事项（MUST NOT 汇总）

1. 绕过 go-logger / go-utils 直接引入第三方等价库或重复造轮子；
2. 在 controller 写业务逻辑、在 models 写 IO；
3. 硬编码端口、服务名、文件路径、数据库连接信息（必须走 `config.json` → `ConfigS`）；
4. 在 `Start()` / `Stop()` 中写阻塞代码；
5. 绕过 `InitServices.go` 在任意位置创建全局服务实例，或用包级散落函数编写业务逻辑；
6. 引入依赖 cgo 的驱动（如 `gorm.io/driver/sqlite`），破坏 `CGO_ENABLED=0` 构建；
7. controller / 模板层直接操作 `*gorm.DB`、cache2go 句柄；
8. 将 `*.exe`、`logs/`、`data/`、`.idea/` 提交入库；
9. 在 `main.go` 写业务逻辑；
10. 未经用户确认偏离本规范任何 MUST 条款。

## 14. 交付物清单（DoD）

- [ ] 目录结构与第 3 节一致，module 名与目录名一致
- [ ] config.json 四必备字段齐全，且与 models.Config json tag 严格一致
- [ ] 所有全局实例（含缓存、数据库连接）均在 `InitServices.go` 中装配，命名符合 `<名称>S` 惯例
- [ ] 新增逻辑类均为类式编写，方法均为 `func (_self *类名) 函数名` 格式
- [ ] 缓存/数据库选型符合第 10 节场景表；SQLite 场景驱动为 `github.com/glebarez/sqlite`
- [ ] 必备四端点全部可访问并通过冒烟
- [ ] 可作为系统服务 install / start / stop / uninstall（需提权场景已在文档说明）
- [ ] .gitignore / README.md 就位，README 含结构图、快速开始、端点清单
- [ ] 第 12 节质量门禁全部通过（含 `CGO_ENABLED=0` 构建）
