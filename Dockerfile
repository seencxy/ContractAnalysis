# ==============================================
# Build Stage
# Go 版本可根据 go.mod 要求调整
# ==============================================
FROM golang:1.25-alpine AS builder

# 安装编译依赖
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源码
COPY . .

# 编译：静态链接，去除调试信息，优化体积
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o futures-analysis \
    main.go

# ==============================================
# Final Stage - 使用 scratch 极简镜像
# ==============================================
FROM scratch

# 从 builder 复制必要文件
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 复制编译好的二进制文件
COPY --from=builder /build/futures-analysis /app/futures-analysis

# 复制 Docker 专用配置
COPY config.docker.yaml /app/config.yaml

# 设置时区
ENV TZ=UTC

# 工作目录
WORKDIR /app

# 暴露端口
EXPOSE 8080 9090

# 启动应用
ENTRYPOINT ["/app/futures-analysis"]
