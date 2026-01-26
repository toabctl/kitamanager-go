package models

import "time"

// Payplan represents a top-level payment plan definition.
// Organizations are assigned to a payplan to determine funding calculations.
type Payplan struct {
	ID        uint            `gorm:"primaryKey" json:"id" example:"1"`
	Name      string          `gorm:"size:255;not null;uniqueIndex" json:"name" example:"Berlin"`
	CreatedAt time.Time       `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time       `json:"updated_at" example:"2024-01-15T10:30:00Z"`
	Periods   []PayplanPeriod `gorm:"foreignKey:PayplanID;constraint:OnDelete:CASCADE" json:"periods,omitempty"`
}

// PayplanPeriod represents a time period within a payplan.
// Each period has its own set of age-based entries with payment amounts.
// Periods within the same payplan must not overlap - this is enforced at the service layer.
// A period with nil To date is considered ongoing (extends indefinitely into the future).
type PayplanPeriod struct {
	ID        uint           `gorm:"primaryKey" json:"id" example:"1"`
	PayplanID uint           `gorm:"not null;index" json:"payplan_id" example:"1"`
	From      time.Time      `gorm:"column:from_date;type:date;not null" json:"from" example:"2023-03-01"`
	To        *time.Time     `gorm:"column:to_date;type:date" json:"to" example:"2024-02-29"`
	Comment   string         `gorm:"size:1000" json:"comment,omitempty" example:"Funding period 2023/2024"`
	CreatedAt time.Time      `json:"created_at" example:"2024-01-15T10:30:00Z"`
	Entries   []PayplanEntry `gorm:"foreignKey:PeriodID;constraint:OnDelete:CASCADE" json:"entries,omitempty"`
}

// PayplanEntry represents an age range entry within a period.
// MinAge is inclusive, MaxAge is exclusive (e.g., MinAge=0, MaxAge=2 covers ages 0 and 1,
// meaning children from birth up to but not including their 2nd birthday).
type PayplanEntry struct {
	ID       uint `gorm:"primaryKey" json:"id" example:"1"`
	PeriodID uint `gorm:"not null;index" json:"period_id" example:"1"`
	// MinAge is the minimum age in years (inclusive). A child whose age >= MinAge qualifies.
	MinAge int `gorm:"not null" json:"min_age" example:"0"`
	// MaxAge is the maximum age in years (exclusive). A child whose age < MaxAge qualifies.
	MaxAge     int               `gorm:"not null" json:"max_age" example:"2"`
	CreatedAt  time.Time         `json:"created_at" example:"2024-01-15T10:30:00Z"`
	Properties []PayplanProperty `gorm:"foreignKey:EntryID;constraint:OnDelete:CASCADE" json:"properties,omitempty"`
}

// PayplanProperty represents a property value with payment and staffing requirement.
// Payment is stored in cents to avoid floating-point issues (e.g., 166847 = 1668.47 EUR).
type PayplanProperty struct {
	ID          uint      `gorm:"primaryKey" json:"id" example:"1"`
	EntryID     uint      `gorm:"not null;index" json:"entry_id" example:"1"`
	Name        string    `gorm:"size:255;not null" json:"name" example:"ganztag"`
	Payment     int       `gorm:"not null" json:"payment" example:"166847"`
	Requirement float64   `gorm:"not null" json:"requirement" example:"0.261"`
	Comment     string    `gorm:"size:500" json:"comment,omitempty" example:"Full-day care funding"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
}

// PayplanCreate represents the request body for creating a payplan.
type PayplanCreate struct {
	Name string `json:"name" binding:"required,max=255" example:"Berlin"`
}

// PayplanUpdate represents the request body for updating a payplan.
type PayplanUpdate struct {
	Name *string `json:"name" binding:"omitempty,max=255" example:"Berlin Updated"`
}

// PayplanPeriodCreate represents the request body for creating a payplan period.
type PayplanPeriodCreate struct {
	From    time.Time  `json:"from" binding:"required" example:"2023-03-01"`
	To      *time.Time `json:"to" example:"2024-02-29"`
	Comment string     `json:"comment" binding:"max=1000" example:"Funding period 2023/2024"`
}

// PayplanPeriodUpdate represents the request body for updating a payplan period.
type PayplanPeriodUpdate struct {
	From    *time.Time `json:"from" example:"2023-03-01"`
	To      *time.Time `json:"to" example:"2024-02-29"`
	Comment *string    `json:"comment" binding:"omitempty,max=1000" example:"Updated comment"`
}

// PayplanEntryCreate represents the request body for creating a payplan entry.
// MinAge is inclusive, MaxAge is exclusive (e.g., MinAge=0, MaxAge=2 covers ages 0 and 1).
type PayplanEntryCreate struct {
	MinAge int `json:"min_age" binding:"gte=0" example:"0"`
	MaxAge int `json:"max_age" binding:"required,gtfield=MinAge" example:"2"`
}

// PayplanEntryUpdate represents the request body for updating a payplan entry.
// MinAge is inclusive, MaxAge is exclusive (e.g., MinAge=0, MaxAge=2 covers ages 0 and 1).
type PayplanEntryUpdate struct {
	MinAge *int `json:"min_age" binding:"omitempty,gte=0" example:"0"`
	MaxAge *int `json:"max_age" example:"3"`
}

// PayplanPropertyCreate represents the request body for creating a payplan property.
type PayplanPropertyCreate struct {
	Name        string  `json:"name" binding:"required,max=255" example:"ganztag"`
	Payment     int     `json:"payment" binding:"gte=0" example:"166847"`
	Requirement float64 `json:"requirement" binding:"gte=0" example:"0.261"`
	Comment     string  `json:"comment" binding:"max=500" example:"Full-day care funding"`
}

// PayplanPropertyUpdate represents the request body for updating a payplan property.
type PayplanPropertyUpdate struct {
	Name        *string  `json:"name" binding:"omitempty,max=255" example:"ganztag"`
	Payment     *int     `json:"payment" binding:"omitempty,gte=0" example:"166847"`
	Requirement *float64 `json:"requirement" binding:"omitempty,gte=0" example:"0.261"`
	Comment     *string  `json:"comment" binding:"omitempty,max=500" example:"Updated comment"`
}

// AssignPayplanRequest represents the request body for assigning a payplan to an organization.
type AssignPayplanRequest struct {
	PayplanID uint `json:"payplan_id" binding:"required" example:"1"`
}
