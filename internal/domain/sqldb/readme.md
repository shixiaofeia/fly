## sqlDb 存放 gorm结构体 & CRUD 方法

### 自动生成curd函数请使用 [autoCurd](./autocrud/main.go)

> 修改 obj 参数为gorm结构体即可自动生成相应函数, 复制到对应文件即可

### 目录建议

> 建议按照模块为维度, 比如

```
├── user                // 用户模块
|  ├── info             // 详情表
|  ├── log              // 操作日志表
|  ├── like             // 收藏表
├── risk                // 风控模块
|  ├── state            // 状态表
|  ├── history          // 历史表
```