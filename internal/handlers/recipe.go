package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "recipe-generator/internal/models"
    "recipe-generator/internal/services"
)

type RecipeHandler struct {
    aiService *services.AIService
}

func NewRecipeHandler(aiService *services.AIService) *RecipeHandler {
    return &RecipeHandler{
        aiService: aiService,
    }
}

func (h *RecipeHandler) GenerateRecipe(c *gin.Context) {
    var request models.RecipeRequest
    
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "無效的請求格式",
        })
        return
    }

    recipe, err := h.aiService.GenerateRecipe(&request)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "食譜生成失敗",
        })
        return
    }

    c.JSON(http.StatusOK, recipe)
} 