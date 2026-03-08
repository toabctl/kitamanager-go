package models

import (
	"testing"
	"time"
)

func TestSection_ToResponse(t *testing.T) {
	now := time.Now()
	minAge := 0
	maxAge := 36

	section := Section{
		ID:             1,
		OrganizationID: 2,
		Name:           "Krippe",
		IsDefault:      false,
		MinAgeMonths:   &minAge,
		MaxAgeMonths:   &maxAge,
		CreatedAt:      now,
		CreatedBy:      "admin@example.com",
		UpdatedAt:      now,
	}

	resp := section.ToResponse()

	if resp.ID != 1 {
		t.Errorf("ID = %d, want 1", resp.ID)
	}
	if resp.OrganizationID != 2 {
		t.Errorf("OrganizationID = %d, want 2", resp.OrganizationID)
	}
	if resp.Name != "Krippe" {
		t.Errorf("Name = %q, want %q", resp.Name, "Krippe")
	}
	if resp.IsDefault != false {
		t.Errorf("IsDefault = %v, want false", resp.IsDefault)
	}
	if resp.MinAgeMonths == nil || *resp.MinAgeMonths != 0 {
		t.Errorf("MinAgeMonths = %v, want 0", resp.MinAgeMonths)
	}
	if resp.MaxAgeMonths == nil || *resp.MaxAgeMonths != 36 {
		t.Errorf("MaxAgeMonths = %v, want 36", resp.MaxAgeMonths)
	}
}

func TestSection_ToResponse_NilAgeFields(t *testing.T) {
	section := Section{
		ID:             1,
		OrganizationID: 2,
		Name:           "Default",
		IsDefault:      true,
	}

	resp := section.ToResponse()

	if resp.MinAgeMonths != nil {
		t.Errorf("MinAgeMonths = %v, want nil", resp.MinAgeMonths)
	}
	if resp.MaxAgeMonths != nil {
		t.Errorf("MaxAgeMonths = %v, want nil", resp.MaxAgeMonths)
	}
}

func TestSection_GetOrganizationID(t *testing.T) {
	section := Section{OrganizationID: 42}
	if got := section.GetOrganizationID(); got != 42 {
		t.Errorf("GetOrganizationID() = %d, want 42", got)
	}
}
