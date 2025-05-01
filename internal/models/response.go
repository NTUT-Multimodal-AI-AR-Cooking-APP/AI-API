package models

type RecipeStep struct {
    Step        string `json:"step"`
    Time        string `json:"time"`
    Temperature string `json:"temperature"`
    Description string `json:"description"`
    Doneness    string `json:"doneness,omitempty"`
}

type RecipeResponse struct {
    DishName        string       `json:"dish_name"`
    DishDescription string       `json:"dish_description"`
    Recipe          []RecipeStep `json:"recipe"`
} 