# Recipe Generator API

一個基於 AI 的食譜生成 API 服務，提供食物識別、食材識別和食譜生成功能。

---

## 功能特點

- 🍳 智能食譜生成
- 🥗 食材與設備圖片辨識
- 📸 食物圖片辨識
- ⚡ 高性能與可擴展性
- 🔒 安全與穩定

---

## 技術架構

- **API 層** (`internal/api/`)
  - HTTP 處理器、路由、中間件
- **AI 服務** (`internal/ai/`)
  - OpenRouter 集成、請求隊列、提示詞處理
- **食譜服務** (`internal/recipe/`)
  - 食材/設備/食物辨識、食譜生成
- **圖片處理** (`internal/image/`)
  - 圖片優化、格式驗證、大小限制
- **快取系統** (`internal/cache/`)
  - 記憶體快取、LRU 策略、TTL 管理
- **監控與日誌** (`internal/metrics/`, `internal/common/`)
  - 健康檢查、日誌、性能指標

---

## 目錄結構

```
.
├── cmd/
│   └── api/                # 主程式入口
├── internal/
│   ├── api/               # API 層
│   ├── ai/                # AI 服務
│   ├── cache/             # 快取系統
│   ├── common/            # 通用工具
│   ├── config/            # 配置管理
│   ├── image/             # 圖片處理
│   ├── metrics/           # 監控指標
│   ├── recipe/            # 食譜服務
│   └── ...                # 其他模組
├── pkg/                   # 可重用套件
├── swagger.yaml           # OpenAPI 文件
├── Dockerfile
├── docker-compose.yml
├── .env
├── .env.example
└── README.md
```

---

## 快速開始

### 環境需求

- Go 1.21 或更高版本
- Docker (可選，用於容器化部署)

### 本地開發

1. 克隆倉庫：
```bash
git clone <repository-url>
cd recipe-generator
```

2. 設置環境變量：
```bash
cp .env.example .env
# 編輯 .env 文件，設置必要的環境變數
```

3. 運行服務：
```bash
go run cmd/api/main.go
```

服務將在 http://localhost:8080 上運行。

### Docker 部署

1. 構建鏡像：
```bash
docker build -t recipe-generator .
```

2. 運行容器：
```bash
docker run -p 8080:8080 --env-file .env recipe-generator
```

### Docker Compose

```bash
docker-compose up --build -d
# 預設服務於 http://localhost:8080
```

---

## API 文件與 Swagger 部署

本專案已提供完整的 OpenAPI (Swagger) 文件，詳見 `swagger.yaml`。

### 如何預覽 API 文件

- **線上預覽**：將 `swagger.yaml` 上傳至 [Swagger Editor](https://editor.swagger.io/) 直接瀏覽。
- **本地預覽**：
  ```sh
  docker run -p 8081:8080 -v $PWD/swagger.yaml:/swagger.yaml swaggerapi/swagger-ui
  # 然後瀏覽 http://localhost:8081
  ```

### 如何部署 Swagger UI（建議用於團隊協作或內部文件）

1. 將 `swagger.yaml` 放在專案根目錄或伺服器上。
2. 使用官方 Docker 映像部署：
   ```sh
   docker run -d -p 8081:8080 -v $PWD/swagger.yaml:/swagger.yaml swaggerapi/swagger-ui
   ```
3. 內網或雲端部署後，團隊可直接瀏覽 API 文件。

### 主要 API 端點

- `POST /api/v1/recipe/food` — 圖片辨識食物
- `POST /api/v1/recipe/ingredient` — 圖片辨識食材與設備
- `POST /api/v1/recipe/generate` — 使用食物名稱與偏好生成詳細新手友善食譜
- `POST /api/v1/recipe/suggest` — 使用食材與設備推薦適合的食譜

詳細請參考 `swagger.yaml` 內的 schema 與範例。

### 型別自動產生與驗證

- 可用 [oapi-codegen](https://github.com/deepmap/oapi-codegen) 產生 Go 型別與驗證。
- 或用 [swaggo/swag](https://github.com/swaggo/swag) 產生 Swagger UI（需在 handler 上加註解）。

---

## 配置說明

- `.env` 內可設定 API 金鑰、埠號、快取、限流等參數
- 圖片最大 5MB，支援 JPEG/PNG
- 服務預設監聽 8080 埠

---

## 監控與維護

- `/health` `/ready` `/live` — 健康檢查端點
- 日誌與錯誤追蹤
- 請求限流與快取策略

---

## 貢獻指南

1. Fork 專案
2. 創建功能分支
3. 提交更改
4. 發起合併請求

---

## 授權

MIT License

---

如需協助或有任何建議，歡迎提 issue 或聯絡作者！

## 環境變量配置

在 `.env` 文件中配置以下環境變量：

### 服務器配置
```
PORT=8080
ENV=development
SERVER_READ_TIMEOUT=10s
SERVER_WRITE_TIMEOUT=10s
SERVER_IDLE_TIMEOUT=120s
```

### 應用程式配置
```
APP_ENV=development
APP_DEBUG=true
LOG_LEVEL=info
APP_VERSION=1.0.0
APP_NAME=recipe-generator
```

### OpenRouter 配置
```
APP_OPENROUTER_API_KEY=your-api-key-here
APP_OPENROUTER_MODEL=google/gemini-2.0-flash-001
```

### 供應商配置
```
PROVIDER_ENABLED=false
PROVIDER_ONLY=  # 例如: Alibaba,OpenAI,Together
PROVIDER_IGNORE=  # 例如: Together
PROVIDER_ORDER=  # 例如: OpenAI,Alibaba,Together
PROVIDER_DATA_COLLECTION=deny  # deny 或 allow
```

### 模型參數配置
```
MODEL_TEMPERATURE=0.7
MODEL_MAX_TOKENS=2048
MODEL_TOP_P=0.9
MODEL_TOP_K=40
MODEL_PRESENCE_PENALTY=0.0
MODEL_FREQUENCY_PENALTY=0.0
```

### 圖片配置
```
MAX_IMAGE_SIZE=5242880
ALLOWED_IMAGE_TYPES=image/jpeg,image/png
```

### 快取配置
```
CACHE_ENABLED=true
CACHE_MAX_SIZE=1000
CACHE_TTL=1h
CACHE_CLEANUP_INTERVAL=10m
```

### 限流配置
```
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

### 隊列配置
```
QUEUE_WORKERS=5
QUEUE_MAX_SIZE=100
```

## 開發

### 代碼風格
- 使用 `gofmt` 格式化代碼
- 遵循 Go 標準代碼風格指南

### 測試
```bash
go test ./...
```

### 構建
```bash
go build -o recipe-generator ./cmd/api
```

