package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"recipe-generator/internal/models"
	"recipe-generator/internal/services/ai"

	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	aiService ai.AIService
}

func NewRecipeHandler(aiService ai.AIService) *RecipeHandler {
	return &RecipeHandler{
		aiService: aiService,
	}
}

// RecognizeFood 處理食物圖片辨識請求
func (h *RecipeHandler) RecognizeFood(c *gin.Context) {
	var request models.FoodRecognitionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	// 檢查圖片格式
	if !strings.HasPrefix(request.Image, "data:image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的圖片格式"})
		return
	}

	response, err := h.aiService.RecognizeFood(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("食物辨識失敗: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RecognizeIngredientsAndEquipment 處理食材和設備圖片辨識請求
func (h *RecipeHandler) RecognizeIngredientsAndEquipment(c *gin.Context) {
	var request models.IngredientEquipmentRecognitionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	// 檢查圖片格式
	if !strings.HasPrefix(request.Image, "data:image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的圖片格式"})
		return
	}

	response, err := h.aiService.RecognizeIngredientsAndEquipment(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("食材設備辨識失敗: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GenerateRecipeFromDishName 處理根據食物名稱生成食譜的請求
func (h *RecipeHandler) GenerateRecipeFromDishName(c *gin.Context) {
	var request models.DishNameRecipeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	response, err := h.aiService.GenerateRecipeFromDishName(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("食譜生成失敗: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GenerateRecipesFromIngredients 處理根據食材和設備生成食譜建議的請求
func (h *RecipeHandler) GenerateRecipesFromIngredients(c *gin.Context) {
	var request models.IngredientsRecipeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	response, err := h.aiService.GenerateRecipesFromIngredients(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("食譜建議生成失敗: %v", err)})
		return
	}

	c.JSON(http.StatusOK, response)
}
