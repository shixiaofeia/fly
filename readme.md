# Welcome To Fly

<p>
<a href="https://www.oscs1024.com/project/shixiaofeia/fly?ref=badge_small">
    <img src="https://www.oscs1024.com/platform/badge/shixiaofeia/fly.svg?size=small" alt="">
</a>
<a href="https://github.com/shixiaofeia/fly">
    <img src="https://badgen.net/badge/Github/fly?icon=github" alt="">
</a>
<a href="https://github.com/shixiaofeia/fly/LICENSE">
    <img alt="GitHub" src="https://img.shields.io/github/license/shixiaofeia/fly?style=flat-square">
</a>
<img src="https://img.shields.io/github/go-mod/go-version/shixiaofeia/fly.svg?style=flat-square" alt="">
<img alt="GitHub last commit" src="https://img.shields.io/github/last-commit/shixiaofeia/fly?style=flat-square">
<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/shixiaofeia/fly?style=social">
</p>

## 简介

一个简单而优雅的后端项目, 封装常用数据库组件及应用示例, 助力后端人员快速开发

[个人博客](https://blog.csdn.net/ywdhzxf/)

## 项目结构

参考 [Go程序布局](https://github.com/golang-standards/project-layout/blob/master/README_zh.md)

```
├── build               // 打包/集成
|  ├── app              // 应用程序名
|  |  ├── Dockerfile    // 集成的配置/脚本
├── cmd                 // 可执行目录
|  ├── app              // 应用程序名
|  |  ├── main.go       // 入口文件
├── configs             // 配置文件
|  ├── config.json      
├── doc                 // 项目文档
├── example             // 示例目录
├── internal            // 私有程序
|  ├── api              // 接口
|  ├── config           // 配置文件解析
|  ├── constvar         // 常量
|  ├── domain           // 表结构
|  ├── httpcode         // 请求处理组件
|  ├── kit              // 公用逻辑函数
|  └── monitor          // 监控定时服务相关
|  └── rpc              // rpc
├── logs                // 日志存放
├── pkg                 // 安全导入的包(可以被任何项目直接导入使用)
|  ├── clickhouse       // ck组件
|  ├── email            // 邮件组件
|  ├── es               // es组件
|  ├── kafka            // kafka组件
|  ├── jwt              // jwt组件
|  ├── libs             // 封装的公用方法
|  ├── logging          // 日志组件
|  ├── mongo            // mongo组件
|  └── mq               // mq组件
|  └── mysql            // mysql组件
|  └── redis            // redis组件
|  └── safego           // 安全运行组件
|  └── ws               // socket组件
├── .dockerignore       // docker忽略文件    
├── .gitignore          // git忽略文件    
├── go.mod              // 包管理    
├── README.md
```

## 优雅的代码

[请先参阅一遍官方的代码规范指南](https://github.com/golang/go/wiki/CodeReviewComments)

[gofmt](https://golang.org/cmd/gofmt/)

[goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)

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

### kafka

[kafka-go](https://github.com/segmentio/kafka-go)

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

> 请先复制 configs 目录下的配置文件, 并修改为自己的配置

### 命令行启动

```
go run cmd/app/main.go -config ./configs/config.yml
```

### Docker启动

```
docker build -f build/app/Dockerfile -t fly:v1.0.0 .
docker run --rm -it -p 8888:8888 -p 9999:9999 --name fly fly:v1.0.0
```