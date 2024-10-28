# 第一阶段：使用官方 Go 镜像进行编译
FROM dk.lixp.dev/golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 Go 模块和代码复制到镜像
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# 编译 Go 程序
RUN go build -o app server/cmd/main.go

# 第二阶段：创建一个更小的镜像
FROM dk.lixp.dev/alpine:latest

# 复制编译后的二进制文件
COPY --from=builder /app/app /app

# 设置运行的入口点
ENTRYPOINT ["/app"]
