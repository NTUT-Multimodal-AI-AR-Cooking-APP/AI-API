package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"recipe-generator/internal/models"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewAIService() (*AIService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed to create AI client: %v", err)
	}

	model := client.GenerativeModel("models/gemini-2.0-flash-001")
	return &AIService{
		client: client,
		model:  model,
	}, nil
}

func (s *AIService) GenerateRecipe(request *models.RecipeRequest) (*models.RecipeResponse, error) {
	ctx := context.Background()

	// 構建提示詞
	prompt := buildPrompt(request)

	// 生成回應
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// 取得回傳的 AI 原始內容
	raw, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("AI 回傳格式不是 genai.Text")
	}
	// 將 raw（genai.Text）轉換為 string
	rawString := string(raw)
	// 印出 AI 回傳內容，幫助 debug
	//fmt.Println("🧠 AI 回傳內容：", rawString)

	// 移除 Markdown 格式的反引號
	cleanJson := strings.TrimPrefix(rawString, "```json\n")
	cleanJson = strings.TrimSuffix(cleanJson, "```")

	// 偵測 JSON 開頭與結尾
	startIdx := 0
	endIdx := len(cleanJson)

	// 提取出合法的 JSON 部分
	validJson := cleanJson[startIdx:endIdx]
	//fmt.Println("📄 有效的 JSON 字串：", validJson)

	// 解析 JSON
	var recipe models.RecipeResponse
	if err := json.Unmarshal([]byte(validJson), &recipe); err != nil {
		fmt.Println("🔴 解析 JSON 失敗：", err) // 顯示解析錯誤
		return nil, fmt.Errorf("failed to parse AI response: %v", err)
	}

	// 顯示解析成功的結果
	//fmt.Println("🍽️ 解析後的食譜：", recipe)

	return &recipe, nil
}

func buildPrompt(request *models.RecipeRequest) string {
	return fmt.Sprintf(`請根據以下資訊生成一道食譜：
    設備：%v
    食材：%v
    偏好：%v
    
    請以 JSON 格式回應，包含：菜名、菜餚描述和詳細的烹飪步驟。每個步驟應包含：步驟說明、所需時間、溫度、描述和熟度等等（如果適用）。
    
    回應格式必須完全符合以下 JSON 結構：
    {
        "dish_name": "菜名",
        "dish_description": "菜餚描述",
        "recipe": [
            {
                "step": "步驟說明",
                "time": "所需時間",
                "temperature": "溫度",
                "description": "描述",
                "doneness": "熟度（如果適用）"
            }
        ]
    }`, request.Equipment, request.Ingredients, request.Preference)
}
