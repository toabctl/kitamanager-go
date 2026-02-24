package service

import (
	"context"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestResolveSection_NilName_DefaultSection(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	defaultSection := getDefaultSection(t, db, org.ID)
	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	id, err := resolveSection(ctx, sectionStore, nil, org.ID, cache)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != defaultSection.ID {
		t.Errorf("expected default section ID %d, got %d", defaultSection.ID, id)
	}
}

func TestResolveSection_EmptyString_DefaultSection(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	defaultSection := getDefaultSection(t, db, org.ID)
	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	empty := ""
	id, err := resolveSection(ctx, sectionStore, &empty, org.ID, cache)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != defaultSection.ID {
		t.Errorf("expected default section ID %d, got %d", defaultSection.ID, id)
	}
}

func TestResolveSection_ExistingSection(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	section := createTestSection(t, db, "Krippe", org.ID, false)
	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	name := "Krippe"
	id, err := resolveSection(ctx, sectionStore, &name, org.ID, cache)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != section.ID {
		t.Errorf("expected section ID %d, got %d", section.ID, id)
	}

	// Verify cache was populated after DB lookup
	cachedID, ok := cache["Krippe"]
	if !ok {
		t.Fatal("expected cache to contain 'Krippe' after DB lookup")
	}
	if cachedID != section.ID {
		t.Errorf("cached ID = %d, want %d", cachedID, section.ID)
	}
}

func TestResolveSection_CacheHit(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	section := createTestSection(t, db, "Krippe", org.ID, false)
	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	name := "Krippe"

	// First call populates cache from DB
	id1, err := resolveSection(ctx, sectionStore, &name, org.ID, cache)
	if err != nil {
		t.Fatalf("first call: expected no error, got %v", err)
	}
	if id1 != section.ID {
		t.Errorf("first call: expected section ID %d, got %d", section.ID, id1)
	}

	// Verify cache was populated
	if _, ok := cache["Krippe"]; !ok {
		t.Fatal("expected cache to contain 'Krippe' after first call")
	}

	// Delete the section from DB to prove the second call uses cache, not DB
	if err := db.Unscoped().Delete(&models.Section{}, section.ID).Error; err != nil {
		t.Fatalf("failed to delete section from DB: %v", err)
	}

	// Second call should return cached ID even though DB record is gone
	id2, err := resolveSection(ctx, sectionStore, &name, org.ID, cache)
	if err != nil {
		t.Fatalf("second call: expected no error (cache hit), got %v", err)
	}
	if id2 != section.ID {
		t.Errorf("second call: expected cached section ID %d, got %d", section.ID, id2)
	}
}

func TestResolveSection_AutoCreatesNewSection(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	name := "NewSection"
	id, err := resolveSection(ctx, sectionStore, &name, org.ID, cache)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id == 0 {
		t.Fatal("expected non-zero section ID after auto-create")
	}

	// Verify the section actually exists in the database
	var section models.Section
	if err := db.First(&section, id).Error; err != nil {
		t.Fatalf("auto-created section not found in DB: %v", err)
	}
	if section.Name != "NewSection" {
		t.Errorf("auto-created section name = %q, want %q", section.Name, "NewSection")
	}
	if section.OrganizationID != org.ID {
		t.Errorf("auto-created section org ID = %d, want %d", section.OrganizationID, org.ID)
	}
}

func TestResolveSection_AutoCreate_CachesResult(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	name := "BrandNew"
	id, err := resolveSection(ctx, sectionStore, &name, org.ID, cache)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify cache has the auto-created entry
	cachedID, ok := cache["BrandNew"]
	if !ok {
		t.Fatal("expected cache to contain 'BrandNew' after auto-create")
	}
	if cachedID != id {
		t.Errorf("cached ID = %d, want %d", cachedID, id)
	}

	// Delete from DB to prove subsequent call uses cache
	if err := db.Unscoped().Delete(&models.Section{}, id).Error; err != nil {
		t.Fatalf("failed to delete section from DB: %v", err)
	}

	// Second call should return same ID from cache
	id2, err := resolveSection(ctx, sectionStore, &name, org.ID, cache)
	if err != nil {
		t.Fatalf("second call: expected no error (cache hit), got %v", err)
	}
	if id2 != id {
		t.Errorf("second call: expected cached ID %d, got %d", id, id2)
	}
}

func TestResolveSection_NoDefaultSection(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Create org directly without helper (no default section)
	org := &models.Organization{Name: "No Default Org", Active: true, State: "berlin"}
	if err := db.Create(org).Error; err != nil {
		t.Fatalf("failed to create org: %v", err)
	}

	sectionStore := store.NewSectionStore(db)
	cache := make(map[string]uint)

	_, err := resolveSection(ctx, sectionStore, nil, org.ID, cache)
	if err == nil {
		t.Fatal("expected error when no default section exists, got nil")
	}
}
