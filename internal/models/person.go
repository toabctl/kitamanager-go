package models

import "time"

// Person represents the base person data shared by Employee and Child.
type Person struct {
	ID             uint          `gorm:"primaryKey" json:"id" example:"1"`
	OrganizationID uint          `gorm:"not null;index" json:"organization_id" example:"1"`
	Organization   *Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization,omitempty"`
	FirstName      string        `gorm:"size:255;not null" json:"first_name" example:"Max"`
	LastName       string        `gorm:"size:255;not null" json:"last_name" example:"Mustermann"`
	Birthdate      time.Time     `gorm:"type:date;not null" json:"birthdate" example:"1990-05-15"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// FullName returns the person's full name.
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}
