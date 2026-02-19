package models

import "time"

// ============================================================
// GORM models (stored in database)
// ============================================================

// GovernmentFundingBillPayment represents a single financial line item for a child in a bill.
type GovernmentFundingBillPayment struct {
	ID      uint   `gorm:"primaryKey" json:"-"`
	ChildID uint   `gorm:"not null;index" json:"-"`
	Key     string `gorm:"size:100;not null" json:"key" example:"care_type"`
	Value   string `gorm:"size:255;not null" json:"value" example:"ganztag"`
	Amount  int    `gorm:"not null" json:"amount" example:"166847"`
}

// GovernmentFundingBillChild represents one child row in a bill period.
type GovernmentFundingBillChild struct {
	ID            uint                           `gorm:"primaryKey" json:"-"`
	PeriodID      uint                           `gorm:"not null;index" json:"-"`
	VoucherNumber string                         `gorm:"size:20;not null" json:"voucher_number" example:"GB-12345678901-02"`
	ChildName     string                         `gorm:"size:255;not null" json:"child_name" example:"Mustermann, Max"`
	BirthDate     string                         `gorm:"size:10;not null" json:"birth_date" example:"01.20"`
	District      int64                          `gorm:"not null" json:"district" example:"1"`
	Payments      []GovernmentFundingBillPayment `gorm:"foreignKey:ChildID;constraint:OnDelete:CASCADE" json:"payments"`
}

// GovernmentFundingBillPeriod represents a single uploaded government funding bill.
type GovernmentFundingBillPeriod struct {
	ID                uint                         `gorm:"primaryKey" json:"id" example:"1"`
	OrganizationID    uint                         `gorm:"not null;index" json:"organization_id" example:"1"`
	Period                                         // from_date, to_date
	FileName          string                       `gorm:"size:255;not null" json:"file_name" example:"Abrechnung_11-25.xlsx"`
	FileSha256        string                       `gorm:"size:64;not null" json:"file_sha256" example:"a1b2c3d4..."`
	FacilityName      string                       `gorm:"size:255;not null" json:"facility_name" example:"Kita Sonnenschein"`
	FacilityTotal     int                          `gorm:"not null" json:"facility_total" example:"500000"`
	ContractBooking   int                          `gorm:"not null" json:"contract_booking" example:"480000"`
	CorrectionBooking int                          `gorm:"not null" json:"correction_booking" example:"20000"`
	CreatedBy         uint                         `gorm:"not null" json:"created_by" example:"1"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
	Children          []GovernmentFundingBillChild `gorm:"foreignKey:PeriodID;constraint:OnDelete:CASCADE" json:"children,omitempty"`
}

// ============================================================
// Response DTOs (enriched at read time, not stored)
// ============================================================

// GovernmentFundingBillAmount represents a single financial line item in a bill response.
type GovernmentFundingBillAmount struct {
	Key    string `json:"key" example:"care_type"`
	Value  string `json:"value" example:"ganztag"`
	Amount int    `json:"amount" example:"166847"`
}

// GovernmentFundingBillChildResponse represents one child from a bill, enriched with match info.
type GovernmentFundingBillChildResponse struct {
	VoucherNumber string                        `json:"voucher_number" example:"GB-12345678901-02"`
	ChildName     string                        `json:"child_name" example:"Mustermann, Max"`
	BirthDate     string                        `json:"birth_date" example:"01.20"`
	District      int64                         `json:"district" example:"1"`
	TotalAmount   int                           `json:"total_amount" example:"166847"`
	Amounts       []GovernmentFundingBillAmount `json:"amounts"`
	ChildID       *uint                         `json:"child_id,omitempty" example:"42"`
	ContractID    *uint                         `json:"contract_id,omitempty" example:"99"`
	Matched       bool                          `json:"matched" example:"true"`
}

// GovernmentFundingBillPeriodResponse is the full detail response for a single bill period.
type GovernmentFundingBillPeriodResponse struct {
	ID                uint                                 `json:"id" example:"1"`
	OrganizationID    uint                                 `json:"organization_id" example:"1"`
	From              string                               `json:"from" example:"2025-11-01"`
	To                string                               `json:"to" example:"2025-11-30"`
	FileName          string                               `json:"file_name" example:"Abrechnung_11-25.xlsx"`
	FileSha256        string                               `json:"file_sha256" example:"a1b2c3d4..."`
	FacilityName      string                               `json:"facility_name" example:"Kita Sonnenschein"`
	FacilityTotal     int                                  `json:"facility_total" example:"500000"`
	ContractBooking   int                                  `json:"contract_booking" example:"480000"`
	CorrectionBooking int                                  `json:"correction_booking" example:"20000"`
	ChildrenCount     int                                  `json:"children_count" example:"25"`
	MatchedCount      int                                  `json:"matched_count" example:"23"`
	UnmatchedCount    int                                  `json:"unmatched_count" example:"2"`
	Surcharges        []GovernmentFundingBillAmount        `json:"surcharges"`
	Children          []GovernmentFundingBillChildResponse `json:"children"`
	CreatedBy         uint                                 `json:"created_by" example:"1"`
	CreatedAt         time.Time                            `json:"created_at"`
}

// GovernmentFundingBillPeriodListResponse is the summary response for list view.
type GovernmentFundingBillPeriodListResponse struct {
	ID                uint      `json:"id" example:"1"`
	From              string    `json:"from" example:"2025-11-01"`
	To                string    `json:"to" example:"2025-11-30"`
	FileName          string    `json:"file_name" example:"Abrechnung_11-25.xlsx"`
	FacilityName      string    `json:"facility_name" example:"Kita Sonnenschein"`
	FacilityTotal     int       `json:"facility_total" example:"500000"`
	ContractBooking   int       `json:"contract_booking" example:"480000"`
	CorrectionBooking int       `json:"correction_booking" example:"20000"`
	ChildrenCount     int       `json:"children_count" example:"25"`
	CreatedAt         time.Time `json:"created_at"`
}

// GovernmentFundingBillResponse is the full response for the ISBJ upload endpoint (backwards compatible).
type GovernmentFundingBillResponse struct {
	ID                uint                                 `json:"id" example:"1"`
	FacilityName      string                               `json:"facility_name" example:"Kita Sonnenschein"`
	FacilityTotal     int                                  `json:"facility_total" example:"500000"`
	ContractBooking   int                                  `json:"contract_booking" example:"480000"`
	CorrectionBooking int                                  `json:"correction_booking" example:"20000"`
	ChildrenCount     int                                  `json:"children_count" example:"25"`
	MatchedCount      int                                  `json:"matched_count" example:"23"`
	UnmatchedCount    int                                  `json:"unmatched_count" example:"2"`
	Surcharges        []GovernmentFundingBillAmount        `json:"surcharges"`
	Children          []GovernmentFundingBillChildResponse `json:"children"`
}
