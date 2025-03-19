# 使用 Go 语言官方镜像作为构建环境
FROM golang:1.23.0 AS builder

# 设置工作目录
WORKDIR /app

# 安装编译依赖
RUN apt-get update && apt-get install -y gcc make

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 下载依赖（如果有）
RUN go mod tidy

# 将当前目录下的所有代码复制到容器的工作目录中
COPY . .

# 编译 Go 项目
RUN make static-build && ls

# 使用更小的基础镜像来运行应用
FROM registry.cn-shanghai.aliyuncs.com/taochen/alpine:latest

RUN apk add --no-cache curl bash wget ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件到运行环境
COPY --from=builder /app/main .

RUN ls

# 暴露应用监听的端口
EXPOSE 8080
VOLUME /root/conf
VOLUME /root/storage
VOLUME /root/static

# 设置启动命令
CMD ["/bin/sh", "-c", "pwd && /root/main"]
