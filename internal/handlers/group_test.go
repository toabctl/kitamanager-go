package handlers

import (
	"net/http"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestGroupHandler_List(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	createTestGroup(t, db, "Group 1")
	createTestGroup(t, db, "Group 2")

	r := setupTestRouter()
	r.GET("/groups", handler.List)

	w := performRequest(r, "GET", "/groups", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var groups []models.Group
	parseResponse(t, w, &groups)

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}
}

func TestGroupHandler_Get(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	group := createTestGroup(t, db, "Test Group")

	r := setupTestRouter()
	r.GET("/groups/:id", handler.Get)

	w := performRequest(r, "GET", "/groups/1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result models.Group
	parseResponse(t, w, &result)

	if result.Name != group.Name {
		t.Errorf("expected name '%s', got '%s'", group.Name, result.Name)
	}
}

func TestGroupHandler_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.GET("/groups/:id", handler.Get)

	w := performRequest(r, "GET", "/groups/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGroupHandler_Create(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	// Create an organization first
	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/groups", handler.Create)

	body := CreateGroupRequest{
		Name:           "New Group",
		OrganizationID: org.ID,
		Active:         true,
	}

	w := performRequest(r, "POST", "/groups", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var result models.GroupResponse
	parseResponse(t, w, &result)

	if result.Name != "New Group" {
		t.Errorf("expected name 'New Group', got '%s'", result.Name)
	}
	if result.CreatedBy != "test@example.com" {
		t.Errorf("expected created_by 'test@example.com', got '%s'", result.CreatedBy)
	}
	if result.OrganizationID != org.ID {
		t.Errorf("expected organization_id %d, got %d", org.ID, result.OrganizationID)
	}
}

func TestGroupHandler_Create_BadRequest(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.POST("/groups", handler.Create)

	// Missing required name field
	body := map[string]interface{}{
		"active": true,
	}

	w := performRequest(r, "POST", "/groups", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGroupHandler_Update(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	createTestGroup(t, db, "Original Name")

	r := setupTestRouter()
	r.PUT("/groups/:id", handler.Update)

	body := UpdateGroupRequest{
		Name: "Updated Name",
	}

	w := performRequest(r, "PUT", "/groups/1", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var result models.Group
	parseResponse(t, w, &result)

	if result.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", result.Name)
	}
}

func TestGroupHandler_Delete(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	createTestGroup(t, db, "To Delete")

	r := setupTestRouter()
	r.DELETE("/groups/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/groups/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify group was deleted
	groups, _ := groupStore.FindAll()
	if len(groups) != 0 {
		t.Error("expected group to be deleted")
	}
}

func TestGroupHandler_Create_BadRequest_MissingOrganization(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.POST("/groups", handler.Create)

	// Missing required organization_id field
	body := map[string]interface{}{
		"name":   "Test Group",
		"active": true,
	}

	w := performRequest(r, "POST", "/groups", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGroupHandler_Get_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.GET("/groups/:id", handler.Get)

	w := performRequest(r, "GET", "/groups/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGroupHandler_Update_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.PUT("/groups/:id", handler.Update)

	body := UpdateGroupRequest{
		Name: "Updated Name",
	}

	w := performRequest(r, "PUT", "/groups/invalid", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGroupHandler_Update_NotFound(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.PUT("/groups/:id", handler.Update)

	body := UpdateGroupRequest{
		Name: "Updated Name",
	}

	w := performRequest(r, "PUT", "/groups/999", body)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGroupHandler_Update_ActiveFlag(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	createTestGroup(t, db, "Test Group")

	r := setupTestRouter()
	r.PUT("/groups/:id", handler.Update)

	active := false
	body := UpdateGroupRequest{
		Active: &active,
	}

	w := performRequest(r, "PUT", "/groups/1", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var result models.Group
	parseResponse(t, w, &result)

	if result.Active != false {
		t.Errorf("expected active false, got %v", result.Active)
	}
}

func TestGroupHandler_Delete_InvalidID(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.DELETE("/groups/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/groups/invalid", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGroupHandler_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	r := setupTestRouter()
	r.GET("/groups", handler.List)

	w := performRequest(r, "GET", "/groups", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var groups []models.GroupResponse
	parseResponse(t, w, &groups)

	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}
