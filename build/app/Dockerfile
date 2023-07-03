FROM golang:1.18-alpine as builder

WORKDIR /data/fly
COPY ../.. /data/fly/

# 打包二进制&&增加执行权限
RUN export GO111MODULE=on \
    && export GOPROXY=https://goproxy.io \
    && go mod tidy \
    && export GOARCH=amd64 \
    && export GOOS=linux \
    && go build -o flyServer cmd/app/main.go \
    && chmod +x flyServer

FROM alpine

#设置东八区，北京时间
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add tzdata

WORKDIR /data/fly
COPY --from=builder /data/fly /data/fly

# 容器入口, 执行命令
CMD ["./flyServer", "-config", "./configs/config.yml"]