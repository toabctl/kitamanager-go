package models

import (
	"testing"
	"time"
)

func TestEmployee_ToResponse(t *testing.T) {
	now := time.Now()

	t.Run("employee without contracts", func(t *testing.T) {
		emp := Employee{
			Person: Person{
				ID:             1,
				OrganizationID: 2,
				FirstName:      "Max",
				LastName:       "Mustermann",
				Gender:         "male",
				Birthdate:      time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
		}

		resp := emp.ToResponse()

		if resp.ID != 1 {
			t.Errorf("ID = %d, want 1", resp.ID)
		}
		if resp.OrganizationID != 2 {
			t.Errorf("OrganizationID = %d, want 2", resp.OrganizationID)
		}
		if resp.FirstName != "Max" {
			t.Errorf("FirstName = %q, want %q", resp.FirstName, "Max")
		}
		if resp.LastName != "Mustermann" {
			t.Errorf("LastName = %q, want %q", resp.LastName, "Mustermann")
		}
		if resp.Gender != "male" {
			t.Errorf("Gender = %q, want %q", resp.Gender, "male")
		}
		if resp.Contracts != nil {
			t.Errorf("Contracts = %v, want nil", resp.Contracts)
		}
	})

	t.Run("employee with contracts", func(t *testing.T) {
		sectionName := "Kita"
		payPlanName := "TVöD-SuE"
		emp := Employee{
			Person: Person{
				ID:             1,
				OrganizationID: 2,
				FirstName:      "Anna",
				LastName:       "Müller",
				Gender:         "female",
				Birthdate:      time.Date(1985, 3, 20, 0, 0, 0, 0, time.UTC),
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			Contracts: []EmployeeContract{
				{
					ID:         10,
					EmployeeID: 1,
					BaseContract: BaseContract{
						Period:    Period{From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
						SectionID: 5,
						Section:   &Section{Name: sectionName},
					},
					StaffCategory: "qualified",
					Grade:         "S8a",
					Step:          3,
					WeeklyHours:   39.0,
					PayPlanID:     1,
					PayPlan:       &PayPlan{Name: payPlanName},
				},
			},
		}

		resp := emp.ToResponse()

		if len(resp.Contracts) != 1 {
			t.Fatalf("len(Contracts) = %d, want 1", len(resp.Contracts))
		}
		c := resp.Contracts[0]
		if c.ID != 10 {
			t.Errorf("Contracts[0].ID = %d, want 10", c.ID)
		}
		if c.StaffCategory != "qualified" {
			t.Errorf("StaffCategory = %q, want %q", c.StaffCategory, "qualified")
		}
		if c.Grade != "S8a" {
			t.Errorf("Grade = %q, want %q", c.Grade, "S8a")
		}
		if c.Step != 3 {
			t.Errorf("Step = %d, want 3", c.Step)
		}
		if c.WeeklyHours != 39.0 {
			t.Errorf("WeeklyHours = %f, want 39.0", c.WeeklyHours)
		}
	})
}

func TestEmployeeContract_ToResponse(t *testing.T) {
	t.Run("without section or pay plan", func(t *testing.T) {
		contract := EmployeeContract{
			ID:         1,
			EmployeeID: 2,
			BaseContract: BaseContract{
				SectionID: 3,
			},
			StaffCategory: "qualified",
			Grade:         "S8a",
			Step:          2,
			WeeklyHours:   39.0,
			PayPlanID:     4,
		}

		resp := contract.ToResponse()

		if resp.SectionName != nil {
			t.Errorf("SectionName = %v, want nil", resp.SectionName)
		}
		if resp.PayPlanName != nil {
			t.Errorf("PayPlanName = %v, want nil", resp.PayPlanName)
		}
	})

	t.Run("with section and pay plan", func(t *testing.T) {
		contract := EmployeeContract{
			ID:         1,
			EmployeeID: 2,
			BaseContract: BaseContract{
				SectionID: 3,
				Section:   &Section{Name: "Krippe"},
			},
			StaffCategory: "qualified",
			PayPlanID:     4,
			PayPlan:       &PayPlan{Name: "TVöD-SuE"},
		}

		resp := contract.ToResponse()

		if resp.SectionName == nil || *resp.SectionName != "Krippe" {
			t.Errorf("SectionName = %v, want %q", resp.SectionName, "Krippe")
		}
		if resp.PayPlanName == nil || *resp.PayPlanName != "TVöD-SuE" {
			t.Errorf("PayPlanName = %v, want %q", resp.PayPlanName, "TVöD-SuE")
		}
	})
}

func TestEmployeeResponse_FullName(t *testing.T) {
	r := EmployeeResponse{FirstName: "Max", LastName: "Mustermann"}
	if got := r.FullName(); got != "Max Mustermann" {
		t.Errorf("FullName() = %q, want %q", got, "Max Mustermann")
	}
}

func TestEmployeeListFilter_Validate(t *testing.T) {
	t.Run("valid with no filters", func(t *testing.T) {
		f := EmployeeListFilter{}
		if err := f.Validate(); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("valid staff category", func(t *testing.T) {
		cat := "qualified"
		f := EmployeeListFilter{StaffCategory: &cat}
		if err := f.Validate(); err != nil {
			t.Errorf("Validate() error = %v, want nil", err)
		}
	})

	t.Run("invalid staff category", func(t *testing.T) {
		cat := "invalid"
		f := EmployeeListFilter{StaffCategory: &cat}
		if err := f.Validate(); err == nil {
			t.Error("Validate() error = nil, want error for invalid staff category")
		}
	})
}
