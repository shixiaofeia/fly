## logging 日志模块

目前支持两种编码模式, `console` 和 `json`; 可在配置文件中设置 `Encoder` 修改;

### console

```
[2022-06-21 18:39:40]   [info]  [main.main:36]  Start Web Server 
[2022-06-21 18:42:00]   [info]  [fly/example/api/example/controller.Hello:21]   [X-Request-Id:62b659cd-373a-49e2-82dc-e63cac629a84]    hello, 司马老贼
[2022-06-21 18:42:00]   [info]  [fly/pkg/httpcode.(*Req).Code:50]       [X-Request-Id:62b659cd-373a-49e2-82dc-e63cac629a84]    api: /v1/hello, run: 330.125µs, param: {"name":"司马老贼"}, code: 200
```

### json

```
{"level":"info","ts":"2022-06-21 18:43:09","caller":"main.main:36","msg":"Start Web Server "}
{"level":"info","ts":"2022-06-21 18:43:13","caller":"fly/example/api/example/controller.Hello:21","msg":"hello, 司马老贼","X-Request-Id":"2e33e08f-4fe9-4424-aeb5-c26151559143"}
{"level":"info","ts":"2022-06-21 18:43:13","caller":"fly/pkg/httpcode.(*Req).Code:50","msg":"api: /v1/hello, run: 305.458µs, param: {\"name\":\"司马老贼\"}, code: 200","X-Request-Id":"2e33e08f-4fe9-4424-aeb5-c26151559143"}
```
