package service

import (
	"context"
	"errors"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestUserOrganizationService_AddUserToOrganization(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")

	resp, err := svc.AddUserToOrganization(ctx, user.ID, org.ID, models.RoleMember, "creator@example.com", admin.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.UserID != user.ID {
		t.Errorf("UserID = %d, want %d", resp.UserID, user.ID)
	}
	if resp.OrganizationID != org.ID {
		t.Errorf("OrganizationID = %d, want %d", resp.OrganizationID, org.ID)
	}
	if resp.Role != models.RoleMember {
		t.Errorf("Role = %v, want %v", resp.Role, models.RoleMember)
	}
	if resp.CreatedBy != "creator@example.com" {
		t.Errorf("CreatedBy = %v, want creator@example.com", resp.CreatedBy)
	}
}

func TestUserOrganizationService_AddUserToOrganization_AllRoles(t *testing.T) {
	roles := []models.Role{models.RoleAdmin, models.RoleManager, models.RoleMember}

	for _, role := range roles {
		t.Run(string(role), func(t *testing.T) {
			db := setupTestDB(t)
			svc := createUserOrganizationService(db)
			ctx := context.Background()

			admin := createTestSuperAdmin(t, db)
			user := createTestUser(t, db, "Test User", "test@example.com", "password")
			org := createTestOrganization(t, db, "Test Org")

			resp, err := svc.AddUserToOrganization(ctx, user.ID, org.ID, role, "test@example.com", admin.ID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if resp.Role != role {
				t.Errorf("Role = %v, want %v", resp.Role, role)
			}
		})
	}
}

func TestUserOrganizationService_AddUserToOrganization_InvalidRole(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")

	_, err := svc.AddUserToOrganization(ctx, user.ID, org.ID, models.Role("invalid"), "test@example.com", admin.ID)
	if err == nil {
		t.Fatal("expected error for invalid role, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestUserOrganizationService_AddUserToOrganization_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	org := createTestOrganization(t, db, "Test Org")

	_, err := svc.AddUserToOrganization(ctx, 999, org.ID, models.RoleMember, "test@example.com", admin.ID)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationService_AddUserToOrganization_AlreadyMember(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")

	// First add
	_, err := svc.AddUserToOrganization(ctx, user.ID, org.ID, models.RoleMember, "test@example.com", admin.ID)
	if err != nil {
		t.Fatalf("first add: expected no error, got %v", err)
	}

	// Second add should fail
	_, err = svc.AddUserToOrganization(ctx, user.ID, org.ID, models.RoleAdmin, "test@example.com", admin.ID)
	if err == nil {
		t.Fatal("expected error for already member, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestUserOrganizationService_UpdateUserOrganizationRole(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	resp, err := svc.UpdateUserOrganizationRole(ctx, user.ID, org.ID, models.RoleAdmin, admin.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Role != models.RoleAdmin {
		t.Errorf("Role = %v, want %v", resp.Role, models.RoleAdmin)
	}
}

func TestUserOrganizationService_UpdateUserOrganizationRole_InvalidRole(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	_, err := svc.UpdateUserOrganizationRole(ctx, user.ID, org.ID, models.Role("invalid"), admin.ID)
	if err == nil {
		t.Fatal("expected error for invalid role, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestUserOrganizationService_UpdateUserOrganizationRole_NotMember(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")
	// No membership created

	_, err := svc.UpdateUserOrganizationRole(ctx, user.ID, org.ID, models.RoleAdmin, admin.ID)
	if err == nil {
		t.Fatal("expected error for non-member, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationService_RemoveUserFromOrganization(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")
	createTestUserOrganization(t, db, user.ID, org.ID, models.RoleMember)

	err := svc.RemoveUserFromOrganization(ctx, user.ID, org.ID, admin.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify user is no longer in org
	resp, _ := svc.GetUserMemberships(ctx, user.ID)
	if len(resp.Memberships) != 0 {
		t.Errorf("expected 0 memberships after removal, got %d", len(resp.Memberships))
	}
}

func TestUserOrganizationService_RemoveUserFromOrganization_NotMember(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org := createTestOrganization(t, db, "Test Org")
	// No membership

	err := svc.RemoveUserFromOrganization(ctx, user.ID, org.ID, admin.ID)
	if err == nil {
		t.Fatal("expected error for non-member, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationService_RemoveUserFromOrganization_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	admin := createTestSuperAdmin(t, db)
	org := createTestOrganization(t, db, "Test Org")

	err := svc.RemoveUserFromOrganization(ctx, 999, org.ID, admin.ID)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationService_GetUserMemberships(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	createTestUserOrganization(t, db, user.ID, org1.ID, models.RoleAdmin)
	createTestUserOrganization(t, db, user.ID, org2.ID, models.RoleMember)

	resp, err := svc.GetUserMemberships(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Memberships) != 2 {
		t.Fatalf("expected 2 memberships, got %d", len(resp.Memberships))
	}
}

func TestUserOrganizationService_GetUserMemberships_NoMemberships(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	user := createTestUser(t, db, "Test User", "test@example.com", "password")

	resp, err := svc.GetUserMemberships(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Memberships) != 0 {
		t.Errorf("expected 0 memberships, got %d", len(resp.Memberships))
	}
}

func TestUserOrganizationService_GetUserMemberships_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	_, err := svc.GetUserMemberships(ctx, 999)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestUserOrganizationService_GetUserMemberships_MultipleOrgs(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	user := createTestUser(t, db, "Test User", "test@example.com", "password")
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")

	// Admin in org1, member in org2
	createTestUserOrganization(t, db, user.ID, org1.ID, models.RoleAdmin)
	createTestUserOrganization(t, db, user.ID, org2.ID, models.RoleMember)

	resp, err := svc.GetUserMemberships(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Memberships) != 2 {
		t.Fatalf("expected 2 memberships, got %d", len(resp.Memberships))
	}

	// Check roles are correct per org
	for _, m := range resp.Memberships {
		if m.OrganizationID == org1.ID && m.Role != models.RoleAdmin {
			t.Errorf("org1 Role = %v, want %v", m.Role, models.RoleAdmin)
		}
		if m.OrganizationID == org2.ID && m.Role != models.RoleMember {
			t.Errorf("org2 Role = %v, want %v", m.Role, models.RoleMember)
		}
	}
}

func TestUserOrganizationService_SetSuperAdmin(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	user := createTestUser(t, db, "Test User", "test@example.com", "password")

	err := svc.SetSuperAdmin(ctx, user.ID, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify it was set
	var dbUser models.User
	db.First(&dbUser, user.ID)
	if !dbUser.IsSuperAdmin {
		t.Error("expected IsSuperAdmin = true")
	}
}

func TestUserOrganizationService_SetSuperAdmin_Toggle(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	user := createTestUser(t, db, "Test User", "test@example.com", "password")

	// Set true
	_ = svc.SetSuperAdmin(ctx, user.ID, true)
	var dbUser models.User
	db.First(&dbUser, user.ID)
	if !dbUser.IsSuperAdmin {
		t.Error("expected IsSuperAdmin = true after set true")
	}

	// Set false
	_ = svc.SetSuperAdmin(ctx, user.ID, false)
	db.First(&dbUser, user.ID)
	if dbUser.IsSuperAdmin {
		t.Error("expected IsSuperAdmin = false after set false")
	}
}

func TestUserOrganizationService_SetSuperAdmin_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	svc := createUserOrganizationService(db)
	ctx := context.Background()

	err := svc.SetSuperAdmin(ctx, 999, true)
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
