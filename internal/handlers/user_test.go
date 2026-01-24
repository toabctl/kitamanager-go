package handlers

import (
	"net/http"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestUserHandler_List(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "User 1", "user1@example.com", "password")
	createTestUser(t, db, "User 2", "user2@example.com", "password")

	r := setupTestRouter()
	r.GET("/users", handler.List)

	w := performRequest(r, "GET", "/users", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var users []models.UserResponse
	parseResponse(t, w, &users)

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestUserHandler_Get(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	user := createTestUser(t, db, "Test User", "test@example.com", "password")

	r := setupTestRouter()
	r.GET("/users/:id", handler.Get)

	w := performRequest(r, "GET", "/users/1", nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result models.UserResponse
	parseResponse(t, w, &result)

	if result.Name != user.Name {
		t.Errorf("expected name '%s', got '%s'", user.Name, result.Name)
	}
}

func TestUserHandler_Get_NotFound(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	r := setupTestRouter()
	r.GET("/users/:id", handler.Get)

	w := performRequest(r, "GET", "/users/999", nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUserHandler_Create(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	r := setupTestRouter()
	r.POST("/users", handler.Create)

	body := models.UserCreate{
		Name:     "New User",
		Email:    "new@example.com",
		Password: "password123",
		Active:   true,
	}

	w := performRequest(r, "POST", "/users", body)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var result models.UserResponse
	parseResponse(t, w, &result)

	if result.Name != "New User" {
		t.Errorf("expected name 'New User', got '%s'", result.Name)
	}
	if result.CreatedBy != "test@example.com" {
		t.Errorf("expected created_by 'test@example.com', got '%s'", result.CreatedBy)
	}
}

func TestUserHandler_Create_BadRequest(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	r := setupTestRouter()
	r.POST("/users", handler.Create)

	// Missing required fields
	body := map[string]interface{}{
		"active": true,
	}

	w := performRequest(r, "POST", "/users", body)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUserHandler_Update(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "Original Name", "test@example.com", "password")

	r := setupTestRouter()
	r.PUT("/users/:id", handler.Update)

	body := models.UserUpdate{
		Name: "Updated Name",
	}

	w := performRequest(r, "PUT", "/users/1", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var result models.UserResponse
	parseResponse(t, w, &result)

	if result.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", result.Name)
	}
}

func TestUserHandler_Delete(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "To Delete", "delete@example.com", "password")

	r := setupTestRouter()
	r.DELETE("/users/:id", handler.Delete)

	w := performRequest(r, "DELETE", "/users/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify user was deleted
	users, _ := userStore.FindAll()
	if len(users) != 0 {
		t.Error("expected user to be deleted")
	}
}

func TestUserHandler_AddToGroup(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "Test User", "test@example.com", "password")
	group := createTestGroup(t, db, "Test Group")

	r := setupTestRouter()
	r.POST("/users/:id/groups", handler.AddToGroup)

	body := AddToGroupRequest{
		GroupID: group.ID,
	}

	w := performRequest(r, "POST", "/users/1/groups", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify user was added to group
	user, _ := userStore.FindByID(1)
	if len(user.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(user.Groups))
	}
}

func TestUserHandler_RemoveFromGroup(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "Test User", "test@example.com", "password")
	group := createTestGroup(t, db, "Test Group")
	_ = userStore.AddToGroup(1, group.ID)

	r := setupTestRouter()
	r.DELETE("/users/:id/groups/:gid", handler.RemoveFromGroup)

	w := performRequest(r, "DELETE", "/users/1/groups/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d: %s", http.StatusNoContent, w.Code, w.Body.String())
	}

	// Verify user was removed from group
	user, _ := userStore.FindByID(1)
	if len(user.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(user.Groups))
	}
}

func TestUserHandler_AddToOrganization(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")

	r := setupTestRouter()
	r.POST("/users/:id/organizations", handler.AddToOrganization)

	body := AddToOrganizationRequest{
		OrganizationID: org.ID,
	}

	w := performRequest(r, "POST", "/users/1/organizations", body)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify user was added to organization
	user, _ := userStore.FindByID(1)
	if len(user.Organizations) != 1 {
		t.Errorf("expected 1 organization, got %d", len(user.Organizations))
	}
}

func TestUserHandler_RemoveFromOrganization(t *testing.T) {
	db := setupTestDB(t)
	userStore := store.NewUserStore(db)
	handler := NewUserHandler(userStore)

	createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")
	_ = userStore.AddToOrganization(1, org.ID)

	r := setupTestRouter()
	r.DELETE("/users/:id/organizations/:oid", handler.RemoveFromOrganization)

	w := performRequest(r, "DELETE", "/users/1/organizations/1", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d: %s", http.StatusNoContent, w.Code, w.Body.String())
	}

	// Verify user was removed from organization
	user, _ := userStore.FindByID(1)
	if len(user.Organizations) != 0 {
		t.Errorf("expected 0 organizations, got %d", len(user.Organizations))
	}
}
