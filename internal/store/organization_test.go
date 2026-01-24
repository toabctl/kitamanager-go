package store

import (
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestOrganizationStore_Create(t *testing.T) {
	db := setupTestDB(t)
	store := NewOrganizationStore(db)

	org := &models.Organization{
		Name:      "Test Org",
		Active:    true,
		CreatedBy: "admin@test.com",
	}

	err := store.Create(org)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if org.ID == 0 {
		t.Error("expected organization ID to be set")
	}
}

func TestOrganizationStore_FindAll(t *testing.T) {
	db := setupTestDB(t)
	store := NewOrganizationStore(db)

	// Create test organizations
	createTestOrganization(t, db, "Org 1")
	createTestOrganization(t, db, "Org 2")

	orgs, err := store.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(orgs) != 2 {
		t.Errorf("expected 2 organizations, got %d", len(orgs))
	}
}

func TestOrganizationStore_FindByID(t *testing.T) {
	db := setupTestDB(t)
	store := NewOrganizationStore(db)

	created := createTestOrganization(t, db, "Test Org")

	found, err := store.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.Name != "Test Org" {
		t.Errorf("expected name 'Test Org', got '%s'", found.Name)
	}
}

func TestOrganizationStore_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	store := NewOrganizationStore(db)

	_, err := store.FindByID(999)
	if err == nil {
		t.Error("expected error for non-existent organization")
	}
}

func TestOrganizationStore_Update(t *testing.T) {
	db := setupTestDB(t)
	store := NewOrganizationStore(db)

	org := createTestOrganization(t, db, "Original Name")
	org.Name = "Updated Name"

	err := store.Update(org)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := store.FindByID(org.ID)
	if found.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", found.Name)
	}
}

func TestOrganizationStore_Delete(t *testing.T) {
	db := setupTestDB(t)
	store := NewOrganizationStore(db)

	org := createTestOrganization(t, db, "To Delete")

	err := store.Delete(org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.FindByID(org.ID)
	if err == nil {
		t.Error("expected error finding deleted organization")
	}
}
