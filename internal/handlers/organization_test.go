package handlers

import (
	"net/http"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestOrganizationHandler_List(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	createTestOrganization(t, db, "Org 1")
	createTestOrganization(t, db, "Org 2")

	r := setupTestRouter()
	r.GET("/organizations", handler.List)

	w := performRequest(r, "GET", "/organizations", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var orgs []models.Organization
	parseResponse(t, w, &orgs)

	if len(orgs) != 2 {
		t.Errorf("expected 2 organizations, got %d", len(orgs))
	}
}

func TestOrganizationHandler_Get(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.GET("/organizations/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result models.Organization
	parseResponse(t, w, &result)

	if result.Name != org.Name {
		t.Errorf("expected name '%s', got '%s'", org.Name, result.Name)
	}
}

func TestOrganizationHandler_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	r := setupTestRouter()
	r.GET("/organizations/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestOrganizationHandler_Get_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	r := setupTestRouter()
	r.GET("/organizations/:id", handler.Get)

	w := performRequest(r, "GET", "/organizations/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOrganizationHandler_Create(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	r := setupTestRouter()
	r.POST("/organizations", handler.Create)

	body := CreateOrganizationRequest{
		Name:   "New Org",
		Active: true,
	}

	w := performRequest(r, "POST", "/organizations", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var result models.Organization
	parseResponse(t, w, &result)

	if result.Name != "New Org" {
		t.Errorf("expected name 'New Org', got '%s'", result.Name)
	}
	if result.CreatedBy != "test@example.com" {
		t.Errorf("expected created_by 'test@example.com', got '%s'", result.CreatedBy)
	}
}

func TestOrganizationHandler_Create_BadRequest(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	r := setupTestRouter()
	r.POST("/organizations", handler.Create)

	// Missing required name field
	body := map[string]interface{}{
		"active": true,
	}

	w := performRequest(r, "POST", "/organizations", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOrganizationHandler_Update(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	createTestOrganization(t, db, "Original Name")

	r := setupTestRouter()
	r.PUT("/organizations/:id", handler.Update)

	body := UpdateOrganizationRequest{
		Name: "Updated Name",
	}

	w := performRequest(r, "PUT", "/organizations/1", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var result models.Organization
	parseResponse(t, w, &result)

	if result.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", result.Name)
	}
}

func TestOrganizationHandler_Delete(t *testing.T) {
	db := setupTestDB(t)
	orgStore := store.NewOrganizationStore(db)
	handler := NewOrganizationHandler(orgStore)

	createTestOrganization(t, db, "To Delete")

	r := setupTestRouter()
	r.DELETE("/organizations/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/organizations/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify organization was deleted
	orgs, _ := orgStore.FindAll()
	if len(orgs) != 0 {
		t.Error("expected organization to be deleted")
	}
}
