package common

import (
	"fmt"
	"strings"
)

// Ingredient 食材
type Ingredient struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	Unit        string `json:"unit"`
	Preparation string `json:"preparation"`
}

// Equipment 設備
type Equipment struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        string `json:"size,omitempty"`
	Material    string `json:"material,omitempty"`
	PowerSource string `json:"power_source,omitempty"`
}

// IngredientRecognitionResult 食材識別結果
type IngredientRecognitionResult struct {
	Ingredients []Ingredient `json:"ingredients"` // 識別出的食材列表
	Equipment   []Equipment  `json:"equipment"`   // 識別出的設備列表
	Summary     string       `json:"summary"`     // 識別內容摘要
}

// RecipePreferences 食譜偏好
type RecipePreferences struct {
	CookingMethod       string   `json:"cooking_method"`
	DietaryRestrictions []string `json:"dietary_restrictions"`
	ServingSize         string   `json:"serving_size"`
}

// RecipeByNameResponse 完全符合 recipe-api.yaml
// Recipe 食譜
// 注意：欄位名稱、型別、巢狀結構、陣列都要一模一樣

type Recipe struct {
	DishName        string       `json:"dish_name"`
	DishDescription string       `json:"dish_description"`
	Ingredients     []Ingredient `json:"ingredients"`
	Equipment       []Equipment  `json:"equipment"`
	Recipe          []RecipeStep `json:"recipe"`
}

type RecipeStep struct {
	StepNumber         int            `json:"step_number"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Actions            []RecipeAction `json:"actions"`
	EstimatedTotalTime string         `json:"estimated_total_time"`
	Temperature        string         `json:"temperature"`
	Warnings           string         `json:"warnings"`
	Notes              string         `json:"notes"`
}

type RecipeAction struct {
	Action            string   `json:"action"`
	ToolRequired      string   `json:"tool_required"`
	MaterialRequired  []string `json:"material_required"`
	TimeMinutes       int      `json:"time_minutes"`
	InstructionDetail string   `json:"instruction_detail"`
}

// RecipeByIngredientsRequest 根據食材推薦食譜的請求
type RecipeByIngredientsRequest struct {
	AvailableIngredients []Ingredient `json:"available_ingredients"`
	AvailableEquipment   []Equipment  `json:"available_equipment"`
	Preference           struct {
		CookingMethod       string   `json:"cooking_method"`
		DietaryRestrictions []string `json:"dietary_restrictions"`
		ServingSize         string   `json:"serving_size"`
	} `json:"preference"`
}

// RecipeByIngredientsResponse 根據食材推薦食譜的回應
type RecipeByIngredientsResponse struct {
	SuggestedRecipes []struct {
		DishName        string       `json:"dish_name"`
		DishDescription string       `json:"dish_description"`
		Ingredients     []Ingredient `json:"ingredients"`
		Equipment       []Equipment  `json:"equipment"`
		Recipe          []struct {
			StepNumber  int    `json:"step_number"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Actions     []struct {
				Action            string   `json:"action"`
				ToolRequired      string   `json:"tool_required"`
				MaterialRequired  []string `json:"material_required"`
				TimeMinutes       int      `json:"time_minutes"`
				InstructionDetail string   `json:"instruction_detail"`
			} `json:"actions"`
			EstimatedTotalTime string  `json:"estimated_total_time"`
			Temperature        string  `json:"temperature"`
			Warnings           *string `json:"warnings"`
			Notes              string  `json:"notes"`
		} `json:"recipe"`
	} `json:"suggested_recipes"`
}

// FormatIngredients 格式化食材列表
func FormatIngredients(ingredients []Ingredient) string {
	var sb strings.Builder
	for _, ing := range ingredients {
		sb.WriteString(fmt.Sprintf("- %s (%s): %s%s, %s\n",
			ing.Name, ing.Type, ing.Amount, ing.Unit, ing.Preparation))
	}
	return sb.String()
}

// FormatEquipment 格式化設備列表
func FormatEquipment(equipment []Equipment) string {
	var sb strings.Builder
	for _, equip := range equipment {
		sb.WriteString(fmt.Sprintf("- %s (%s): %s, %s, %s\n",
			equip.Name, equip.Type, equip.Size, equip.Material, equip.PowerSource))
	}
	return sb.String()
}

// FoodRecognitionResult 食物辨識結果
type FoodRecognitionResult struct {
	RecognizedFoods []RecognizedFood `json:"recognized_foods"`
}

// RecognizedFood 辨識到的食物
type RecognizedFood struct {
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	PossibleIngredients []PossibleIngredient `json:"possible_ingredients"`
	PossibleEquipment   []PossibleEquipment  `json:"possible_equipment"`
}

// PossibleIngredient 可能的食材
type PossibleIngredient struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// PossibleEquipment 可能的設備
type PossibleEquipment struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
