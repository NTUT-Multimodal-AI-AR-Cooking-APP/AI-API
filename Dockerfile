# 構建階段
FROM golang:1.21-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 安裝必要的構建工具
RUN apk add --no-cache git

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -o recipe-generator ./cmd/api

# 運行階段
FROM alpine:latest

# 安裝必要的運行時依賴
RUN apk add --no-cache ca-certificates tzdata

# 設置時區
ENV TZ=Asia/Taipei

# 設置工作目錄
WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/recipe-generator .

# 暴露端口
EXPOSE 8080

# 設置環境變量
ENV APP_ENV=production

# 運行應用
CMD ["./recipe-generator"]

