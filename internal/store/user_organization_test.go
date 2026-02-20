package store

import (
	"errors"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestUserOrganizationStore_AddUserToOrg(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")

	uo, err := store.AddUserToOrg(ctx, user.ID, org.ID, models.RoleMember, "creator@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if uo.UserID != user.ID {
		t.Errorf("UserID = %d, want %d", uo.UserID, user.ID)
	}
	if uo.OrganizationID != org.ID {
		t.Errorf("OrganizationID = %d, want %d", uo.OrganizationID, org.ID)
	}
	if uo.Role != models.RoleMember {
		t.Errorf("Role = %v, want %v", uo.Role, models.RoleMember)
	}
	if uo.CreatedBy != "creator@example.com" {
		t.Errorf("CreatedBy = %v, want creator@example.com", uo.CreatedBy)
	}
}

func TestUserOrganizationStore_AddUserToOrg_AllRoles(t *testing.T) {
	roles := []models.Role{models.RoleAdmin, models.RoleManager, models.RoleMember}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			db := setupTestDB(t)
			store := NewUserOrganizationStore(db)

			user := createTestUser(t, db, "Test User", "test@example.com")
			org := createTestOrganization(t, db, "Test Org")

			uo, err := store.AddUserToOrg(ctx, user.ID, org.ID, role, "test@example.com")
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if uo.Role != role {
				t.Errorf("Role = %v, want %v", uo.Role, role)
			}
		})
	}
}

func TestUserOrganizationStore_AddUserToOrg_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")

	// First addition should succeed
	_, err := store.AddUserToOrg(ctx, user.ID, org.ID, models.RoleMember, "test@example.com")
	if err != nil {
		t.Fatalf("first add: expected no error, got %v", err)
	}

	// Second addition should fail due to composite primary key constraint
	_, err = store.AddUserToOrg(ctx, user.ID, org.ID, models.RoleAdmin, "test@example.com")
	if err == nil {
		t.Error("expected error for duplicate user-org, got nil")
	}
}

func TestUserOrganizationStore_UpdateRole(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	err := store.UpdateRole(ctx, user.ID, org.ID, models.RoleAdmin)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the role was updated
	uo, _ := store.FindByUserAndOrg(ctx, user.ID, org.ID)
	if uo.Role != models.RoleAdmin {
		t.Errorf("Role = %v after update, want %v", uo.Role, models.RoleAdmin)
	}
}

func TestUserOrganizationStore_UpdateRole_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	err := store.UpdateRole(ctx, 999, 999, models.RoleAdmin)
	if err == nil {
		t.Error("expected error for non-existent relationship, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationStore_RemoveUserFromOrg(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	err := store.RemoveUserFromOrg(ctx, user.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the relationship was removed
	exists, _ := store.Exists(ctx, user.ID, org.ID)
	if exists {
		t.Error("expected relationship to be removed")
	}
}

func TestUserOrganizationStore_RemoveUserFromOrg_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	// Removing non-existent relationship should not error (idempotent)
	err := store.RemoveUserFromOrg(ctx, 999, 999)
	if err != nil {
		t.Errorf("expected no error for non-existent relationship, got %v", err)
	}
}

func TestUserOrganizationStore_RemoveUserFromOrg_UserStillExists(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)
	userStore := NewUserStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	_ = store.RemoveUserFromOrg(ctx, user.ID, org.ID)

	// Verify user still exists
	found, err := userStore.FindByID(ctx, user.ID)
	if err != nil {
		t.Errorf("expected user to still exist, got error: %v", err)
	}
	if found.ID != user.ID {
		t.Error("user should still exist after removing from organization")
	}
}

func TestUserOrganizationStore_RemoveUserFromOrg_PreservesOtherOrgs(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	createTestUserOrganization(t, db, user.ID, org1.ID, models.RoleAdmin)
	createTestUserOrganization(t, db, user.ID, org2.ID, models.RoleMember)

	// Remove from org1 only
	err := store.RemoveUserFromOrg(ctx, user.ID, org1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify org1 membership is gone
	exists1, _ := store.Exists(ctx, user.ID, org1.ID)
	if exists1 {
		t.Error("expected org1 membership to be removed")
	}

	// Verify org2 membership is preserved
	exists2, _ := store.Exists(ctx, user.ID, org2.ID)
	if !exists2 {
		t.Error("expected org2 membership to be preserved")
	}
}

func TestUserOrganizationStore_FindByUserAndOrg(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleAdmin)

	found, err := store.FindByUserAndOrg(ctx, user.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.UserID != user.ID {
		t.Errorf("UserID = %d, want %d", found.UserID, user.ID)
	}
	if found.OrganizationID != org.ID {
		t.Errorf("OrganizationID = %d, want %d", found.OrganizationID, org.ID)
	}
	if found.Role != models.RoleAdmin {
		t.Errorf("Role = %v, want %v", found.Role, models.RoleAdmin)
	}
}

func TestUserOrganizationStore_FindByUserAndOrg_NotFound(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	_, err := store.FindByUserAndOrg(ctx, 999, 999)
	if err == nil {
		t.Error("expected error for non-existent relationship")
	}
}

func TestUserOrganizationStore_FindByUser(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	createTestUserOrganization(t, db, user.ID, org1.ID, models.RoleAdmin)
	createTestUserOrganization(t, db, user.ID, org2.ID, models.RoleMember)

	memberships, err := store.FindByUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(memberships) != 2 {
		t.Fatalf("expected 2 memberships, got %d", len(memberships))
	}
}

func TestUserOrganizationStore_FindByUser_NoOrgs(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")

	memberships, err := store.FindByUser(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(memberships) != 0 {
		t.Errorf("expected 0 memberships, got %d", len(memberships))
	}
}

func TestUserOrganizationStore_FindByUser_PreloadsOrganization(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")

	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	memberships, _ := store.FindByUser(ctx, user.ID)

	if len(memberships) != 1 {
		t.Fatalf("expected 1 membership, got %d", len(memberships))
	}
	if memberships[0].Organization == nil {
		t.Fatal("expected Organization to be preloaded")
	}
	if memberships[0].Organization.Name != "Test Org" {
		t.Errorf("Organization.Name = %v, want Test Org", memberships[0].Organization.Name)
	}
}

func TestUserOrganizationStore_GetRoleInOrg(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleAdmin)

	role, err := store.GetRoleInOrg(ctx, user.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if role != models.RoleAdmin {
		t.Errorf("role = %v, want %v", role, models.RoleAdmin)
	}
}

func TestUserOrganizationStore_GetRoleInOrg_NotInOrg(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")

	role, err := store.GetRoleInOrg(ctx, user.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if role != "" {
		t.Errorf("role = %v, want empty string", role)
	}
}

func TestUserOrganizationStore_GetRoleInOrg_MultipleOrgs(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Admin in org1, member in org2
	createTestUserOrganization(t, db, user.ID, org1.ID, models.RoleAdmin)
	createTestUserOrganization(t, db, user.ID, org2.ID, models.RoleMember)

	// Check org1
	role1, _ := store.GetRoleInOrg(ctx, user.ID, org1.ID)
	if role1 != models.RoleAdmin {
		t.Errorf("org1 role = %v, want %v", role1, models.RoleAdmin)
	}

	// Check org2
	role2, _ := store.GetRoleInOrg(ctx, user.ID, org2.ID)
	if role2 != models.RoleMember {
		t.Errorf("org2 role = %v, want %v", role2, models.RoleMember)
	}
}

func TestUserOrganizationStore_GetUserOrganizationsWithRoles(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	createTestUserOrganization(t, db, user.ID, org1.ID, models.RoleAdmin)
	createTestUserOrganization(t, db, user.ID, org2.ID, models.RoleMember)

	orgRoles, err := store.GetUserOrganizationsWithRoles(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(orgRoles) != 2 {
		t.Fatalf("expected 2 orgs, got %d", len(orgRoles))
	}
	if orgRoles[org1.ID] != models.RoleAdmin {
		t.Errorf("org1 role = %v, want %v", orgRoles[org1.ID], models.RoleAdmin)
	}
	if orgRoles[org2.ID] != models.RoleMember {
		t.Errorf("org2 role = %v, want %v", orgRoles[org2.ID], models.RoleMember)
	}
}

func TestUserOrganizationStore_GetUserOrganizationsWithRoles_NoOrgs(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")

	orgRoles, err := store.GetUserOrganizationsWithRoles(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(orgRoles) != 0 {
		t.Errorf("expected empty map, got %v", orgRoles)
	}
}

func TestUserOrganizationStore_SetSuperAdmin(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")

	// Set to true
	err := store.SetSuperAdmin(ctx, user.ID, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	isSuperAdmin, _ := store.IsSuperAdmin(ctx, user.ID)
	if !isSuperAdmin {
		t.Error("expected IsSuperAdmin = true")
	}

	// Set to false
	err = store.SetSuperAdmin(ctx, user.ID, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	isSuperAdmin, _ = store.IsSuperAdmin(ctx, user.ID)
	if isSuperAdmin {
		t.Error("expected IsSuperAdmin = false")
	}
}

func TestUserOrganizationStore_SetSuperAdmin_Toggle(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")

	// true -> false -> true
	_ = store.SetSuperAdmin(ctx, user.ID, true)
	is1, _ := store.IsSuperAdmin(ctx, user.ID)
	if !is1 {
		t.Error("expected true after first set")
	}

	_ = store.SetSuperAdmin(ctx, user.ID, false)
	is2, _ := store.IsSuperAdmin(ctx, user.ID)
	if is2 {
		t.Error("expected false after second set")
	}

	_ = store.SetSuperAdmin(ctx, user.ID, true)
	is3, _ := store.IsSuperAdmin(ctx, user.ID)
	if !is3 {
		t.Error("expected true after third set")
	}
}

func TestUserOrganizationStore_SetSuperAdmin_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	err := store.SetSuperAdmin(ctx, 999, true)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationStore_IsSuperAdmin(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")

	// Initially should be false
	isSuperAdmin, err := store.IsSuperAdmin(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if isSuperAdmin {
		t.Error("expected IsSuperAdmin = false initially")
	}
}

func TestUserOrganizationStore_IsSuperAdmin_NonExistent(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	isSuperAdmin, err := store.IsSuperAdmin(ctx, 999)
	if err != nil {
		t.Errorf("expected no error for non-existent user, got %v", err)
	}
	if isSuperAdmin {
		t.Error("expected false for non-existent user")
	}
}

func TestUserOrganizationStore_Exists(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	user := createTestUser(t, db, "Test User", "test@example.com")
	org := createTestOrganization(t, db, "Test Org")

	// Should not exist initially
	exists, err := store.Exists(ctx, user.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected Exists = false initially")
	}

	// Create relationship
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	// Should exist now
	exists, err = store.Exists(ctx, user.ID, org.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !exists {
		t.Error("expected Exists = true after creation")
	}
}

func TestUserOrganizationStore_Exists_NonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	store := NewUserOrganizationStore(db)

	exists, err := store.Exists(ctx, 999, 999)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exists {
		t.Error("expected Exists = false for non-existent user/org")
	}
}
