package store

import (
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestGroupStore_Create(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	group := &models.Group{
		Name:      "Test Group",
		Active:    true,
		CreatedBy: "admin@test.com",
	}

	err := store.Create(group)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if group.ID == 0 {
		t.Error("expected group ID to be set")
	}
}

func TestGroupStore_FindAll(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	createTestGroup(t, db, "Group 1")
	createTestGroup(t, db, "Group 2")

	groups, err := store.FindAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupStore_FindByID(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	created := createTestGroup(t, db, "Test Group")

	found, err := store.FindByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.Name != "Test Group" {
		t.Errorf("expected name 'Test Group', got '%s'", found.Name)
	}
}

func TestGroupStore_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	_, err := store.FindByID(999)
	if err == nil {
		t.Error("expected error for non-existent group")
	}
}

func TestGroupStore_Update(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	group := createTestGroup(t, db, "Original Name")
	group.Name = "Updated Name"

	err := store.Update(group)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := store.FindByID(group.ID)
	if found.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", found.Name)
	}
}

func TestGroupStore_Delete(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	group := createTestGroup(t, db, "To Delete")

	err := store.Delete(group.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.FindByID(group.ID)
	if err == nil {
		t.Error("expected error finding deleted group")
	}
}

func TestGroupStore_AddToOrganization(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	group := createTestGroup(t, db, "Test Group")
	org := createTestOrganization(t, db, "Test Org")

	err := store.AddToOrganization(group.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := store.FindByID(group.ID)
	if len(found.Organizations) != 1 {
		t.Errorf("expected 1 organization, got %d", len(found.Organizations))
	}
}

func TestGroupStore_RemoveFromOrganization(t *testing.T) {
	db := setupTestDB(t)
	store := NewGroupStore(db)

	group := createTestGroup(t, db, "Test Group")
	org := createTestOrganization(t, db, "Test Org")

	_ = store.AddToOrganization(group.ID, org.ID)

	err := store.RemoveFromOrganization(group.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := store.FindByID(group.ID)
	if len(found.Organizations) != 0 {
		t.Errorf("expected 0 organizations, got %d", len(found.Organizations))
	}
}
