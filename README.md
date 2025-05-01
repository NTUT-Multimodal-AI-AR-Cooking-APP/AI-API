# 食譜生成器 API

這是一個使用 Go 和 Gemini AI 開發的智慧食譜生成器 API 服務。只要提供您家中現有的廚具和食材，API 就能為您生成合適的食譜建議。
## 測試
我有提供一個`swiftExsemple.wsift`前端測試範例可以參考
## 系統功能

1. **智慧配對**
   - 根據現有廚具類型生成適合的烹飪方式
   - 依據食材組合推薦最佳烹飪方案

2. **彈性輸入**
   - 支援多種廚具組合
   - 可接受不同份量的食材清單
   - 可設定個人烹飪偏好

3. **標準化輸出**
   - 清晰的步驟說明
   - 精確的烹飪時間
   - 詳細的火候控制
   - 完整的成品說明

## 使用 Docker 快速部署

### 1. 環境準備
```bash
# 複製環境設定檔
cp .env.example .env

# 編輯 .env 檔案，設定您的 Gemini API 金鑰
# GEMINI_API_KEY=您的API金鑰
```

### 2. 建立與運行
```bash
# 建立 Docker 映像檔
docker build -t recipe-generator .

# 運行容器（前台模式）
docker run -p 8080:8080 recipe-generator

# 或使用背景模式運行
docker run -d -p 8080:8080 recipe-generator
```

### 3. API 使用範例

發送 POST 請求到 `http://localhost:8080/generate-recipe`：

```json
{
    "equipment": [
        {
            "name": "平底鍋",
            "type": "鍋具",
            "size": "中型",
            "material": "不鏽鋼"
        }
    ],
    "ingredients": [
        {
            "name": "油",
            "type": "食材",
            "amount": "2湯匙",
            "unit": "湯匙"
        }
    ],
    "preference": {
        "cooking_method": "煎",
        "doneness": "中等熟"
    }
}
```

您將收到包含完整食譜的 JSON 回應：
```json
{
    "dish_name": "菜名",
    "dish_description": "菜餚描述",
    "recipe": [
        {
            "step": "步驟說明",
            "time": "所需時間",
            "temperature": "溫度",
            "description": "描述",
            "doneness": "熟度"
        }
    ]
}
```

### 4. 常見問題排解

1. 如果無法連接 API：
   - 確認容器是否正常運行：`docker ps`
   - 檢查日誌：`docker logs <container_id>`

2. 如果需要停止服務：
   - 查看容器 ID：`docker ps`
   - 停止容器：`docker stop <container_id>`

3. 如果需要更新服務：
   - 停止舊容器
   - 重新建立映像檔
   - 啟動新容器
