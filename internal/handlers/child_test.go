package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestChildHandler_List(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Child1", LastName: "Last", Birthdate: time.Now()},
	})
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Child2", LastName: "Last", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	w := performRequest(r, "GET", "/organizations/1/children", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if len(response.Data) != 2 {
		t.Errorf("expected 2 children, got %d", len(response.Data))
	}
}

func TestChildHandler_Get(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/1/children/1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var result models.Child
	parseResponse(t, w, &result)

	if result.FirstName != "Test" {
		t.Errorf("expected first name 'Test', got '%s'", result.FirstName)
	}
}

func TestChildHandler_Create(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := models.ChildCreateRequest{
		FirstName: "New",
		LastName:  "Child",
		Birthdate: time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var result models.Child
	parseResponse(t, w, &result)

	if result.FirstName != "New" {
		t.Errorf("expected first name 'New', got '%s'", result.FirstName)
	}
	if result.OrganizationID != org.ID {
		t.Errorf("expected organization ID %d, got %d", org.ID, result.OrganizationID)
	}
}

func TestChildHandler_Update(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Original", LastName: "Child", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.PUT("/organizations/:orgId/children/:id", handler.Update)

	newName := "Updated"
	body := models.ChildUpdateRequest{
		FirstName: &newName,
	}

	w := performRequest(r, "PUT", "/organizations/1/children/1", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var result models.Child
	parseResponse(t, w, &result)

	if result.FirstName != "Updated" {
		t.Errorf("expected first name 'Updated', got '%s'", result.FirstName)
	}
}

func TestChildHandler_Delete(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "ToDelete", LastName: "Child", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/organizations/1/children/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestChildHandler_ListContracts(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	db.Create(&models.ChildContract{
		ChildID: child.ID,
		Period:  models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	})

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts", handler.ListContracts)

	w := performRequest(r, "GET", "/organizations/1/children/1/contracts", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var contracts []models.ChildContract
	parseResponse(t, w, &contracts)

	if len(contracts) != 1 {
		t.Errorf("expected 1 contract, got %d", len(contracts))
	}
}

func TestChildHandler_GetCurrentContract(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	db.Create(&models.ChildContract{
		ChildID:    child.ID,
		Period:     models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), To: nil},
		Attributes: []string{"ganztags"},
	})

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts/current", handler.GetCurrentContract)

	w := performRequest(r, "GET", "/organizations/1/children/1/contracts/current", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var contract models.ChildContract
	parseResponse(t, w, &contract)

	if len(contract.Attributes) != 1 || contract.Attributes[0] != "ganztags" {
		t.Errorf("expected attributes [ganztags], got %v", contract.Attributes)
	}
}

func TestChildHandler_GetCurrentContract_NotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts/current", handler.GetCurrentContract)

	w := performRequest(r, "GET", "/organizations/1/children/1/contracts/current", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_CreateContract(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	body := models.ChildContractCreateRequest{
		From:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		To:         nil,
		Attributes: []string{"ganztags", "ndh"},
	}

	w := performRequest(r, "POST", "/organizations/1/children/1/contracts", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var contract models.ChildContract
	parseResponse(t, w, &contract)

	if len(contract.Attributes) != 2 {
		t.Errorf("expected 2 attributes, got %d", len(contract.Attributes))
	}
}

func TestChildHandler_CreateContract_Overlap(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	// Create existing contract
	db.Create(&models.ChildContract{
		ChildID:    child.ID,
		Period:     models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), To: nil},
		Attributes: []string{"ganztags"},
	})

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	// Try to create overlapping contract
	body := models.ChildContractCreateRequest{
		From:       time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
		To:         nil,
		Attributes: []string{"teilzeit"},
	}

	w := performRequest(r, "POST", "/organizations/1/children/1/contracts", body)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d: %s", http.StatusConflict, w.Code, w.Body.String())
	}
}

func TestChildHandler_DeleteContract(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	contract := &models.ChildContract{
		ChildID: child.ID,
		Period:  models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	db.Create(contract)

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id/contracts/:contractId", handler.DeleteContract)

	w := performRequest(r, "DELETE", "/organizations/1/children/1/contracts/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

// Edge case tests

func TestChildHandler_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	_ = org

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/1/children/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_Get_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/1/children/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_Get_ZeroID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/1/children/0", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d for zero ID, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_Create_MissingRequiredFields(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := map[string]interface{}{}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for missing fields, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_Create_EmptyFirstName(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := models.ChildCreateRequest{
		FirstName: "",
		LastName:  "Child",
		Birthdate: time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for empty first name, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_Create_EmptyLastName(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := models.ChildCreateRequest{
		FirstName: "Test",
		LastName:  "",
		Birthdate: time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for empty last name, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.PUT("/organizations/:orgId/children/:id", handler.Update)

	newName := "Updated"
	body := models.ChildUpdateRequest{
		FirstName: &newName,
	}

	w := performRequest(r, "PUT", "/organizations/1/children/999", body)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_Update_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.PUT("/organizations/:orgId/children/:id", handler.Update)

	newName := "Updated"
	body := models.ChildUpdateRequest{
		FirstName: &newName,
	}

	w := performRequest(r, "PUT", "/organizations/1/children/invalid", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/organizations/1/children/999", nil)

	if w.Code != http.StatusNoContent && w.Code != http.StatusNotFound {
		t.Errorf("expected status %d or %d, got %d", http.StatusNoContent, http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_Delete_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/organizations/1/children/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	w := performRequest(r, "GET", "/organizations/1/children", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if len(response.Data) != 0 {
		t.Errorf("expected empty list, got %d children", len(response.Data))
	}
}

func TestChildHandler_ListContracts_ChildNotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts", handler.ListContracts)

	w := performRequest(r, "GET", "/organizations/1/children/999/contracts", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_ListContracts_Empty(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts", handler.ListContracts)

	w := performRequest(r, "GET", "/organizations/1/children/1/contracts", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var contracts []models.ChildContract
	parseResponse(t, w, &contracts)

	if len(contracts) != 0 {
		t.Errorf("expected empty list, got %d contracts", len(contracts))
	}
}

func TestChildHandler_CreateContract_ChildNotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	body := models.ChildContractCreateRequest{
		From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children/999/contracts", body)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_CreateContract_InvalidChildID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	body := models.ChildContractCreateRequest{
		From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children/invalid/contracts", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_CreateContract_ContractBoundaryTouch(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	// Create contract ending on specific date
	endDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	db.Create(&models.ChildContract{
		ChildID:    child.ID,
		Period:     models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), To: &endDate},
		Attributes: []string{"ganztags"},
	})

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	// Create contract starting the day after
	body := models.ChildContractCreateRequest{
		From:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		Attributes: []string{"teilzeit"},
	}

	w := performRequest(r, "POST", "/organizations/1/children/1/contracts", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d for non-overlapping contract, got %d: %s",
			http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestChildHandler_DeleteContract_NotFound(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id/contracts/:contractId", handler.DeleteContract)

	w := performRequest(r, "DELETE", "/organizations/1/children/1/contracts/999", nil)

	if w.Code != http.StatusNoContent && w.Code != http.StatusNotFound {
		t.Errorf("expected status %d or %d, got %d", http.StatusNoContent, http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_DeleteContract_InvalidContractID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	db.Create(&models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	})

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id/contracts/:contractId", handler.DeleteContract)

	w := performRequest(r, "DELETE", "/organizations/1/children/1/contracts/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_GetCurrentContract_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts/current", handler.GetCurrentContract)

	w := performRequest(r, "GET", "/organizations/1/children/invalid/contracts/current", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// Validation edge case tests

func TestChildHandler_Create_FutureBirthdate(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := models.ChildCreateRequest{
		FirstName: "Test",
		LastName:  "Child",
		Birthdate: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for future birthdate, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestChildHandler_Create_WhitespaceOnlyFirstName(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := models.ChildCreateRequest{
		FirstName: "   ",
		LastName:  "Child",
		Birthdate: time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for whitespace-only first name, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestChildHandler_Create_WhitespaceOnlyLastName(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children", handler.Create)

	body := models.ChildCreateRequest{
		FirstName: "Test",
		LastName:  "   ",
		Birthdate: time.Date(2020, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	w := performRequest(r, "POST", "/organizations/1/children", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for whitespace-only last name, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestChildHandler_CreateContract_FromAfterTo(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")
	child := &models.Child{
		Person: models.Person{OrganizationID: org.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	toDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	body := models.ChildContractCreateRequest{
		From: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		To:   &toDate,
	}

	w := performRequest(r, "POST", "/organizations/1/children/1/contracts", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for from > to, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

// SECURITY TESTS: Cross-organization access attempts

func TestChildHandler_Get_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id", handler.Get)

	// Try to access org1's child via org2's URL
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children/%d", org2.ID, child.ID), nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when accessing child from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_Update_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Original", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.PUT("/organizations/:orgId/children/:id", handler.Update)

	newName := "Hacker"
	body := models.ChildUpdateRequest{
		FirstName: &newName,
	}

	// Try to update org1's child via org2's URL
	w := performRequest(r, "PUT", fmt.Sprintf("/organizations/%d/children/%d", org2.ID, child.ID), body)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when updating child from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_Delete_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id", handler.Delete)

	// Try to delete org1's child via org2's URL
	w := performRequest(r, "DELETE", fmt.Sprintf("/organizations/%d/children/%d", org2.ID, child.ID), nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when deleting child from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_ListContracts_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts", handler.ListContracts)

	// Try to list contracts for org1's child via org2's URL
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children/%d/contracts", org2.ID, child.ID), nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when listing contracts from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_CreateContract_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	r := setupTestRouter()
	r.POST("/organizations/:orgId/children/:id/contracts", handler.CreateContract)

	body := models.ChildContractCreateRequest{
		From: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Try to create contract for org1's child via org2's URL
	w := performRequest(r, "POST", fmt.Sprintf("/organizations/%d/children/%d/contracts", org2.ID, child.ID), body)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when creating contract from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_DeleteContract_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	// Create contract for child
	contract := &models.ChildContract{
		ChildID: child.ID,
		Period:  models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	db.Create(contract)

	r := setupTestRouter()
	r.DELETE("/organizations/:orgId/children/:id/contracts/:contractId", handler.DeleteContract)

	// Try to delete contract for org1's child via org2's URL
	w := performRequest(r, "DELETE", fmt.Sprintf("/organizations/%d/children/%d/contracts/%d", org2.ID, child.ID, contract.ID), nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when deleting contract from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

func TestChildHandler_GetCurrentContract_WrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Create child in org1
	child := &models.Child{
		Person: models.Person{OrganizationID: org1.ID, FirstName: "Test", LastName: "Child", Birthdate: time.Now()},
	}
	db.Create(child)

	// Create active contract
	db.Create(&models.ChildContract{
		ChildID: child.ID,
		Period:  models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	})

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children/:id/contracts/current", handler.GetCurrentContract)

	// Try to get current contract for org1's child via org2's URL
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children/%d/contracts/current", org2.ID, child.ID), nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("SECURITY: expected status %d when getting current contract from wrong org, got %d", http.StatusNotFound, w.Code)
	}
}

// PAGINATION TESTS

func TestChildHandler_List_Pagination(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	// Create 15 children
	for i := 1; i <= 15; i++ {
		db.Create(&models.Child{
			Person: models.Person{
				OrganizationID: org.ID,
				FirstName:      fmt.Sprintf("Child%02d", i),
				LastName:       "Last",
				Birthdate:      time.Now(),
			},
		})
	}

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	// Test page 1 with limit 5
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=1&limit=5", org.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if len(response.Data) != 5 {
		t.Errorf("expected 5 children on page 1, got %d", len(response.Data))
	}
	if response.Page != 1 {
		t.Errorf("expected page 1, got %d", response.Page)
	}
	if response.Limit != 5 {
		t.Errorf("expected limit 5, got %d", response.Limit)
	}
	if response.Total != 15 {
		t.Errorf("expected total 15, got %d", response.Total)
	}
	if response.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", response.TotalPages)
	}
}

func TestChildHandler_List_Pagination_SecondPage(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	// Create 15 children
	for i := 1; i <= 15; i++ {
		db.Create(&models.Child{
			Person: models.Person{
				OrganizationID: org.ID,
				FirstName:      fmt.Sprintf("Child%02d", i),
				LastName:       "Last",
				Birthdate:      time.Now(),
			},
		})
	}

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	// Get page 1
	w1 := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=1&limit=5", org.ID), nil)
	var response1 models.PaginatedResponse[models.Child]
	parseResponse(t, w1, &response1)

	// Get page 2
	w2 := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=2&limit=5", org.ID), nil)

	if w2.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w2.Code)
	}

	var response2 models.PaginatedResponse[models.Child]
	parseResponse(t, w2, &response2)

	if len(response2.Data) != 5 {
		t.Errorf("expected 5 children on page 2, got %d", len(response2.Data))
	}
	if response2.Page != 2 {
		t.Errorf("expected page 2, got %d", response2.Page)
	}

	// Verify page 1 and page 2 have different children
	page1IDs := make(map[uint]bool)
	for _, child := range response1.Data {
		page1IDs[child.ID] = true
	}
	for _, child := range response2.Data {
		if page1IDs[child.ID] {
			t.Errorf("child ID %d appears on both page 1 and page 2", child.ID)
		}
	}
}

func TestChildHandler_List_Pagination_LastPage(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	// Create 12 children (3 pages of 5, last page has 2)
	for i := 1; i <= 12; i++ {
		db.Create(&models.Child{
			Person: models.Person{
				OrganizationID: org.ID,
				FirstName:      fmt.Sprintf("Child%02d", i),
				LastName:       "Last",
				Birthdate:      time.Now(),
			},
		})
	}

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	// Get last page (page 3)
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=3&limit=5", org.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if len(response.Data) != 2 {
		t.Errorf("expected 2 children on last page, got %d", len(response.Data))
	}
	if response.TotalPages != 3 {
		t.Errorf("expected 3 total pages, got %d", response.TotalPages)
	}
}

func TestChildHandler_List_Pagination_BeyondLastPage(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	// Create 5 children
	for i := 1; i <= 5; i++ {
		db.Create(&models.Child{
			Person: models.Person{
				OrganizationID: org.ID,
				FirstName:      fmt.Sprintf("Child%02d", i),
				LastName:       "Last",
				Birthdate:      time.Now(),
			},
		})
	}

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	// Request page 10 (beyond available data)
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=10&limit=5", org.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if len(response.Data) != 0 {
		t.Errorf("expected 0 children beyond last page, got %d", len(response.Data))
	}
	if response.Total != 5 {
		t.Errorf("expected total 5, got %d", response.Total)
	}
}

func TestChildHandler_List_Pagination_DefaultValues(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	// Create 25 children
	for i := 1; i <= 25; i++ {
		db.Create(&models.Child{
			Person: models.Person{
				OrganizationID: org.ID,
				FirstName:      fmt.Sprintf("Child%02d", i),
				LastName:       "Last",
				Birthdate:      time.Now(),
			},
		})
	}

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	// Request without pagination params (should use defaults: page=1, limit=20)
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children", org.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if len(response.Data) != 20 {
		t.Errorf("expected 20 children with default limit, got %d", len(response.Data))
	}
	if response.Page != 1 {
		t.Errorf("expected default page 1, got %d", response.Page)
	}
	if response.Limit != 20 {
		t.Errorf("expected default limit 20, got %d", response.Limit)
	}
}

func TestChildHandler_List_Pagination_InvalidNegativePage(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=-1&limit=10", org.ID), nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for negative page, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_List_Pagination_InvalidNegativeLimit(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=1&limit=-5", org.ID), nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for negative limit, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_List_Pagination_LimitExceedsMax(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=1&limit=101", org.ID), nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for limit > 100, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestChildHandler_List_Pagination_MaxLimit(t *testing.T) {
	db := setupTestDB(t)
	childService := createChildService(db)
	handler := NewChildHandler(childService)

	org := createTestOrganization(t, db, "Test Org")

	// Create 5 children
	for i := 1; i <= 5; i++ {
		db.Create(&models.Child{
			Person: models.Person{
				OrganizationID: org.ID,
				FirstName:      fmt.Sprintf("Child%02d", i),
				LastName:       "Last",
				Birthdate:      time.Now(),
			},
		})
	}

	r := setupTestRouter()
	r.GET("/organizations/:orgId/children", handler.List)

	// Request with max limit (100)
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/children?page=1&limit=100", org.ID), nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d for limit=100, got %d", http.StatusOK, w.Code)
	}

	var response models.PaginatedResponse[models.Child]
	parseResponse(t, w, &response)

	if response.Limit != 100 {
		t.Errorf("expected limit 100, got %d", response.Limit)
	}
}
