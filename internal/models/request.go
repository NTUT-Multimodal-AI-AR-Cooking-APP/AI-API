package models

type Equipment struct {
    Name     string `json:"name"`
    Type     string `json:"type"`
    Size     string `json:"size"`
    Material string `json:"material"`
}

type Ingredient struct {
    Name   string `json:"name"`
    Type   string `json:"type"`
    Amount string `json:"amount"`
    Unit   string `json:"unit"`
    Weight string `json:"weight,omitempty"`
}

type Preference struct {
    CookingMethod string `json:"cooking_method"`
    Doneness      string `json:"doneness"`
}

type RecipeRequest struct {
    Equipment   []Equipment   `json:"equipment"`
    Ingredients []Ingredient  `json:"ingredients"`
    Preference  Preference    `json:"preference"`
} 