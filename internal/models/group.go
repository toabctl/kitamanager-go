package models

import (
	"time"
)

type Group struct {
	ID            uint           `gorm:"primaryKey" json:"id" example:"1"`
	Name          string         `gorm:"size:255;not null" json:"name" binding:"required" example:"Administrators"`
	Active        bool           `gorm:"default:true" json:"active" example:"true"`
	CreatedAt     time.Time      `json:"created_at" example:"2024-01-15T10:30:00Z"`
	CreatedBy     string         `gorm:"size:255" json:"created_by" example:"admin@example.com"`
	UpdatedAt     time.Time      `json:"updated_at" example:"2024-01-15T10:30:00Z"`
	Users         []User         `gorm:"many2many:user_groups;" json:"users,omitempty"`
	Organizations []Organization `gorm:"many2many:group_organizations;" json:"organizations,omitempty"`
}
