package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"recipe-generator/internal/config"
	"recipe-generator/internal/models"

	"github.com/revrost/go-openrouter"
)

// AIService 定義 AI 服務的介面
type AIService interface {
	RecognizeFood(ctx context.Context, request *models.FoodRecognitionRequest) (*models.FoodRecognitionResponse, error)
	RecognizeIngredientsAndEquipment(ctx context.Context, request *models.IngredientEquipmentRecognitionRequest) (*models.IngredientEquipmentRecognitionResponse, error)
	GenerateRecipeFromDishName(ctx context.Context, request *models.DishNameRecipeRequest) (*models.RecipeResponse, error)
	GenerateRecipesFromIngredients(ctx context.Context, request *models.IngredientsRecipeRequest) (*models.SuggestedRecipesResponse, error)
}

type OpenRouterService struct {
	client *openrouter.Client
	model  string
	// 添加默認請求參數
	defaultParams openrouter.ChatCompletionRequest
}

func NewOpenRouterService(cfg *config.Config) (*OpenRouterService, error) {
	client := openrouter.NewClient(cfg.OpenRouter.APIKey)

	// 設置默認請求參數
	defaultParams := openrouter.ChatCompletionRequest{
		Model:       cfg.OpenRouter.Model,
		MaxTokens:   1000,  // 限制回應長度
		Temperature: 0.7,   // 平衡創造性和準確性
		Stream:      false, // 不使用串流，因為需要完整的 JSON 回應
	}

	return &OpenRouterService{
		client:        client,
		model:         cfg.OpenRouter.Model,
		defaultParams: defaultParams,
	}, nil
}

// cleanJSON 清理和驗證 JSON 字符串
func cleanJSON(input string) (string, error) {
	// 移除可能的 markdown 代碼塊標記
	clean := strings.TrimPrefix(input, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	// 移除可能的 BOM 和其他特殊字符
	clean = regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`).ReplaceAllString(clean, "")

	// 驗證是否為有效的 JSON
	var test map[string]interface{}
	if err := json.Unmarshal([]byte(clean), &test); err != nil {
		return "", fmt.Errorf("無效的 JSON 格式: %v", err)
	}

	return clean, nil
}

// RecognizeFood 識別圖片中的食物
func (s *OpenRouterService) RecognizeFood(ctx context.Context, request *models.FoodRecognitionRequest) (*models.FoodRecognitionResponse, error) {
	prompt := fmt.Sprintf(`請分析這張食物圖片，並以嚴格的 JSON 格式回應。請確保回應只包含有效的 JSON 數據，不要包含任何其他文字或標記。

要求：
1. 使用繁體中文
2. 只返回 JSON 格式數據
3. 不要包含任何 markdown 標記
4. 不要包含任何解釋或說明文字

JSON 格式：
{
    "recognized_foods": [{
        "name": "食物名稱",
        "description": "特色、口感、適合場合",
        "possible_ingredients": [{"name": "食材名", "type": "分類"}],
        "possible_equipment": [{"name": "設備名", "type": "分類"}]
    }]
}

%s`, request.DescriptionHint)

	req := s.defaultParams
	req.Messages = []openrouter.ChatCompletionMessage{
		{
			Role: "system",
			Content: openrouter.Content{
				Text: "你是一個專業的食物識別助手。你的任務是分析圖片並返回嚴格的 JSON 格式數據。請確保回應只包含有效的 JSON，不要包含任何其他文字。",
			},
		},
		{
			Role: "user",
			Content: openrouter.Content{
				Multi: []openrouter.ChatMessagePart{
					{Type: "text", Text: prompt},
					{Type: "image_url", ImageURL: &openrouter.ChatMessageImageURL{URL: request.Image}},
				},
			},
		},
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("食物辨識失敗: %v", err)
	}

	content := resp.Choices[0].Message.Content.Text
	cleanJson, err := cleanJSON(content)
	if err != nil {
		return nil, fmt.Errorf("食物辨識失敗: %v", err)
	}

	var response models.FoodRecognitionResponse
	if err := json.Unmarshal([]byte(cleanJson), &response); err != nil {
		return nil, fmt.Errorf("解析 AI 回應失敗: %v\n原始回應: %s", err, content)
	}

	// 驗證必要字段
	if len(response.RecognizedFoods) == 0 {
		return nil, fmt.Errorf("AI 回應缺少必要的食物信息")
	}

	return &response, nil
}

// RecognizeIngredientsAndEquipment 識別圖片中的食材和設備
func (s *OpenRouterService) RecognizeIngredientsAndEquipment(ctx context.Context, request *models.IngredientEquipmentRecognitionRequest) (*models.IngredientEquipmentRecognitionResponse, error) {
	prompt := fmt.Sprintf(`請分析這張圖片中的食材和設備，並以嚴格的 JSON 格式回應。請確保回應只包含有效的 JSON 數據，不要包含任何其他文字或標記。

要求：
1. 使用繁體中文
2. 只返回 JSON 格式數據
3. 不要包含任何 markdown 標記
4. 不要包含任何解釋或說明文字

JSON 格式：
{
    "ingredients": [{
        "name": "食材名",
        "type": "分類",
        "amount": "數量",
        "unit": "單位",
        "preparation": "處理方式"
    }],
    "equipment": [{
        "name": "設備名",
        "type": "分類",
        "size": "大小",
        "material": "材質",
        "power_source": "能源類型"
    }],
    "summary": "整體摘要"
}

%s`, request.DescriptionHint)

	req := s.defaultParams
	req.Messages = []openrouter.ChatCompletionMessage{
		{
			Role: "system",
			Content: openrouter.Content{
				Text: "你是一個專業的食材和設備識別助手。你的任務是分析圖片並返回嚴格的 JSON 格式數據。請確保回應只包含有效的 JSON，不要包含任何其他文字。",
			},
		},
		{
			Role: "user",
			Content: openrouter.Content{
				Multi: []openrouter.ChatMessagePart{
					{Type: "text", Text: prompt},
					{Type: "image_url", ImageURL: &openrouter.ChatMessageImageURL{URL: request.Image}},
				},
			},
		},
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("食材設備辨識失敗: %v", err)
	}

	content := resp.Choices[0].Message.Content.Text
	cleanJson, err := cleanJSON(content)
	if err != nil {
		return nil, fmt.Errorf("食材設備辨識失敗: %v", err)
	}

	var response models.IngredientEquipmentRecognitionResponse
	if err := json.Unmarshal([]byte(cleanJson), &response); err != nil {
		// 添加更詳細的錯誤信息
		return nil, fmt.Errorf("解析 AI 回應失敗: %v\n原始回應: %s", err, content)
	}

	// 驗證必要字段
	if len(response.Ingredients) == 0 && len(response.Equipment) == 0 {
		return nil, fmt.Errorf("AI 回應缺少必要的食材或設備信息")
	}

	return &response, nil
}

// GenerateRecipeFromDishName 根據食物名稱生成食譜
func (s *OpenRouterService) GenerateRecipeFromDishName(ctx context.Context, request *models.DishNameRecipeRequest) (*models.RecipeResponse, error) {
	prompt := fmt.Sprintf(`根據以下資訊生成食譜（JSON格式）：
菜名：%s
偏好食材：%v
排除食材：%v
偏好設備：%v
烹飪偏好：%s/%s/%s

回應格式：
{
    "dish_name": "菜名",
    "dish_description": "描述",
    "ingredients": [{"name": "食材名", "type": "分類", "amount": "數量", "unit": "單位", "preparation": "處理方式"}],
    "equipment": [{"name": "設備名", "type": "分類", "size": "大小", "material": "材質", "power_source": "能源類型"}],
    "recipe": [{
        "step_number": 1,
        "title": "步驟標題",
        "description": "說明",
        "actions": [{"action": "動作", "tool_required": "工具", "material_required": ["材料"], "time_minutes": 分鐘數, "instruction_detail": "詳細操作"}],
        "estimated_total_time": "時間",
        "temperature": "溫度",
        "warnings": ["注意事項"],
        "notes": "備註"
    }]
}`,
		request.DishName,
		request.PreferredIngredients,
		request.ExcludedIngredients,
		request.PreferredEquipment,
		request.Preference.CookingMethod,
		request.Preference.Doneness,
		request.Preference.ServingSize)

	req := s.defaultParams
	req.MaxTokens = 2000 // 食譜生成需要更長的回應
	req.Messages = []openrouter.ChatCompletionMessage{
		{Role: "user", Content: openrouter.Content{Text: prompt}},
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("食譜生成失敗: %v", err)
	}

	content := resp.Choices[0].Message.Content.Text
	cleanJson := strings.TrimPrefix(content, "```json\n")
	cleanJson = strings.TrimSuffix(cleanJson, "```")

	var response models.RecipeResponse
	if err := json.Unmarshal([]byte(cleanJson), &response); err != nil {
		return nil, fmt.Errorf("解析 AI 回應失敗: %v", err)
	}

	return &response, nil
}

// GenerateRecipesFromIngredients 根據食材和設備生成食譜建議
func (s *OpenRouterService) GenerateRecipesFromIngredients(ctx context.Context, request *models.IngredientsRecipeRequest) (*models.SuggestedRecipesResponse, error) {
	ingredientsJson, _ := json.Marshal(request.AvailableIngredients)
	equipmentJson, _ := json.Marshal(request.AvailableEquipment)

	prompt := fmt.Sprintf(`根據以下資訊生成多個食譜建議（JSON格式）：
可用食材：%s
可用設備：%s
烹飪偏好：%s/%v/%s

回應格式：
{
    "suggested_recipes": [{
        "dish_name": "菜名",
        "dish_description": "描述",
        "ingredients": [{"name": "食材名", "type": "分類", "amount": "數量", "unit": "單位", "preparation": "處理方式"}],
        "equipment": [{"name": "設備名", "type": "分類", "size": "大小", "material": "材質", "power_source": "能源類型"}],
        "recipe": [{
            "step_number": 1,
            "title": "步驟標題",
            "description": "說明",
            "actions": [{"action": "動作", "tool_required": "工具", "material_required": ["材料"], "time_minutes": 分鐘數, "instruction_detail": "詳細操作"}],
            "estimated_total_time": "時間",
            "temperature": "溫度",
            "warnings": ["注意事項"],
            "notes": "備註"
        }]
    }]
}`,
		string(ingredientsJson),
		string(equipmentJson),
		request.Preference.CookingMethod,
		request.Preference.DietaryRestrictions,
		request.Preference.ServingSize)

	req := s.defaultParams
	req.MaxTokens = 3000 // 多個食譜建議需要更長的回應
	req.Messages = []openrouter.ChatCompletionMessage{
		{Role: "user", Content: openrouter.Content{Text: prompt}},
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("食譜建議生成失敗: %v", err)
	}

	content := resp.Choices[0].Message.Content.Text
	cleanJson := strings.TrimPrefix(content, "```json\n")
	cleanJson = strings.TrimSuffix(cleanJson, "```")

	var response models.SuggestedRecipesResponse
	if err := json.Unmarshal([]byte(cleanJson), &response); err != nil {
		return nil, fmt.Errorf("解析 AI 回應失敗: %v", err)
	}

	return &response, nil
}
