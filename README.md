# goGateway微服务网关
Gin best practices, gin development scaffolding, too late to explain, get on the bus.

基于gin_scaffold使用gin构建了企业级脚手架，代码简洁易读，可快速进行高效web开发。
主要功能有：
1. 请求链路日志打印，涵盖mysql/redis/request
2. 支持多语言错误信息提示及自定义错误提示。
3. 支持了多配置环境
4. 封装了 log/redis/mysql/http.client 常用方法
5. 支持swagger文档生成

项目地址：https://github.com/PlutoaCharon/goGateway
### 现在开始
- 安装软件依赖
```
git clone https://github.com/PlutoaCharon/goGateway.git
cd goGateway
go mod tidy
```
- 确保正确配置了 conf/mysql_map.toml、conf/redis_map.toml：

- 运行

```
go run main.go

➜  gin_scaffold git:(master) ✗ go run main.go
------------------------------------------------------------------------
[INFO]  config=./conf/dev/
[INFO]  start loading resources.
[INFO]  success loading resources.
------------------------------------------------------------------------
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /demo/index               --> github.com/PlutoaCharon/goGateway/controller.(*Demo).Index-fm (6 handlers)
[GIN-debug] GET    /demo/bind                --> github.com/PlutoaCharon/goGateway/controller.(*Demo).Bind-fm (6 handlers)
[GIN-debug] GET    /demo/dao                 --> github.com/PlutoaCharon/goGateway/controller.(*Demo).Dao-fm (6 handlers)
[GIN-debug] GET    /demo/redis               --> github.com/PlutoaCharon/goGateway/controller.(*Demo).Redis-fm (6 handlers)
 [INFO] HttpServerRun::8880
```
- 测试mysql与请求链路

### 文件分层
```
├── README.md
├── conf            配置文件夹
│   └── dev
│       ├── base.toml
│       ├── mysql_map.toml
│       └── redis_map.toml
├── controller      控制器
│   └── demo.go
├── dao             DB数据层
│   └── demo.go
├── docs            swagger文件层
├── dto             输入输出结构层
│   └── demo.go
├── go.mod
├── go.sum
├── main.go         入口文件
├── middleware      中间件层
│   ├── panic.go
│   ├── response.go
│   ├── token_auth.go
│   └── translation.go
├── public          公共文件
│   ├── log.go
│   ├── mysql.go
│   └── validate.go
└── router          路由层
    ├── httpserver.go
    └── route.go
```

### log / redis / mysql / http.client 常用方法

参考文档：https://github.com/e421083458/golang_common


### swagger文档生成

https://github.com/swaggo/swag/releases

- 下载对应操作系统的执行文件到$GOPATH/bin下面

如下：
```
➜  gin_scaffold git:(master) ✗ ll -r $GOPATH/bin
total 434168
-rwxr-xr-x  1 niuyufu  staff    13M  4  3 17:38 swag
```

- 设置接口文档参考： `controller/demo.go` 的 Bind方法的注释设置

```
// ListPage godoc
// @Summary 测试数据绑定
// @Description 测试数据绑定
// @Tags 用户
// @ID /demo/bind
// @Accept  json
// @Produce  json
// @Param polygon body dto.DemoInput true "body"
// @Success 200 {object} middleware.Response{data=dto.DemoInput} "success"
// @Router /demo/bind [post]
```

- 生成接口文档：`swag init`
- 然后启动服务器：`go run main.go`，浏览地址: http://127.0.0.1:8880/swagger/index.html