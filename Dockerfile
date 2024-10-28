# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /build

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o seed-sync .

# 最终阶段
FROM alpine:latest

WORKDIR /app

# 安装 tzdata 包
RUN apk add --no-cache tzdata

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /build/seed-sync .

# 只复制 .toml 配置文件
COPY --from=builder /build/config/*.toml /app/config/

# 确保二进制文件有执行权限
RUN chmod +x /app/seed-sync

# 暴露端口
EXPOSE 8705

# 运行应用
CMD ["./seed-sync"]