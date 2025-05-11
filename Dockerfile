# 使用官方 Go 映像作為基礎
FROM golang:1.24-alpine AS builder

# 設置工作目錄
WORKDIR /app

# 安裝必要的系統依賴
RUN apk add --no-cache git gcc musl-dev

# 複製 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製應用程式代碼
COPY . .

# 創建輸出目錄
RUN mkdir -p /app/bin

# 編譯應用程式（靜態連結）
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/server ./cmd/main.go

# 使用輕量級的 alpine 作為最終映像
FROM alpine:latest

# 安裝 CA 證書（用於 HTTPS 請求）
RUN apk --no-cache add ca-certificates

WORKDIR /app

# 從 builder 階段複製編譯好的執行檔和配置文件
COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/.env.example .env

# 暴露端口
EXPOSE 8080

# 運行應用程式
CMD ["/app/server"]

