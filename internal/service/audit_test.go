package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestAuditService_NewAndShutdown(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	// Shutdown should complete cleanly without panic
	svc.Shutdown()
}

func TestAuditService_LogLogin(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	var userID uint = 1
	svc.LogLogin(userID, "user@example.com", "127.0.0.1", "TestAgent/1.0")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	log := logs[0]
	if log.Action != models.AuditActionLogin {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionLogin)
	}
	if log.UserID == nil || *log.UserID != userID {
		t.Errorf("UserID = %v, want %d", log.UserID, userID)
	}
	if log.UserEmail != "user@example.com" {
		t.Errorf("UserEmail = %v, want user@example.com", log.UserEmail)
	}
	if log.IPAddress != "127.0.0.1" {
		t.Errorf("IPAddress = %v, want 127.0.0.1", log.IPAddress)
	}
	if log.UserAgent != "TestAgent/1.0" {
		t.Errorf("UserAgent = %v, want TestAgent/1.0", log.UserAgent)
	}
	if !log.Success {
		t.Error("expected Success = true")
	}
	if log.Timestamp.IsZero() {
		t.Error("expected non-zero Timestamp")
	}
}

func TestAuditService_LogLoginFailed(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogLoginFailed("bad@example.com", "10.0.0.1", "BadAgent/1.0", "invalid password")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionLoginFailed {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionLoginFailed)
	}
	if log.UserEmail != "bad@example.com" {
		t.Errorf("UserEmail = %v, want bad@example.com", log.UserEmail)
	}
	if log.Success {
		t.Error("expected Success = false")
	}

	var details map[string]string
	if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
		t.Fatalf("failed to unmarshal details: %v", err)
	}
	if details["reason"] != "invalid password" {
		t.Errorf("details[reason] = %v, want 'invalid password'", details["reason"])
	}
}

func TestAuditService_LogSuperAdminChange(t *testing.T) {
	ctx := context.Background()

	t.Run("grant", func(t *testing.T) {
		db := setupTestDB(t)
		auditStore := store.NewAuditStore(db)
		svc := NewAuditService(auditStore)

		svc.LogSuperAdminChange(1, 2, "target@example.com", true, "127.0.0.1")
		svc.Shutdown()

		logs, _, err := store.NewAuditStore(db).FindByAction(ctx, models.AuditActionSuperAdminGrant, 100, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 log, got %d", len(logs))
		}
		if logs[0].Action != models.AuditActionSuperAdminGrant {
			t.Errorf("Action = %v, want %v", logs[0].Action, models.AuditActionSuperAdminGrant)
		}
		if logs[0].ResourceType != "user" {
			t.Errorf("ResourceType = %v, want user", logs[0].ResourceType)
		}
		if logs[0].ResourceID == nil || *logs[0].ResourceID != 2 {
			t.Errorf("ResourceID = %v, want 2", logs[0].ResourceID)
		}
	})

	t.Run("revoke", func(t *testing.T) {
		db := setupTestDB(t)
		auditStore := store.NewAuditStore(db)
		svc := NewAuditService(auditStore)

		svc.LogSuperAdminChange(1, 3, "revoked@example.com", false, "127.0.0.1")
		svc.Shutdown()

		logs, _, err := store.NewAuditStore(db).FindByAction(ctx, models.AuditActionSuperAdminRevoke, 100, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(logs) != 1 {
			t.Fatalf("expected 1 log, got %d", len(logs))
		}
		if logs[0].Action != models.AuditActionSuperAdminRevoke {
			t.Errorf("Action = %v, want %v", logs[0].Action, models.AuditActionSuperAdminRevoke)
		}
	})
}

func TestAuditService_LogUserCreate(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogUserCreate(1, 10, "new@example.com", "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionUserCreate {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionUserCreate)
	}
	if log.ResourceType != "user" {
		t.Errorf("ResourceType = %v, want user", log.ResourceType)
	}
	if log.ResourceID == nil || *log.ResourceID != 10 {
		t.Errorf("ResourceID = %v, want 10", log.ResourceID)
	}
	if !log.Success {
		t.Error("expected Success = true")
	}
}

func TestAuditService_LogUserDelete(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogUserDelete(1, 20, "deleted@example.com", "10.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionUserDelete {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionUserDelete)
	}
	if log.ResourceType != "user" {
		t.Errorf("ResourceType = %v, want user", log.ResourceType)
	}
	if log.ResourceID == nil || *log.ResourceID != 20 {
		t.Errorf("ResourceID = %v, want 20", log.ResourceID)
	}

	var details map[string]interface{}
	if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
		t.Fatalf("failed to unmarshal details: %v", err)
	}
	if details["deleted_user_email"] != "deleted@example.com" {
		t.Errorf("details[deleted_user_email] = %v, want deleted@example.com", details["deleted_user_email"])
	}
}

func TestAuditService_LogUserAddToGroup(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogUserAddToGroup(1, 5, 10, "admin", "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionUserAddToGroup {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionUserAddToGroup)
	}
	if log.ResourceType != "user_group" {
		t.Errorf("ResourceType = %v, want user_group", log.ResourceType)
	}

	var details map[string]interface{}
	if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
		t.Fatalf("failed to unmarshal details: %v", err)
	}
	if details["group_id"].(float64) != 10 {
		t.Errorf("details[group_id] = %v, want 10", details["group_id"])
	}
	if details["role"] != "admin" {
		t.Errorf("details[role] = %v, want admin", details["role"])
	}
}

func TestAuditService_LogUserRemoveFromGroup(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogUserRemoveFromGroup(1, 5, 10, "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionUserRemoveFromGroup {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionUserRemoveFromGroup)
	}
	if log.ResourceType != "user_group" {
		t.Errorf("ResourceType = %v, want user_group", log.ResourceType)
	}

	var details map[string]interface{}
	if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
		t.Fatalf("failed to unmarshal details: %v", err)
	}
	if details["group_id"].(float64) != 10 {
		t.Errorf("details[group_id] = %v, want 10", details["group_id"])
	}
}

func TestAuditService_LogRoleChange(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogRoleChange(1, 5, 10, "manager", "admin", "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionRoleChange {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionRoleChange)
	}

	var details map[string]interface{}
	if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
		t.Fatalf("failed to unmarshal details: %v", err)
	}
	if details["old_role"] != "manager" {
		t.Errorf("details[old_role] = %v, want manager", details["old_role"])
	}
	if details["new_role"] != "admin" {
		t.Errorf("details[new_role] = %v, want admin", details["new_role"])
	}
}

func TestAuditService_LogResourceDelete(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		wantAction   models.AuditAction
	}{
		{"employee", "employee", models.AuditActionEmployeeDelete},
		{"child", "child", models.AuditActionChildDelete},
		{"organization", "organization", models.AuditActionOrgDelete},
		{"unknown type", "widget", "widget_delete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			auditStore := store.NewAuditStore(db)
			svc := NewAuditService(auditStore)
			ctx := context.Background()

			svc.LogResourceDelete(1, tt.resourceType, 42, "Test Resource", "127.0.0.1")
			svc.Shutdown()

			logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if total != 1 {
				t.Fatalf("expected total 1, got %d", total)
			}

			log := logs[0]
			if log.Action != tt.wantAction {
				t.Errorf("Action = %v, want %v", log.Action, tt.wantAction)
			}
			if log.ResourceType != tt.resourceType {
				t.Errorf("ResourceType = %v, want %v", log.ResourceType, tt.resourceType)
			}
			if log.ResourceID == nil || *log.ResourceID != 42 {
				t.Errorf("ResourceID = %v, want 42", log.ResourceID)
			}

			var details map[string]interface{}
			if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
				t.Fatalf("failed to unmarshal details: %v", err)
			}
			if details["resource_name"] != "Test Resource" {
				t.Errorf("details[resource_name] = %v, want Test Resource", details["resource_name"])
			}
		})
	}
}

func TestAuditService_LogResourceCreate(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogResourceCreate(1, "employee", 50, "John Doe", "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != "employee_create" {
		t.Errorf("Action = %v, want employee_create", log.Action)
	}
	if log.ResourceType != "employee" {
		t.Errorf("ResourceType = %v, want employee", log.ResourceType)
	}
	if log.ResourceID == nil || *log.ResourceID != 50 {
		t.Errorf("ResourceID = %v, want 50", log.ResourceID)
	}
}

func TestAuditService_LogResourceUpdate(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogResourceUpdate(1, "child", 30, "Jane Doe", "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != "child_update" {
		t.Errorf("Action = %v, want child_update", log.Action)
	}
	if log.ResourceType != "child" {
		t.Errorf("ResourceType = %v, want child", log.ResourceType)
	}
	if log.ResourceID == nil || *log.ResourceID != 30 {
		t.Errorf("ResourceID = %v, want 30", log.ResourceID)
	}
}

func TestAuditService_LogOrgCreate(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	svc.LogOrgCreate(1, 100, "Test Org", "127.0.0.1")
	svc.Shutdown()

	logs, total, err := store.NewAuditStore(db).FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}

	log := logs[0]
	if log.Action != models.AuditActionOrgCreate {
		t.Errorf("Action = %v, want %v", log.Action, models.AuditActionOrgCreate)
	}
	if log.ResourceType != "organization" {
		t.Errorf("ResourceType = %v, want organization", log.ResourceType)
	}
	if log.ResourceID == nil || *log.ResourceID != 100 {
		t.Errorf("ResourceID = %v, want 100", log.ResourceID)
	}

	var details map[string]interface{}
	if err := json.Unmarshal([]byte(log.Details), &details); err != nil {
		t.Fatalf("failed to unmarshal details: %v", err)
	}
	if details["org_name"] != "Test Org" {
		t.Errorf("details[org_name] = %v, want Test Org", details["org_name"])
	}
}

func TestAuditService_GetLogs(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	// Add multiple logs
	for i := 0; i < 5; i++ {
		svc.LogLogin(uint(i+1), "user@example.com", "127.0.0.1", "Agent")
	}
	svc.Shutdown()

	// Use a read-only service (no channel needed)
	readSvc := &AuditService{store: store.NewAuditStore(db)}

	// Verify total
	logs, total, err := readSvc.GetLogs(ctx, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(logs) != 5 {
		t.Errorf("expected 5 logs, got %d", len(logs))
	}

	// Test pagination - limit 2
	logs, total, err = readSvc.GetLogs(ctx, 2, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5 with limit, got %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs with limit, got %d", len(logs))
	}

	// Test pagination - offset 3
	logs, total, err = readSvc.GetLogs(ctx, 100, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5 with offset, got %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs with offset 3, got %d", len(logs))
	}
}

func TestAuditService_GetLogsByUser(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	// Log for user 1
	svc.LogLogin(1, "user1@example.com", "127.0.0.1", "Agent")
	svc.LogLogin(1, "user1@example.com", "127.0.0.1", "Agent")

	// Log for user 2
	svc.LogLogin(2, "user2@example.com", "127.0.0.1", "Agent")

	svc.Shutdown()

	readSvc := &AuditService{store: store.NewAuditStore(db)}

	// Filter by user 1
	logs, total, err := readSvc.GetLogsByUser(ctx, 1, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2 for user 1, got %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs for user 1, got %d", len(logs))
	}

	// Filter by user 2
	logs, total, err = readSvc.GetLogsByUser(ctx, 2, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1 for user 2, got %d", total)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log for user 2, got %d", len(logs))
	}

	// Non-existent user
	logs, total, err = readSvc.GetLogsByUser(ctx, 999, 100, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0 for non-existent user, got %d", total)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 logs for non-existent user, got %d", len(logs))
	}
}

func TestAuditService_CountRecentFailedLogins(t *testing.T) {
	db := setupTestDB(t)
	auditStore := store.NewAuditStore(db)
	svc := NewAuditService(auditStore)
	ctx := context.Background()

	// Add failed login attempts
	svc.LogLoginFailed("fail@example.com", "127.0.0.1", "Agent", "bad password")
	svc.LogLoginFailed("fail@example.com", "127.0.0.1", "Agent", "bad password")
	svc.LogLoginFailed("fail@example.com", "127.0.0.1", "Agent", "bad password")
	// Different email
	svc.LogLoginFailed("other@example.com", "127.0.0.1", "Agent", "bad password")

	svc.Shutdown()

	readSvc := &AuditService{store: store.NewAuditStore(db)}

	count, err := readSvc.CountRecentFailedLogins(ctx, "fail@example.com", 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}

	// Different email should have 1
	count, err = readSvc.CountRecentFailedLogins(ctx, "other@example.com", 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Non-existent email should be 0
	count, err = readSvc.CountRecentFailedLogins(ctx, "none@example.com", 1*time.Hour)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

func TestAuditService_NilSafety(t *testing.T) {
	ctx := context.Background()

	t.Run("nil service", func(t *testing.T) {
		var svc *AuditService

		// GetLogs returns empty
		logs, total, err := svc.GetLogs(ctx, 100, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logs != nil {
			t.Errorf("expected nil logs, got %v", logs)
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}

		// GetLogsByUser returns empty
		logs, total, err = svc.GetLogsByUser(ctx, 1, 100, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if logs != nil {
			t.Errorf("expected nil logs, got %v", logs)
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}

		// CountRecentFailedLogins returns 0
		count, err := svc.CountRecentFailedLogins(ctx, "test@example.com", 1*time.Hour)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}

		// Shutdown doesn't panic
		svc.Shutdown()
	})

	t.Run("nil channel", func(t *testing.T) {
		svc := &AuditService{}

		// Log methods should not panic with nil channel
		svc.LogLogin(1, "test@example.com", "127.0.0.1", "Agent")
		svc.LogLoginFailed("test@example.com", "127.0.0.1", "Agent", "reason")
		svc.LogSuperAdminChange(1, 2, "test@example.com", true, "127.0.0.1")
		svc.LogUserCreate(1, 2, "test@example.com", "127.0.0.1")
		svc.LogUserDelete(1, 2, "test@example.com", "127.0.0.1")
		svc.LogUserAddToGroup(1, 2, 3, "admin", "127.0.0.1")
		svc.LogUserRemoveFromGroup(1, 2, 3, "127.0.0.1")
		svc.LogRoleChange(1, 2, 3, "old", "new", "127.0.0.1")
		svc.LogResourceDelete(1, "employee", 2, "name", "127.0.0.1")
		svc.LogResourceCreate(1, "employee", 2, "name", "127.0.0.1")
		svc.LogResourceUpdate(1, "employee", 2, "name", "127.0.0.1")
		svc.LogOrgCreate(1, 2, "org", "127.0.0.1")

		// Shutdown doesn't panic
		svc.Shutdown()
	})
}
