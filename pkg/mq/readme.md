## mq连接支持

> 使用请参考 [test文件](index_test.go)

### 注意

1. mq的channel是有数量限制的(默认2047), 超过限制之后无法创建新的channel, 会报错(channel id space exhausted), 所以建议单个服务不要过多创建新的channel;
2. 建议项目中频率不高的生产者共用一个channel, 消费者的话建议一个消费者对应一个channel(1 consumer per channel);
