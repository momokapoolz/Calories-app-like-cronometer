package dto

type NutritionSummaryDTO struct {
	UserID                 uint    `json:"user_id"`
	DateRange              string  `json:"date_range"`
	TotalCalories          float64 `json:"total_calories"`
	MacroNutrientBreakDown []MacronutrientBreakdownDTO
	MicroNutrientBreakDown []MicronutrientDTO
	MealBreakdown          []MealNutritionDTO
}

type MacronutrientBreakdownDTO struct {
	Energy       float64 `json:"energy"`
	Protein      float64 `json:"protein"`
	Fat          float64 `json:"total_lipid_fe"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fiber        float64 `json:"fiber"`
	Cholesterol  float64 `json:"cholesteroid"`
	Vitamin_A    float64 `json:"vitamin_a"`
	Vitamin_B12  float64 `json:"vitamin_b"`
	Calcium      float64 `json:"calcium"`
	Iron         float64 `json:"iron"`
}

type MicronutrientDTO struct {
	NutrientID   uint    `json:"nutrient_id"`
	NutrientName string  `json:"nutrient_name"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
}

type MealNutritionDTO struct {
	MealLogID    uint    `json:"meal_log_id"`
	MealType     string  `json:"meal_type"`
	Date         string  `json:"date"`
	Calories     float64 `json:"calories"`
	Protein      float64 `json:"protein"`
	Carbohydrate float64 `json:"carbohydrate"`
	Fat          float64 `json:"fat"`
	FoodCount    int     `json:"food_count"`
}

// MealNutritionDetailDTO represents detailed nutrition information for a single meal
type MealNutritionDetailDTO struct {
	MealLogID              uint                        `json:"meal_log_id"`
	UserID                 uint                        `json:"user_id"`
	MealType               string                      `json:"meal_type"`
	Date                   string                      `json:"date"`
	TotalCalories          float64                     `json:"total_calories"`
	FoodCount              int                         `json:"food_count"`
	MacroNutrientBreakDown []MacronutrientBreakdownDTO `json:"MacroNutrientBreakDown"`
	MicroNutrientBreakDown []MicronutrientDTO          `json:"MicroNutrientBreakDown"`
}
