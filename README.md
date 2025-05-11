# AI 食譜生成器 API

這是一個基於 AI 的食譜生成器 API，提供食物圖片辨識、食材設備辨識、食譜生成等功能。

## 功能特點

- 食物圖片辨識：分析食物圖片，提供食物名稱、描述、可能使用的食材和設備
- 食材設備辨識：分析圖片中的食材和設備，提供詳細的清單和摘要
- 食譜生成：根據食物名稱生成詳細的食譜
- 食譜建議：根據可用的食材和設備，提供多個適合的食譜建議

## API 端點

### 1. 食物圖片辨識

分析食物圖片，提供詳細的食物資訊。

**端點：** `POST /api/recognize-food`

**請求格式：**
```json
{
    "image": "data:image/jpeg;base64,...",
    "description_hint": "可選的描述提示"
}
```

**回應格式：**
```json
{
    "recognized_foods": [
        {
            "name": "食物名稱",
            "description": "詳細的食物描述",
            "possible_ingredients": [
                {
                    "name": "食材名稱",
                    "type": "食材分類"
                }
            ],
            "possible_equipment": [
                {
                    "name": "設備名稱",
                    "type": "設備分類"
                }
            ]
        }
    ]
}
```

### 2. 食材設備辨識

分析圖片中的食材和設備，提供詳細的清單。

**端點：** `POST /api/recognize-ingredients`

**請求格式：**
```json
{
    "image": "data:image/jpeg;base64,...",
    "description_hint": "可選的描述提示"
}
```

**回應格式：**
```json
{
    "ingredients": [
        {
            "name": "食材名稱",
            "type": "食材分類",
            "amount": "數量",
            "unit": "單位",
            "preparation": "處理方式"
        }
    ],
    "equipment": [
        {
            "name": "設備名稱",
            "type": "設備分類",
            "size": "大小",
            "material": "材質",
            "power_source": "能源類型"
        }
    ],
    "summary": "整體摘要說明"
}
```

### 3. 食譜生成

根據食物名稱生成詳細的食譜。

**端點：** `POST /api/generate-recipe`

**請求格式：**
```json
{
    "dish_name": "菜名",
    "preferred_ingredients": ["食材1", "食材2"],
    "excluded_ingredients": ["食材3", "食材4"],
    "preferred_equipment": ["設備1", "設備2"],
    "preference": {
        "cooking_method": "烹調方式",
        "doneness": "熟度",
        "serving_size": "份量"
    }
}
```

**回應格式：**
```json
{
    "dish_name": "菜名",
    "dish_description": "菜餚描述",
    "ingredients": [
        {
            "name": "食材名稱",
            "type": "食材分類",
            "amount": "數量",
            "unit": "單位",
            "preparation": "處理方式"
        }
    ],
    "equipment": [
        {
            "name": "設備名稱",
            "type": "設備分類",
            "size": "大小",
            "material": "材質",
            "power_source": "能源類型"
        }
    ],
    "recipe": [
        {
            "step_number": 1,
            "title": "步驟標題",
            "description": "步驟說明",
            "actions": [
                {
                    "action": "動作名稱",
                    "tool_required": "使用工具",
                    "material_required": ["使用材料"],
                    "time_minutes": 整數分鐘,
                    "instruction_detail": "詳細操作方法"
                }
            ],
            "estimated_total_time": "預估時間",
            "temperature": "溫度或火力",
            "warnings": ["注意事項"],
            "notes": "補充備註"
        }
    ]
}
```

### 4. 食譜建議

根據可用的食材和設備，提供多個適合的食譜建議。

**端點：** `POST /api/suggest-recipes`

**請求格式：**
```json
{
    "available_ingredients": [
        {
            "name": "食材名稱",
            "type": "食材分類",
            "amount": "數量",
            "unit": "單位",
            "preparation": "處理方式"
        }
    ],
    "available_equipment": [
        {
            "name": "設備名稱",
            "type": "設備分類",
            "size": "大小",
            "material": "材質",
            "power_source": "能源類型"
        }
    ],
    "preference": {
        "cooking_method": "烹調方式",
        "dietary_restrictions": ["限制1", "限制2"],
        "serving_size": "份量"
    }
}
```

**回應格式：**
```json
{
    "suggested_recipes": [
        {
            "dish_name": "菜名",
            "dish_description": "菜餚描述",
            "ingredients": [
                {
                    "name": "食材名稱",
                    "type": "食材分類",
                    "amount": "數量",
                    "unit": "單位",
                    "preparation": "處理方式"
                }
            ],
            "equipment": [
                {
                    "name": "設備名稱",
                    "type": "設備分類",
                    "size": "大小",
                    "material": "材質",
                    "power_source": "能源類型"
                }
            ],
            "recipe": [
                {
                    "step_number": 1,
                    "title": "步驟標題",
                    "description": "步驟說明",
                    "actions": [
                        {
                            "action": "動作名稱",
                            "tool_required": "使用工具",
                            "material_required": ["使用材料"],
                            "time_minutes": 整數分鐘,
                            "instruction_detail": "詳細操作方法"
                        }
                    ],
                    "estimated_total_time": "預估時間",
                    "temperature": "溫度或火力",
                    "warnings": ["注意事項"],
                    "notes": "補充備註"
                }
            ]
        }
    ]
}
```

## 注意事項

1. 圖片大小限制為 5MB
2. 支援的圖片格式：jpg、jpeg、png
3. 所有 API 回應都使用繁體中文
4. 建議使用高品質、清晰的圖片以獲得更好的辨識結果
5. 圖片內容應該清晰可見，避免模糊或過暗

## 錯誤處理

API 可能返回以下錯誤：

- 400 Bad Request：請求格式無效
  - 無效的圖片格式
  - 圖片大小超過限制
  - 不支援的檔案類型
- 500 Internal Server Error：伺服器內部錯誤
  - 食物辨識失敗
  - 食材設備辨識失敗
  - 食譜生成失敗
  - 食譜建議生成失敗

## 開發環境設定

1. 複製專案
```bash
git clone [repository-url]
cd recipe-generator
```

2. 設定環境變數
```bash
cp .env.example .env
# 編輯 .env 檔案，填入必要的設定
```

3. 使用 Docker 建置和運行
```bash
docker-compose up --build
```

## 技術堆疊

- Go
- Gin Web Framework
- OpenRouter AI API
- Docker

## 授權

MIT License
