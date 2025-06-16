package models

import (
	"time"
)

// UserBiometric represents the user_biometrics table in the database
type UserBiometric struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	Type      string    `gorm:"column:type;not null" json:"type"`
	Value     float64   `gorm:"column:value;not null" json:"value"`
	Unit      string    `gorm:"column:unit;not null" json:"unit"`
}

// TableName specifies the table name for the UserBiometric model
func (UserBiometric) TableName() string {
	return "user_biometrics"
}

// BiometricTypes contains constants for different biometric types
type BiometricTypes struct {
	Weight                 string
	Height                 string
	BodyFatPercentage      string
	MuscleMass             string
	BMI                    string
	WaistCircumference     string
	HipCircumference       string
	ChestCircumference     string
	ArmCircumference       string
	ThighCircumference     string
	BloodPressureSystolic  string
	BloodPressureDiastolic string
	RestingHeartRate       string
	BodyWaterPercentage    string
	BoneDensity            string
}

// GetBiometricTypes returns the available biometric types
func GetBiometricTypes() BiometricTypes {
	return BiometricTypes{
		Weight:                 "weight",
		Height:                 "height",
		BodyFatPercentage:      "body_fat_percentage",
		MuscleMass:             "muscle_mass",
		BMI:                    "bmi",
		WaistCircumference:     "waist_circumference",
		HipCircumference:       "hip_circumference",
		ChestCircumference:     "chest_circumference",
		ArmCircumference:       "arm_circumference",
		ThighCircumference:     "thigh_circumference",
		BloodPressureSystolic:  "blood_pressure_systolic",
		BloodPressureDiastolic: "blood_pressure_diastolic",
		RestingHeartRate:       "resting_heart_rate",
		BodyWaterPercentage:    "body_water_percentage",
		BoneDensity:            "bone_density",
	}
}

// ProgressData represents progress data for visualization
type ProgressData struct {
	Date   time.Time `json:"date"`
	Value  float64   `json:"value"`
	Change float64   `json:"change"`
	Trend  string    `json:"trend"` // "up", "down", "stable"
}

// BiometricProgress represents overall progress for a biometric type
type BiometricProgress struct {
	Type          string         `json:"type"`
	Unit          string         `json:"unit"`
	CurrentValue  float64        `json:"current_value"`
	PreviousValue float64        `json:"previous_value"`
	OverallChange float64        `json:"overall_change"`
	PercentChange float64        `json:"percent_change"`
	Trend         string         `json:"trend"`
	DataPoints    []ProgressData `json:"data_points"`
	StartDate     time.Time      `json:"start_date"`
	EndDate       time.Time      `json:"end_date"`
}

// ChartData represents data formatted for charts
type ChartData struct {
	Type      string    `json:"type"`
	Unit      string    `json:"unit"`
	Labels    []string  `json:"labels"`
	Values    []float64 `json:"values"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Goal represents a biometric goal
type Goal struct {
	Type         string    `json:"type"`
	TargetValue  float64   `json:"target_value"`
	CurrentValue float64   `json:"current_value"`
	Unit         string    `json:"unit"`
	TargetDate   time.Time `json:"target_date"`
	Progress     float64   `json:"progress"`
	IsAchieved   bool      `json:"is_achieved"`
}

// GoalProgress represents progress towards goals
type GoalProgress struct {
	Goals           []Goal  `json:"goals"`
	OverallProgress float64 `json:"overall_progress"`
	AchievedGoals   int     `json:"achieved_goals"`
	TotalGoals      int     `json:"total_goals"`
}

// BiometricSummary represents a summary of all biometrics for a user
type BiometricSummary struct {
	UserID           uint                         `json:"user_id"`
	LatestBiometrics map[string]UserBiometric     `json:"latest_biometrics"`
	ProgressData     map[string]BiometricProgress `json:"progress_data"`
	Goals            GoalProgress                 `json:"goals"`
	LastUpdated      time.Time                    `json:"last_updated"`
}

// AdvancedMetrics represents calculated advanced metrics
type AdvancedMetrics struct {
	BMI                 float64 `json:"bmi"`
	BodyFatPercentage   float64 `json:"body_fat_percentage"`
	MuscleMass          float64 `json:"muscle_mass"`
	WaistToHipRatio     float64 `json:"waist_to_hip_ratio"`
	BodyWaterPercentage float64 `json:"body_water_percentage"`
	BMICategory         string  `json:"bmi_category"`
	BodyFatCategory     string  `json:"body_fat_category"`
	HealthRisk          string  `json:"health_risk"`
}
