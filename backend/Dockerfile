# 使用最新版 alpine 镜像
FROM golang:alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装基础工具
RUN apk add --no-cache git

# 复制依赖
COPY go.mod go.sum ./
# 设置代理，防止下载超时
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

# 复制源码并构建
COPY . .
RUN go build -o main .

# 运行阶段
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
# 暴露端口
EXPOSE 8080
CMD ["./main"]