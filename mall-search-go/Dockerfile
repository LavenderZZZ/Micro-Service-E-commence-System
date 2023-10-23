# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.18 AS build-env

# 设置工作目录
WORKDIR /app

# 将 go module 和 go sum 文件复制并下载依赖项
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码到容器中
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用轻量级 alpine 镜像
FROM alpine:3.15

# 添加 ca-certificates，这可能是必要的，如果你的应用程序与使用 SSL 的外部服务进行通信
RUN apk --no-cache add ca-certificates

# 从构建环境复制构建的应用程序到当前容器
COPY --from=build-env /app/main /app/

# 设置工作目录
WORKDIR /app

# 设置启动命令
CMD ["./main"]

# 暴露端口
EXPOSE 8081
