package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestChildService_List(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	createTestChild(t, db, "John", "Doe", org.ID)
	createTestChild(t, db, "Jane", "Doe", org.ID)

	children, total, err := svc.List(ctx, 10, 0)
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

func TestChildService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	found, err := svc.GetByID(ctx, child.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID != child.ID {
		t.Errorf("ID = %d, want %d", found.ID, child.ID)
	}
	if found.FirstName != "John" {
		t.Errorf("FirstName = %v, want John", found.FirstName)
	}
}

func TestChildService_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	_, err := svc.GetByID(ctx, 999, org.ID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// SECURITY TEST: Cross-organization access attempt
func TestChildService_GetByID_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	// Try to access child from wrong organization - should return not found
	_, err := svc.GetByID(ctx, child.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error when accessing child from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}
}

func TestChildService_Create(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	req := &models.ChildCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Birthdate: time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	child, err := svc.Create(ctx, org.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if child.ID == 0 {
		t.Error("expected ID to be set")
	}
	if child.FirstName != "John" {
		t.Errorf("FirstName = %v, want John", child.FirstName)
	}
	if child.LastName != "Doe" {
		t.Errorf("LastName = %v, want Doe", child.LastName)
	}
	if child.OrganizationID != org.ID {
		t.Errorf("OrganizationID = %d, want %d", child.OrganizationID, org.ID)
	}
}

func TestChildService_Create_WhitespaceOnlyNames(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	tests := []struct {
		name string
		req  *models.ChildCreateRequest
	}{
		{"empty first name", &models.ChildCreateRequest{FirstName: "", LastName: "Doe", Birthdate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}},
		{"whitespace first name", &models.ChildCreateRequest{FirstName: "   ", LastName: "Doe", Birthdate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}},
		{"empty last name", &models.ChildCreateRequest{FirstName: "John", LastName: "", Birthdate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}},
		{"whitespace last name", &models.ChildCreateRequest{FirstName: "John", LastName: "   ", Birthdate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(ctx, org.ID, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !errors.Is(err, apperror.ErrBadRequest) {
				t.Errorf("expected ErrBadRequest, got %v", err)
			}
		})
	}
}

func TestChildService_Create_TrimmedNames(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	req := &models.ChildCreateRequest{
		FirstName: "  John  ",
		LastName:  "  Doe  ",
		Birthdate: time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	child, err := svc.Create(ctx, org.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if child.FirstName != "John" {
		t.Errorf("FirstName = %v, want 'John' (trimmed)", child.FirstName)
	}
	if child.LastName != "Doe" {
		t.Errorf("LastName = %v, want 'Doe' (trimmed)", child.LastName)
	}
}

func TestChildService_Create_FutureBirthdate(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	req := &models.ChildCreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Birthdate: time.Now().AddDate(1, 0, 0), // 1 year in future
	}

	_, err := svc.Create(ctx, org.ID, req)
	if err == nil {
		t.Fatal("expected error for future birthdate, got nil")
	}

	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestChildService_Update(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	newFirstName := "Jane"
	req := &models.ChildUpdateRequest{
		FirstName: &newFirstName,
	}

	updated, err := svc.Update(ctx, child.ID, org.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.FirstName != "Jane" {
		t.Errorf("FirstName = %v, want Jane", updated.FirstName)
	}
	if updated.LastName != "Doe" {
		t.Errorf("LastName should not change, got %v", updated.LastName)
	}
}

func TestChildService_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	newName := "Jane"
	req := &models.ChildUpdateRequest{
		FirstName: &newName,
	}

	_, err := svc.Update(ctx, 999, org.ID, req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// SECURITY TEST: Cross-organization update attempt
func TestChildService_Update_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	newName := "Hacked"
	req := &models.ChildUpdateRequest{
		FirstName: &newName,
	}

	// Try to update child from wrong organization
	_, err := svc.Update(ctx, child.ID, org2.ID, req)
	if err == nil {
		t.Fatal("expected error when updating child from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}

	// Verify child was not actually updated
	found, _ := svc.GetByID(ctx, child.ID, org1.ID)
	if found.FirstName != "John" {
		t.Errorf("child was modified despite wrong org, FirstName = %v", found.FirstName)
	}
}

func TestChildService_Delete(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	err := svc.Delete(ctx, child.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it's deleted
	_, err = svc.GetByID(ctx, child.ID, org.ID)
	if err == nil {
		t.Error("expected child to be deleted")
	}
}

// SECURITY TEST: Cross-organization delete attempt
func TestChildService_Delete_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	// Try to delete child from wrong organization
	err := svc.Delete(ctx, child.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error when deleting child from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}

	// Verify child still exists
	found, err := svc.GetByID(ctx, child.ID, org1.ID)
	if err != nil {
		t.Errorf("child was deleted despite wrong org: %v", err)
	}
	if found.FirstName != "John" {
		t.Error("child data was corrupted")
	}
}

func TestChildService_CreateContract(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	req := &models.ChildContractCreateRequest{
		From:       from,
		To:         &to,
		Attributes: []string{"ganztags", "ndh"},
	}

	contract, err := svc.CreateContract(ctx, child.ID, org.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if contract.ID == 0 {
		t.Error("expected ID to be set")
	}
	if contract.ChildID != child.ID {
		t.Errorf("ChildID = %d, want %d", contract.ChildID, child.ID)
	}
	if len(contract.Attributes) != 2 {
		t.Errorf("Attributes = %v, want 2 elements", contract.Attributes)
	}
}

func TestChildService_CreateContract_ChildNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &models.ChildContractCreateRequest{
		From: from,
	}

	_, err := svc.CreateContract(ctx, 999, org.ID, req)
	if err == nil {
		t.Fatal("expected error for non-existent child, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// SECURITY TEST: Cross-organization contract creation attempt
func TestChildService_CreateContract_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &models.ChildContractCreateRequest{
		From: from,
	}

	// Try to create contract via wrong organization
	_, err := svc.CreateContract(ctx, child.ID, org2.ID, req)
	if err == nil {
		t.Fatal("expected error when creating contract from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}
}

func TestChildService_CreateContract_InvalidPeriod(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	from := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Before from

	req := &models.ChildContractCreateRequest{
		From: from,
		To:   &to,
	}

	_, err := svc.CreateContract(ctx, child.ID, org.ID, req)
	if err == nil {
		t.Fatal("expected error for invalid period (to before from), got nil")
	}

	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestChildService_CreateContract_OverlappingContract(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	// Create first contract
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	req1 := &models.ChildContractCreateRequest{
		From:       from1,
		To:         &to1,
		Attributes: []string{"ganztags"},
	}
	_, err := svc.CreateContract(ctx, child.ID, org.ID, req1)
	if err != nil {
		t.Fatalf("first contract: expected no error, got %v", err)
	}

	// Try to create overlapping contract
	from2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC) // Overlaps with first
	to2 := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	req2 := &models.ChildContractCreateRequest{
		From:       from2,
		To:         &to2,
		Attributes: []string{"teilzeit"},
	}

	_, err = svc.CreateContract(ctx, child.ID, org.ID, req2)
	if err == nil {
		t.Fatal("expected error for overlapping contract, got nil")
	}

	if !errors.Is(err, apperror.ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestChildService_CreateContract_OngoingContract(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// No 'to' date means ongoing contract
	req := &models.ChildContractCreateRequest{
		From: from,
		To:   nil,
	}

	contract, err := svc.CreateContract(ctx, child.ID, org.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if contract.To != nil {
		t.Errorf("To = %v, want nil (ongoing)", contract.To)
	}
}

func TestChildService_ListContracts(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	// Create two contracts
	from1 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
	req1 := &models.ChildContractCreateRequest{From: from1, To: &to1, Attributes: []string{"teilzeit"}}
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, req1)

	from2 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	req2 := &models.ChildContractCreateRequest{From: from2, Attributes: []string{"ganztags"}}
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, req2)

	contracts, err := svc.ListContracts(ctx, child.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(contracts) != 2 {
		t.Errorf("expected 2 contracts, got %d", len(contracts))
	}
}

func TestChildService_ListContracts_ChildNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	_, err := svc.ListContracts(ctx, 999, org.ID)
	if err == nil {
		t.Fatal("expected error for non-existent child, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// SECURITY TEST: Cross-organization list contracts attempt
func TestChildService_ListContracts_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	// Create a contract for child in org1
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &models.ChildContractCreateRequest{From: from, Attributes: []string{"ganztags"}}
	_, _ = svc.CreateContract(ctx, child.ID, org1.ID, req)

	// Try to list contracts from wrong organization
	_, err := svc.ListContracts(ctx, child.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error when listing contracts from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}
}

func TestChildService_ListByOrganization(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	createTestChild(t, db, "John", "Doe", org1.ID)
	createTestChild(t, db, "Jane", "Doe", org1.ID)
	createTestChild(t, db, "Bob", "Smith", org2.ID)

	children, total, err := svc.ListByOrganization(ctx, org1.ID, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(children) != 2 {
		t.Errorf("expected 2 children in org1, got %d", len(children))
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
}

// SECURITY TEST: Verify ListByOrganization doesn't leak data from other orgs
func TestChildService_ListByOrganization_IsolatesData(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create children in both orgs
	createTestChild(t, db, "John", "Doe", org1.ID)
	createTestChild(t, db, "Secret", "Child", org2.ID)

	// List children in org1
	children, _, err := svc.ListByOrganization(ctx, org1.ID, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify no children from org2 are returned
	for _, child := range children {
		if child.OrganizationID != org1.ID {
			t.Errorf("got child from wrong org: %d (expected %d)", child.OrganizationID, org1.ID)
		}
		if child.FirstName == "Secret" {
			t.Error("data leaked from other organization")
		}
	}
}

// SECURITY TEST: DeleteContract cross-org
func TestChildService_DeleteContract_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	// Create a contract
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &models.ChildContractCreateRequest{From: from, Attributes: []string{"ganztags"}}
	contract, _ := svc.CreateContract(ctx, child.ID, org1.ID, req)

	// Try to delete contract from wrong organization
	err := svc.DeleteContract(ctx, contract.ID, child.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error when deleting contract from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}

	// Verify contract still exists
	contracts, _ := svc.ListContracts(ctx, child.ID, org1.ID)
	if len(contracts) != 1 {
		t.Error("contract was deleted despite wrong org")
	}
}

// SECURITY TEST: GetCurrentContract cross-org
func TestChildService_GetCurrentContract_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	child := createTestChild(t, db, "John", "Doe", org1.ID)

	// Create an ongoing contract
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	req := &models.ChildContractCreateRequest{From: from, Attributes: []string{"ganztags"}}
	_, _ = svc.CreateContract(ctx, child.ID, org1.ID, req)

	// Try to get current contract from wrong organization
	_, err := svc.GetCurrentContract(ctx, child.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error when getting current contract from wrong org, got nil")
	}

	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound (not forbidden - security), got %v", err)
	}
}

// =========================================
// Funding Calculation Tests
// =========================================

func TestChildService_CalculateFunding_BasicCalculation(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	// Create org with government funding
	org := createTestOrganization(t, db, "Test Org")
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")

	// Assign funding to org
	db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("government_funding_id", funding.ID)

	// Create funding period covering our test date
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)

	// Create properties with age filter (ages 3-6)
	createTestFundingProperty(t, db, period.ID, "ganztags", 100000, 3, 7) // 1000.00 EUR
	createTestFundingProperty(t, db, period.ID, "ndh", 50000, 3, 7)       // 500.00 EUR

	// Create child (born 2022-01-15, age 3 on 2025-01-27)
	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	child.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child)

	// Create contract with attributes
	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags", "ndh"},
	})
	if err != nil {
		t.Fatalf("failed to create contract: %v", err)
	}

	// Calculate funding
	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(result.Children))
	}

	cf := result.Children[0]
	if cf.ChildID != child.ID {
		t.Errorf("ChildID = %d, want %d", cf.ChildID, child.ID)
	}
	if cf.Funding != 150000 { // 1000.00 + 500.00 = 1500.00 EUR = 150000 cents
		t.Errorf("Funding = %d, want 150000 (cents)", cf.Funding)
	}
	if len(cf.MatchedAttributes) != 2 {
		t.Errorf("MatchedAttributes = %v, want 2 items", cf.MatchedAttributes)
	}
	if len(cf.UnmatchedAttributes) != 0 {
		t.Errorf("UnmatchedAttributes = %v, want 0 items", cf.UnmatchedAttributes)
	}
}

func TestChildService_CalculateFunding_NoFundingAssigned(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	// Create org WITHOUT government funding
	org := createTestOrganization(t, db, "Test Org")

	// Create child with contract
	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	child.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child)

	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags", "ndh"},
	})
	if err != nil {
		t.Fatalf("failed to create contract: %v", err)
	}

	// Calculate funding
	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(result.Children))
	}

	cf := result.Children[0]
	if cf.Funding != 0 {
		t.Errorf("Funding = %d, want 0 (no funding assigned)", cf.Funding)
	}
	if len(cf.UnmatchedAttributes) != 2 {
		t.Errorf("UnmatchedAttributes = %v, want 2 items", cf.UnmatchedAttributes)
	}
}

func TestChildService_CalculateFunding_NoMatchingPeriod(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	// Create org with funding
	org := createTestOrganization(t, db, "Test Org")
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("government_funding_id", funding.ID)

	// Create period that doesn't cover our test date
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &to)

	// Create child with contract
	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	child.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child)

	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags"},
	})

	// Calculate funding for date outside period
	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cf := result.Children[0]
	if cf.Funding != 0 {
		t.Errorf("Funding = %d, want 0 (no matching period)", cf.Funding)
	}
}

func TestChildService_CalculateFunding_NoMatchingAgeProperty(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	// Create org with funding
	org := createTestOrganization(t, db, "Test Org")
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("government_funding_id", funding.ID)

	// Create period
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)

	// Create property for ages 0-2 only
	createTestFundingProperty(t, db, period.ID, "ganztags", 100000, 0, 2)

	// Create child age 3 (doesn't match 0-2 property)
	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	child.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC) // Age 3 on 2025-01-27
	db.Save(child)

	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags"},
	})

	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cf := result.Children[0]
	if cf.Funding != 0 {
		t.Errorf("Funding = %d, want 0 (no matching age property)", cf.Funding)
	}
	if len(cf.UnmatchedAttributes) != 1 {
		t.Errorf("UnmatchedAttributes = %v, want [ganztags]", cf.UnmatchedAttributes)
	}
}

func TestChildService_CalculateFunding_PartialAttributeMatch(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	// Create org with funding
	org := createTestOrganization(t, db, "Test Org")
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("government_funding_id", funding.ID)

	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	createTestFundingProperty(t, db, period.ID, "ganztags", 100000, 3, 7)
	// "xyz" property does NOT exist

	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	child.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child)

	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags", "xyz"}, // xyz doesn't match any property
	})

	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cf := result.Children[0]
	if cf.Funding != 100000 {
		t.Errorf("Funding = %d, want 100000", cf.Funding)
	}
	if len(cf.MatchedAttributes) != 1 || cf.MatchedAttributes[0] != "ganztags" {
		t.Errorf("MatchedAttributes = %v, want [ganztags]", cf.MatchedAttributes)
	}
	if len(cf.UnmatchedAttributes) != 1 || cf.UnmatchedAttributes[0] != "xyz" {
		t.Errorf("UnmatchedAttributes = %v, want [xyz]", cf.UnmatchedAttributes)
	}
}

func TestChildService_CalculateFunding_DuplicateAttributes(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("government_funding_id", funding.ID)

	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	createTestFundingProperty(t, db, period.ID, "ganztags", 100000, 3, 7)

	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	child.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child)

	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags", "ganztags", "ganztags"}, // Duplicates
	})

	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cf := result.Children[0]
	// Should only count once even with duplicates
	if cf.Funding != 100000 {
		t.Errorf("Funding = %d, want 100000 (counted once despite duplicates)", cf.Funding)
	}
	if len(cf.MatchedAttributes) != 1 {
		t.Errorf("MatchedAttributes = %v, want 1 item (deduplicated)", cf.MatchedAttributes)
	}
}

func TestChildService_CalculateFunding_ChildNoContractOnDate(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("government_funding_id", funding.ID)

	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	createTestFundingProperty(t, db, period.ID, "ganztags", 100000, 3, 7)

	// Child with active contract
	childActive := createTestChild(t, db, "Active", "Child", org.ID)
	childActive.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(childActive)
	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, childActive.ID, org.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags"},
	})

	// Child with NO contract (should not appear in results)
	childNoContract := createTestChild(t, db, "NoContract", "Child", org.ID)
	childNoContract.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(childNoContract)

	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should only include child with active contract
	if len(result.Children) != 1 {
		t.Errorf("expected 1 child (with active contract), got %d", len(result.Children))
	}
	if result.Children[0].ChildName != "Active Child" {
		t.Errorf("expected Active Child, got %s", result.Children[0].ChildName)
	}
}

// SECURITY TEST: Cross-organization funding calculation
func TestChildService_CalculateFunding_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	funding := createTestGovernmentFunding(t, db, "Funding")
	db.Model(&models.Organization{}).Where("id = ?", org1.ID).Update("government_funding_id", funding.ID)
	db.Model(&models.Organization{}).Where("id = ?", org2.ID).Update("government_funding_id", funding.ID)

	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil)
	createTestFundingProperty(t, db, period.ID, "ganztags", 100000, 3, 7)

	// Child in org1
	child1 := createTestChild(t, db, "Org1", "Child", org1.ID)
	child1.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child1)
	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child1.ID, org1.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags"},
	})

	// Child in org2
	child2 := createTestChild(t, db, "Org2", "Child", org2.ID)
	child2.Birthdate = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)
	db.Save(child2)
	_, _ = svc.CreateContract(ctx, child2.ID, org2.ID, &models.ChildContractCreateRequest{
		From:       fromDate,
		Attributes: []string{"ganztags"},
	})

	// Calculate funding for org1 - should NOT include org2's child
	refDate := time.Date(2025, 1, 27, 0, 0, 0, 0, time.UTC)
	result, err := svc.CalculateFunding(ctx, org1.ID, refDate)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Children) != 1 {
		t.Errorf("expected 1 child from org1, got %d", len(result.Children))
	}

	for _, cf := range result.Children {
		if cf.ChildName == "Org2 Child" {
			t.Error("SECURITY: org2's child leaked into org1's funding calculation")
		}
	}
}

// =========================================
// Monthly Statistics Tests
// =========================================

func TestChildService_GetContractCountByMonth_BasicCounting(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	// Create contract starting 2025-01-01
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From: from,
	})
	if err != nil {
		t.Fatalf("failed to create contract: %v", err)
	}

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2025, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stats.Years) != 1 {
		t.Fatalf("expected 1 year, got %d", len(stats.Years))
	}

	if stats.Years[0].Year != 2025 {
		t.Errorf("expected year 2025, got %d", stats.Years[0].Year)
	}

	// All months should have 1 child
	for month, count := range stats.Years[0].Counts {
		if count != 1 {
			t.Errorf("expected count 1 for month %d, got %d", month+1, count)
		}
	}
}

func TestChildService_GetContractCountByMonth_ContractEndsJuly(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	// Create contract ending July 31, 2025
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From: from,
		To:   &to,
	})
	if err != nil {
		t.Fatalf("failed to create contract: %v", err)
	}

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2025, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Jan-Jul (0-6) should have 1, Aug-Dec (7-11) should have 0
	for month, count := range stats.Years[0].Counts {
		if month <= 6 { // Jan-Jul
			if count != 1 {
				t.Errorf("expected count 1 for month %d (contract active), got %d", month+1, count)
			}
		} else { // Aug-Dec
			if count != 0 {
				t.Errorf("expected count 0 for month %d (contract ended), got %d", month+1, count)
			}
		}
	}
}

func TestChildService_GetContractCountByMonth_MultipleYears(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	child := createTestChild(t, db, "John", "Doe", org.ID)

	// Create ongoing contract
	from := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child.ID, org.ID, &models.ChildContractCreateRequest{
		From: from,
	})

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2023, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(stats.Years) != 3 {
		t.Errorf("expected 3 years, got %d", len(stats.Years))
	}

	// Check 2023: Jan-May = 0, Jun-Dec = 1
	for month, count := range stats.Years[0].Counts {
		expected := 0
		if month >= 5 { // Jun onwards
			expected = 1
		}
		if count != expected {
			t.Errorf("2023 month %d: expected %d, got %d", month+1, expected, count)
		}
	}

	// 2024 and 2025 should all be 1
	for _, yearData := range stats.Years[1:] {
		for month, count := range yearData.Counts {
			if count != 1 {
				t.Errorf("year %d month %d: expected 1, got %d", yearData.Year, month+1, count)
			}
		}
	}
}

func TestChildService_GetContractCountByMonth_NoChildren(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2025, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// All counts should be 0
	for _, count := range stats.Years[0].Counts {
		if count != 0 {
			t.Errorf("expected count 0 for no children, got %d", count)
		}
	}
}

func TestChildService_GetContractCountByMonth_ChildWithNoContract(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	createTestChild(t, db, "John", "Doe", org.ID) // No contract

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2025, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// All counts should be 0 (child has no contract)
	for month, count := range stats.Years[0].Counts {
		if count != 0 {
			t.Errorf("expected count 0 for month %d (no contract), got %d", month+1, count)
		}
	}
}

// SECURITY TEST: Cross-organization isolation
func TestChildService_GetContractCountByMonth_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1 with contract
	child := createTestChild(t, db, "John", "Doe", org1.ID)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child.ID, org1.ID, &models.ChildContractCreateRequest{
		From: from,
	})

	// Query stats for org2 - should not include org1's children
	stats, err := svc.GetContractCountByMonth(ctx, org2.ID, 2025, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for month, count := range stats.Years[0].Counts {
		if count != 0 {
			t.Errorf("SECURITY: expected count 0 for org2 month %d (child in org1), got %d", month+1, count)
		}
	}
}

func TestChildService_GetContractCountByMonth_MultipleChildrenDifferentContracts(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	// Child 1: contract Jan 2025 - Jul 2025
	child1 := createTestChild(t, db, "Child1", "Test", org.ID)
	from1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child1.ID, org.ID, &models.ChildContractCreateRequest{
		From: from1,
		To:   &to1,
	})

	// Child 2: contract Aug 2025 onwards
	child2 := createTestChild(t, db, "Child2", "Test", org.ID)
	from2 := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child2.ID, org.ID, &models.ChildContractCreateRequest{
		From: from2,
	})

	// Child 3: contract all year
	child3 := createTestChild(t, db, "Child3", "Test", org.ID)
	from3 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _ = svc.CreateContract(ctx, child3.ID, org.ID, &models.ChildContractCreateRequest{
		From: from3,
	})

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2025, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{
		2, // Jan: child1 + child3
		2, // Feb
		2, // Mar
		2, // Apr
		2, // May
		2, // Jun
		2, // Jul: child1 + child3
		2, // Aug: child2 + child3
		2, // Sep
		2, // Oct
		2, // Nov
		2, // Dec
	}

	for month, count := range stats.Years[0].Counts {
		if count != expected[month] {
			t.Errorf("month %d: expected %d, got %d", month+1, expected[month], count)
		}
	}
}

func TestChildService_GetContractCountByMonth_PeriodFormat(t *testing.T) {
	db := setupTestDB(t)
	svc := createChildService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")

	stats, err := svc.GetContractCountByMonth(ctx, org.ID, 2023, 2025)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if stats.Period.Start != "2023-01-01" {
		t.Errorf("expected period start '2023-01-01', got '%s'", stats.Period.Start)
	}
	if stats.Period.End != "2025-12-31" {
		t.Errorf("expected period end '2025-12-31', got '%s'", stats.Period.End)
	}
}
