# Welcome To Fly

## 简介

一个简单而优雅的后端框架, 封装常用数据库组件及应用示例, 助力后端人员快速开发

[个人博客](https://blog.csdn.net/ywdhzxf/)

## 项目结构

参考 [Go程序布局](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)

```
├── README.md
├── config
|  ├── dev.json         // 各环境配置文件
|  └── config.go        // 配置初始化
|  └── model.go         // 配置结构体
├── api                 // 对外接口
|  ├── v1               // 版本号
|  |  ├── router.go     // 主路由
|  |  └── user          // 模块分组
|  |     ├── controller // 控制器
|  |     ├── service    // 接口逻辑
|  |     └── models     // 结构体存放
|  |     ├── router.go  // 模块子路由
├── rpc                 // grpc
|  ├── v1               // grpc版本
|  |  ├── router.go     // 主路由
|  |  └── user          // 模块分组
|  |     ├── controller // 控制器
|  |     ├── service    // 接口逻辑
|  |     └── models     // 结构体存放
|  |     └── pb         // proto文件
|  |     ├── router.go  // 模块子路由
├── internal            // 私有程序
|  ├── cache            // 缓存相关
|  ├── constvar         // 常量
|  └── monitor          // 监控定时服务相关
|  └── models           // 公用结构体
|  └── utils            // 公用方法(不能调用任何内部对象)
├── domain              // 数据库相关
|  ├── sqldb            // mysql相关, 包含对sql相关的操作
├── pkg                 // 安全导入的包(可以被任何项目直接导入使用)
|  ├── clickhouse       // ck组件
|  ├── email            // 邮件组件
|  ├── es               // es组件
|  ├── httpcode         // 请求处理组件
|  ├── jwt              // jwt组件
|  ├── logging          // 日志组件
|  ├── mongo            // mongo组件
|  └── mq               // mq组件
|  └── mysql            // mysql组件
|  └── redis            // redis组件
|  └── safego           // 安全运行组件
|  └── ws               // socket组件
├── go.mod              // 包管理    
├── main.go             // 入口文件     
├── Dockerfile          // Dockerfile     


├── cmd                 // 如果该项目有多个入口文件, 则根据业务存放在cmd里
|  ├── xxx-interface
|  |   ├── main.go      // 入口文件
|  |   ├── Dockerfile   // Dockerfile  
```

## 技术选型

### web框架

[iris 号称最快的Web框架](https://github.com/kataras/iris)

### rpc

[grpc](https://pkg.go.dev/google.golang.org/grpc)

### socket

[gorilla](https://github.com/gorilla/websocket)

### mysql

[gorm](https://gorm.io/)

### clickhouse

[dbr](https://github.com/mailru/dbr)

### es

[elastic](https://github.com/olivere/elastic/v6)

### mongo

[mongo](https://github.com/go-mgo/mgo/tree/v2)

### redis

[redis](https://github.com/go-redis/redis)

### mq

[amqp](https://github.com/streadway/amqp)

### log

[zap](https://pkg.go.dev/go.uber.org/zap)

### config

[viper](https://github.com/spf13/viper)

## 启动方式

### 命令行启动

```
go run main.go
```

### Docker启动

```
docker build -t fly:v1.0.0 .
docker run --rm -it -p 8888:8888 -p 9999:9999 --name fly fly:v1.0.0
```