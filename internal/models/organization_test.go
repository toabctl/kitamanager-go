package models

import (
	"testing"
	"time"
)

func TestOrganization_ToResponse(t *testing.T) {
	now := time.Now()
	org := Organization{
		ID:        1,
		Name:      "Kita Sonnenschein",
		Active:    true,
		State:     "berlin",
		CreatedAt: now,
		CreatedBy: "admin@example.com",
		UpdatedAt: now,
	}

	resp := org.ToResponse()

	if resp.ID != 1 {
		t.Errorf("ID = %d, want 1", resp.ID)
	}
	if resp.Name != "Kita Sonnenschein" {
		t.Errorf("Name = %q, want %q", resp.Name, "Kita Sonnenschein")
	}
	if resp.Active != true {
		t.Errorf("Active = %v, want true", resp.Active)
	}
	if resp.State != "berlin" {
		t.Errorf("State = %q, want %q", resp.State, "berlin")
	}
	if resp.CreatedBy != "admin@example.com" {
		t.Errorf("CreatedBy = %q, want %q", resp.CreatedBy, "admin@example.com")
	}
}
