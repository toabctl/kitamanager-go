package models

import (
	"testing"
	"time"
)

func TestPayPlan_GetOrganizationID(t *testing.T) {
	pp := PayPlan{OrganizationID: 42}
	if got := pp.GetOrganizationID(); got != 42 {
		t.Errorf("GetOrganizationID() = %d, want 42", got)
	}
}

func TestPayPlan_ToResponse(t *testing.T) {
	pp := PayPlan{
		ID:             1,
		OrganizationID: 2,
		Name:           "TVöD-SuE",
	}

	resp := pp.ToResponse()

	if resp.ID != 1 {
		t.Errorf("ID = %d, want 1", resp.ID)
	}
	if resp.Name != "TVöD-SuE" {
		t.Errorf("Name = %q, want %q", resp.Name, "TVöD-SuE")
	}
}

func TestPayPlan_ToDetailResponse(t *testing.T) {
	minYears := 3
	pp := PayPlan{
		ID:             1,
		OrganizationID: 2,
		Name:           "TVöD-SuE",
		Periods: []PayPlanPeriod{
			{
				ID:          10,
				PayPlanID:   1,
				Period:      Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				WeeklyHours: 39.0,
				Entries: []PayPlanEntry{
					{
						ID:            100,
						PeriodID:      10,
						Grade:         "S8a",
						Step:          3,
						MonthlyAmount: 350000,
						StepMinYears:  &minYears,
					},
				},
			},
		},
	}

	resp := pp.ToDetailResponse()

	if len(resp.Periods) != 1 {
		t.Fatalf("len(Periods) = %d, want 1", len(resp.Periods))
	}
	if resp.Periods[0].WeeklyHours != 39.0 {
		t.Errorf("Periods[0].WeeklyHours = %f, want 39.0", resp.Periods[0].WeeklyHours)
	}
	if len(resp.Periods[0].Entries) != 1 {
		t.Fatalf("len(Entries) = %d, want 1", len(resp.Periods[0].Entries))
	}
	entry := resp.Periods[0].Entries[0]
	if entry.Grade != "S8a" {
		t.Errorf("Grade = %q, want %q", entry.Grade, "S8a")
	}
	if entry.MonthlyAmount != 350000 {
		t.Errorf("MonthlyAmount = %d, want 350000", entry.MonthlyAmount)
	}
	if entry.StepMinYears == nil || *entry.StepMinYears != 3 {
		t.Errorf("StepMinYears = %v, want 3", entry.StepMinYears)
	}
}

func TestPayPlanPeriod_ToResponse_EmployerContribution(t *testing.T) {
	period := PayPlanPeriod{
		ID:                       10,
		PayPlanID:                1,
		Period:                   Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		WeeklyHours:              39.0,
		EmployerContributionRate: 2200, // 22.00%
	}

	resp := period.ToResponse()

	if resp.EmployerContributionRate != 2200 {
		t.Errorf("EmployerContributionRate = %d, want 2200", resp.EmployerContributionRate)
	}
}

func TestPayPlanEntry_ToResponse_NilStepMinYears(t *testing.T) {
	entry := PayPlanEntry{
		ID:            1,
		PeriodID:      10,
		Grade:         "S4",
		Step:          1,
		MonthlyAmount: 280000,
		StepMinYears:  nil,
	}

	resp := entry.ToResponse()

	if resp.StepMinYears != nil {
		t.Errorf("StepMinYears = %v, want nil", resp.StepMinYears)
	}
}
