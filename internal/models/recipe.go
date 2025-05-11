package models

// 圖片辨識食物請求
type FoodRecognitionRequest struct {
	Image           string `json:"image" binding:"required"`
	DescriptionHint string `json:"description_hint,omitempty"`
}

// 圖片辨識食物回應
type FoodRecognitionResponse struct {
	RecognizedFoods []RecognizedFood `json:"recognized_foods"`
}

type RecognizedFood struct {
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	PossibleIngredients []Ingredient `json:"possible_ingredients"`
	PossibleEquipment   []Equipment  `json:"possible_equipment"`
}

// 圖片辨識食材與設備請求
type IngredientEquipmentRecognitionRequest struct {
	Image           string `json:"image" binding:"required"`
	DescriptionHint string `json:"description_hint,omitempty"`
}

// 圖片辨識食材與設備回應
type IngredientEquipmentRecognitionResponse struct {
	Ingredients []Ingredient `json:"ingredients"`
	Equipment   []Equipment  `json:"equipment"`
	Summary     string       `json:"summary"`
}

// 食材定義
type Ingredient struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Amount      string `json:"amount,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Weight      string `json:"weight,omitempty"`
	Preparation string `json:"preparation,omitempty"`
}

// 設備定義
type Equipment struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        string `json:"size,omitempty"`
	Material    string `json:"material,omitempty"`
	PowerSource string `json:"power_source,omitempty"`
}

// 烹飪偏好設定
type CookingPreference struct {
	CookingMethod       string   `json:"cooking_method,omitempty"`
	Doneness            string   `json:"doneness,omitempty"`
	ServingSize         string   `json:"serving_size,omitempty"`
	DietaryRestrictions []string `json:"dietary_restrictions,omitempty"`
}

// 食物名稱生成食譜請求
type DishNameRecipeRequest struct {
	DishName             string            `json:"dish_name" binding:"required"`
	PreferredIngredients []string          `json:"preferred_ingredients,omitempty"`
	ExcludedIngredients  []string          `json:"excluded_ingredients,omitempty"`
	PreferredEquipment   []string          `json:"preferred_equipment,omitempty"`
	Preference           CookingPreference `json:"preference"`
}

// 食材設備生成食譜請求
type IngredientsRecipeRequest struct {
	AvailableIngredients []Ingredient      `json:"available_ingredients" binding:"required"`
	AvailableEquipment   []Equipment       `json:"available_equipment" binding:"required"`
	Preference           CookingPreference `json:"preference"`
}

// 食譜步驟動作
type RecipeAction struct {
	Action            string   `json:"action"`
	ToolRequired      string   `json:"tool_required"`
	MaterialRequired  []string `json:"material_required"`
	TimeMinutes       int      `json:"time_minutes"`
	InstructionDetail string   `json:"instruction_detail"`
}

// 食譜步驟
type RecipeStep struct {
	StepNumber         int            `json:"step_number"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Actions            []RecipeAction `json:"actions"`
	EstimatedTotalTime string         `json:"estimated_total_time"`
	Temperature        string         `json:"temperature,omitempty"`
	Warnings           []string       `json:"warnings,omitempty"`
	Notes              string         `json:"notes,omitempty"`
}

// 單一食譜回應
type RecipeResponse struct {
	DishName        string       `json:"dish_name"`
	DishDescription string       `json:"dish_description"`
	Ingredients     []Ingredient `json:"ingredients"`
	Equipment       []Equipment  `json:"equipment"`
	Recipe          []RecipeStep `json:"recipe"`
}

// 多個食譜建議回應
type SuggestedRecipesResponse struct {
	SuggestedRecipes []RecipeResponse `json:"suggested_recipes"`
}
