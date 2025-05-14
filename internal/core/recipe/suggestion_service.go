package recipe

import (
	"context"
	"fmt"
	"strings"

	"recipe-generator/internal/core/ai/cache"
	"recipe-generator/internal/core/ai/service"
	"recipe-generator/internal/pkg/common"

	"go.uber.org/zap"
)

// SuggestionService 食譜推薦服務
type SuggestionService struct {
	aiService    *service.Service
	cacheManager *cache.CacheManager
}

// NewSuggestionService 創建新的食譜推薦服務
func NewSuggestionService(aiService *service.Service, cacheManager *cache.CacheManager) *SuggestionService {
	return &SuggestionService{
		aiService:    aiService,
		cacheManager: cacheManager,
	}
}

// SuggestRecipes 根據可用食材和設備推薦食譜
func (s *SuggestionService) SuggestRecipes(ctx context.Context, req *common.RecipeByIngredientsRequest) (*common.Recipe, error) {
	// 驗證必要欄位
	if req.Preference.CookingMethod == "" || req.Preference.ServingSize == "" {
		return nil, fmt.Errorf("missing required fields: cooking_method and serving_size are required")
	}

	prompt := fmt.Sprintf(`請根據以下可用食材和設備，推薦適合的食譜(並且用繁體中文回答）。

可用食材：
%s

可用設備：
%s

烹飪偏好：
- 烹飪方式：%s
- 飲食限制：%s
- 份量：%s

要求：
1. 只根據提供的食材和設備推薦內容，不要添加未出現的食材或設備
2. 不要使用預設值或猜測值，若無法確定請填寫 "未知"
3. 每個步驟都要非常詳細，適合新手操作
4. 動作描述要具體明確，包含具體的時間和溫度
5. 注意事項要特別提醒新手容易忽略的細節
6. 所有字段都必須使用雙引號
7. 不需要考慮可讀性，請省略所有空格和換行，返回最緊湊的 JSON 格式
8. 推薦的食譜要優先使用已有的食材和設備
9. 如果某些食材或設備不足，可以建議替代方案
10. 每個食譜都要考慮到烹飪難度和時間
11. time_minutes 欄位必須是整數，不能有小數點（以秒為單位）
12. warnings 欄位必須是字串類型，如果沒有警告事項請填寫 "null"
13. 每個步驟都必須包含 warnings 欄位，不能省略此欄位
14. 不要使用\n，不需要換行

請以以下 JSON 格式返回（僅作為範例，請勿直接複製內容）：
{
    "dish_name": "菜名",
    "dish_description": "描述",
    "ingredients": [
        {
            "name": "食材名稱",
            "type": "食材類型",
            "amount": "數量",
            "unit": "單位",
            "preparation": "處理方式"
        }
    ],
    "equipment": [
        {
            "name": "設備名稱",
            "type": "設備類型",
            "size": "尺寸",
            "material": "材質",
            "power_source": "能源類型"
        }
    ],
    "recipe": [
        {
            "step_number": 1,
            "title": "步驟標題",
            "description": "步驟描述",
            "actions": [
                {
                    "action": "動作",
                    "tool_required": "工具",
                    "material_required": ["材料"],
                    "time_minutes": 1,
                    "instruction_detail": "細節"
                }
            ],
            "estimated_total_time": "時間",
            "temperature": "火侯",
            "warnings": "無",
            "notes": "備註"
        }
    ]
}`,
		common.FormatIngredients(req.AvailableIngredients),
		common.FormatEquipment(req.AvailableEquipment),
		req.Preference.CookingMethod,
		strings.Join(req.Preference.DietaryRestrictions, "、"),
		req.Preference.ServingSize)

	common.LogDebug("SuggestRecipes 組裝的 prompt", zap.String("prompt", prompt))

	resp, err := s.aiService.ProcessRequest(ctx, prompt, "")
	if err != nil {
		return nil, fmt.Errorf("AI service error: %w", err)
	}

	if resp == nil || resp.Content == "" {
		return nil, fmt.Errorf("empty AI response")
	}

	content := resp.Content
	content = strings.TrimSpace(content)
	// 強化 markdown 去除：直接抓第一個 { 到最後一個 } 之間的內容
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start != -1 && end != -1 && end > start {
		content = content[start : end+1]
	}

	var result common.Recipe
	if err := common.ParseJSON(content, &result); err != nil {
		// 只 log 前 100 字與長度，不印全部內容
		aiRespPreview := content
		common.LogError("AI 回應解析失敗",
			zap.Error(err),
			zap.Int("ai_response_length", len(content)),
			zap.String("ai_response_preview", aiRespPreview),
		)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	// 確保每個步驟都有 warnings 欄位
	for i := range result.Recipe {
		if result.Recipe[i].Warnings == "" {
			result.Recipe[i].Warnings = "null"
		}
	}

	return &result, nil
}
