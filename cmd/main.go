package main

import (
	"log"

	"recipe-generator/internal/config"
	"recipe-generator/internal/handlers"
	"recipe-generator/internal/services/ai"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// 載入配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("無法載入設定: %v", err)
	}

	// 初始化 AI 服務
	aiService, err := ai.NewOpenRouterService(cfg)
	if err != nil {
		log.Fatalf("無法初始化 AI 服務: %v", err)
	}

	// 初始化 handlers
	recipeHandler := handlers.NewRecipeHandler(aiService)

	// 設定路由
	r := gin.Default()

	// 設定 CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API 路由
	api := r.Group("/api")
	{
		// 圖片辨識相關
		api.POST("/recognize-food", recipeHandler.RecognizeFood)
		api.POST("/recognize-ingredients", recipeHandler.RecognizeIngredientsAndEquipment)

		// 食譜生成相關
		api.POST("/generate-recipe", recipeHandler.GenerateRecipeFromDishName)
		api.POST("/suggest-recipes", recipeHandler.GenerateRecipesFromIngredients)
	}

	// 啟動伺服器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("無法啟動伺服器: %v", err)
	}
}
