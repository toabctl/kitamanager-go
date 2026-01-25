package store

import (
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestChildStore_Create(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Emma",
			LastName:       "Schmidt",
			Birthdate:      time.Date(2020, 3, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	err := store.Create(child)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if child.ID == 0 {
		t.Error("expected child ID to be set")
	}
}

func TestChildStore_FindAll(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Child1", LastName: "Last", Birthdate: time.Now()},
	})
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Child2", LastName: "Last", Birthdate: time.Now()},
	})

	children, total, err := store.FindAll(100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestChildStore_FindByID(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Test",
			LastName:       "Child",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	found, err := store.FindByID(child.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.FirstName != "Test" {
		t.Errorf("expected first name 'Test', got '%s'", found.FirstName)
	}
}

func TestChildStore_FindByOrganization(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Child1", LastName: "Last", Birthdate: time.Now()},
	})
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Child2", LastName: "Last", Birthdate: time.Now()},
	})
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org2.ID, FirstName: "Child3", LastName: "Last", Birthdate: time.Now()},
	})

	children, total, err := store.FindByOrganization(org1.ID, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(children) != 2 {
		t.Errorf("expected 2 children for org1, got %d", len(children))
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

func TestChildStore_Update(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Original",
			LastName:       "Name",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	child.FirstName = "Updated"
	err := store.Update(child)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := store.FindByID(child.ID)
	if found.FirstName != "Updated" {
		t.Errorf("expected first name 'Updated', got '%s'", found.FirstName)
	}
}

func TestChildStore_Delete(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "ToDelete",
			LastName:       "Child",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	err := store.Delete(child.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.FindByID(child.ID)
	if err == nil {
		t.Error("expected error finding deleted child")
	}
}

func TestChildStore_CreateContract(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Test",
			LastName:       "Child",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	contract := &models.ChildContract{
		ChildID: child.ID,
		Period: models.Period{
			From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   nil,
		},
		CareHoursPerWeek: 35,
		MealsIncluded:    true,
	}

	err := store.CreateContract(contract)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if contract.ID == 0 {
		t.Error("expected contract ID to be set")
	}
}

func TestChildStore_DeleteContract(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Test",
			LastName:       "Child",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	contract := &models.ChildContract{
		ChildID: child.ID,
		Period: models.Period{
			From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		CareHoursPerWeek: 35,
	}
	db.Create(contract)

	err := store.DeleteContract(contract.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = store.FindContractByID(contract.ID)
	if err == nil {
		t.Error("expected error finding deleted contract")
	}
}

func TestChildStore_DeleteAlsoDeletesContracts(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Test",
			LastName:       "Child",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	contract := &models.ChildContract{
		ChildID: child.ID,
		Period: models.Period{
			From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		CareHoursPerWeek: 35,
	}
	db.Create(contract)
	contractID := contract.ID

	err := store.Delete(child.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify contract is also deleted
	_, err = store.FindContractByID(contractID)
	if err == nil {
		t.Error("expected contract to be deleted with child")
	}
}

func TestChildStore_ContractOverlapValidation(t *testing.T) {
	db := setupTestDB(t)
	store := NewChildStore(db)
	org := createTestOrganization(t, db, "Test Org")

	child := &models.Child{
		Person: models.Person{
			OrganizationID: org.ID,
			FirstName:      "Test",
			LastName:       "Child",
			Birthdate:      time.Now(),
		},
	}
	db.Create(child)

	// Create existing contract
	existing := &models.ChildContract{
		ChildID: child.ID,
		Period: models.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   datePtr(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
		},
		CareHoursPerWeek: 35,
	}
	db.Create(existing)

	// Try to create overlapping contract
	err := store.Contracts().ValidateNoOverlap(
		child.ID,
		time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		datePtr(time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)),
		nil,
	)

	if err == nil {
		t.Error("expected overlap error")
	}

	// Non-overlapping contract should succeed
	err = store.Contracts().ValidateNoOverlap(
		child.ID,
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		nil,
		nil,
	)

	if err != nil {
		t.Errorf("expected no error for non-overlapping contract, got %v", err)
	}
}
