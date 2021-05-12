# Welcome To Fly

## 项目结构

参考 [Go程序布局](!https://github.com/golang-standards/project-layout/blob/master/README_zh.md)

```
├── README.md
├── config
|  ├── dev.json         // 各环境配置文件
|  └── config.go        // 配置初始化
|  └── model.go         // 配置结构体
├── interface           // 对外接口
|  ├── v1               // api版本
|  |  ├── router.go     // 主路由
|  |  └── user          // 模块分组
|  |     ├── controller // 控制器
|  |     ├── service    // 接口逻辑
|  |     └── models     // 结构体存放
|  |     ├── router.go  // 模块子路由
├── internal            // 私有程序
|  ├── cache            // 缓存相关
|  ├── const            // 常量
|  └── monitor          // 监控定时服务相关
|  └── models           // 公用结构体
|  └── utils            // 公用方法(不能调用任何内部对象)
├── db                  // 数据库相关
|  ├── sqldb            // mysql相关, 包含对sql相关的操作
├── pkg                 // 安全导入的包(可以被任何项目直接导入使用)
|  ├── cache            // 缓存相关
|  └── monitor          // 监控定时服务相关
|  └── models           // 公用结构体
├── go.mod              // 包管理    
├── main.go             // 入口文件     
├── Dockerfile          // Dockerfile     


├── cmd                 // 如果该项目有多个入口文件, 则舍弃main文件根据业务存放在cmd里
|  ├── xxx-interface
|  |   ├── main.go      // 入口文件
|  |   ├── Dockerfile   // docker file  
```

## 技术选型

### web框架

[iris 号称最快的Web框架](!https://github.com/kataras/iris)

### mysql

[gorm](!https://gorm.io/)

### redis

[redis](!https://github.com/go-redis/redis)

### mq

[amqp](!https://github.com/streadway/amqp)

### log

[zap](!https://pkg.go.dev/go.uber.org/zap)

### config

[viper](!https://github.com/spf13/viper)