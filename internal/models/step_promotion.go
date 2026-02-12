package models

// StepPromotionResponse represents a single employee eligible for step promotion.
type StepPromotionResponse struct {
	EmployeeID       uint    `json:"employee_id" example:"1"`
	EmployeeName     string  `json:"employee_name" example:"Anna Müller"`
	CurrentStep      int     `json:"current_step" example:"2"`
	EligibleStep     int     `json:"eligible_step" example:"3"`
	YearsOfService   float64 `json:"years_of_service" example:"3.5"`
	ServiceStart     string  `json:"service_start" example:"2022-01-01"`
	Grade            string  `json:"grade" example:"S8a"`
	CurrentAmount    int     `json:"current_amount" example:"329947"`
	NewAmount        int     `json:"new_amount" example:"350089"`
	MonthlyCostDelta int     `json:"monthly_cost_delta" example:"20142"`
	PayPlanID        uint    `json:"payplan_id" example:"1"`
	PayPlanName      string  `json:"payplan_name" example:"TVöD-SuE 2024"`
}

// StepPromotionsResponse is the response for the step promotions endpoint.
type StepPromotionsResponse struct {
	Date                  string                  `json:"date" example:"2025-06-15"`
	TotalMonthlyCostDelta int                     `json:"total_monthly_cost_delta" example:"40284"`
	Promotions            []StepPromotionResponse `json:"promotions"`
}
