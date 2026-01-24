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

	r := setupTestRouter()
	r.POST("/groups", handler.Create)

	body := CreateGroupRequest{
		Name:   "New Group",
		Active: true,
	}

	w := performRequest(r, "POST", "/groups", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var result models.Group
	parseResponse(t, w, &result)

	if result.Name != "New Group" {
		t.Errorf("expected name 'New Group', got '%s'", result.Name)
	}
	if result.CreatedBy != "test@example.com" {
		t.Errorf("expected created_by 'test@example.com', got '%s'", result.CreatedBy)
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

func TestGroupHandler_AddToOrganization(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	createTestGroup(t, db, "Test Group")
	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/groups/:id/organizations", handler.AddToOrganization)

	body := AddGroupToOrganizationRequest{
		OrganizationID: org.ID,
	}

	w := performRequest(r, "POST", "/groups/1/organizations", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify group was added to organization
	group, _ := groupStore.FindByID(1)
	if len(group.Organizations) != 1 {
		t.Errorf("expected 1 organization, got %d", len(group.Organizations))
	}
}

func TestGroupHandler_RemoveFromOrganization(t *testing.T) {
	db := setupTestDB(t)
	groupStore := store.NewGroupStore(db)
	handler := NewGroupHandler(groupStore)

	createTestGroup(t, db, "Test Group")
	org := createTestOrganization(t, db, "Test Org")
	_ = groupStore.AddToOrganization(1, org.ID)

	r := setupTestRouter()
	r.DELETE("/groups/:id/organizations/:oid", handler.RemoveFromOrganization)

	w := performRequest(r, "DELETE", "/groups/1/organizations/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d: %s", http.StatusNoContent, w.Code, w.Body.String())
	}

	// Verify group was removed from organization
	group, _ := groupStore.FindByID(1)
	if len(group.Organizations) != 0 {
		t.Errorf("expected 0 organizations, got %d", len(group.Organizations))
	}
}
