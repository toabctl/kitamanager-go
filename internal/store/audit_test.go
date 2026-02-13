package store

import (
	"context"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func uintPtr(v uint) *uint { return &v }

func createTestAuditLog(t *testing.T, s *AuditStore, userID uint, action models.AuditAction) *models.AuditLog {
	t.Helper()
	log := &models.AuditLog{
		UserID:       uintPtr(userID),
		UserEmail:    "test@example.com",
		Action:       action,
		ResourceType: "user",
		Details:      "test details",
		IPAddress:    "127.0.0.1",
		Timestamp:    time.Now(),
	}
	if err := s.Create(context.Background(), log); err != nil {
		t.Fatalf("failed to create audit log: %v", err)
	}
	return log
}

func TestAuditStore_Create(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	log := &models.AuditLog{
		UserID:       uintPtr(1),
		UserEmail:    "admin@example.com",
		Action:       models.AuditActionLogin,
		ResourceType: "auth",
		Details:      "successful login",
		IPAddress:    "192.168.1.1",
	}

	err := store.Create(context.Background(), log)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if log.ID == 0 {
		t.Error("expected log ID to be set")
	}
	if log.Timestamp.IsZero() {
		t.Error("expected timestamp to be auto-set")
	}
}

func TestAuditStore_FindByUser(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	createTestAuditLog(t, store, 1, models.AuditActionLogin)
	createTestAuditLog(t, store, 1, models.AuditActionUserCreate)
	createTestAuditLog(t, store, 2, models.AuditActionLogin)

	logs, total, err := store.FindByUser(context.Background(), 1, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(logs))
	}
}

func TestAuditStore_FindByUser_Pagination(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	for i := 0; i < 5; i++ {
		createTestAuditLog(t, store, 1, models.AuditActionLogin)
	}

	logs, total, err := store.FindByUser(context.Background(), 1, 2, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs (limit), got %d", len(logs))
	}

	logs2, _, err := store.FindByUser(context.Background(), 1, 2, 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(logs2) != 2 {
		t.Errorf("expected 2 logs (offset), got %d", len(logs2))
	}
}

func TestAuditStore_FindByAction(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	createTestAuditLog(t, store, 1, models.AuditActionLogin)
	createTestAuditLog(t, store, 1, models.AuditActionUserCreate)
	createTestAuditLog(t, store, 2, models.AuditActionUserCreate)

	logs, total, err := store.FindByAction(context.Background(), models.AuditActionUserCreate, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(logs))
	}
}

func TestAuditStore_FindByDateRange(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	now := time.Now()

	old := &models.AuditLog{
		UserID: uintPtr(1), UserEmail: "test@example.com", Action: models.AuditActionLogin,
		ResourceType: "auth", Timestamp: now.Add(-48 * time.Hour),
	}
	_ = store.Create(context.Background(), old)

	recent := &models.AuditLog{
		UserID: uintPtr(1), UserEmail: "test@example.com", Action: models.AuditActionLogin,
		ResourceType: "auth", Timestamp: now.Add(-1 * time.Hour),
	}
	_ = store.Create(context.Background(), recent)

	from := now.Add(-24 * time.Hour)
	to := now

	logs, total, err := store.FindByDateRange(context.Background(), from, to, 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
}

func TestAuditStore_FindAll(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	createTestAuditLog(t, store, 1, models.AuditActionLogin)
	createTestAuditLog(t, store, 2, models.AuditActionUserCreate)
	createTestAuditLog(t, store, 3, models.AuditActionUserDelete)

	logs, total, err := store.FindAll(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(logs) != 3 {
		t.Errorf("expected 3 logs, got %d", len(logs))
	}
}

func TestAuditStore_FindFailedLogins(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	now := time.Now()

	failed := &models.AuditLog{
		UserEmail: "hacker@example.com", Action: models.AuditActionLoginFailed,
		ResourceType: "auth", Timestamp: now.Add(-5 * time.Minute),
	}
	_ = store.Create(context.Background(), failed)

	success := &models.AuditLog{
		UserID: uintPtr(1), UserEmail: "user@example.com", Action: models.AuditActionLogin,
		ResourceType: "auth", Timestamp: now.Add(-3 * time.Minute),
	}
	_ = store.Create(context.Background(), success)

	// Find all failed logins
	logs, err := store.FindFailedLogins(context.Background(), "", now.Add(-10*time.Minute), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 failed login, got %d", len(logs))
	}

	// Filter by email
	logs, err = store.FindFailedLogins(context.Background(), "hacker@example.com", now.Add(-10*time.Minute), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 failed login for email, got %d", len(logs))
	}

	// Filter by different email
	logs, err = store.FindFailedLogins(context.Background(), "other@example.com", now.Add(-10*time.Minute), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 failed logins, got %d", len(logs))
	}
}

func TestAuditStore_CountFailedLoginsSince(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	now := time.Now()
	email := "attacker@example.com"

	for i := 0; i < 3; i++ {
		log := &models.AuditLog{
			UserEmail: email, Action: models.AuditActionLoginFailed,
			ResourceType: "auth", Timestamp: now.Add(-time.Duration(i) * time.Minute),
		}
		_ = store.Create(context.Background(), log)
	}

	count, err := store.CountFailedLoginsSince(context.Background(), email, now.Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3, got %d", count)
	}

	count, err = store.CountFailedLoginsSince(context.Background(), "other@example.com", now.Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestAuditStore_Cleanup(t *testing.T) {
	db := setupTestDB(t)
	store := NewAuditStore(db)

	now := time.Now()

	old := &models.AuditLog{
		UserID: uintPtr(1), UserEmail: "test@example.com", Action: models.AuditActionLogin,
		ResourceType: "auth", Timestamp: now.Add(-48 * time.Hour),
	}
	_ = store.Create(context.Background(), old)

	recent := &models.AuditLog{
		UserID: uintPtr(1), UserEmail: "test@example.com", Action: models.AuditActionLogin,
		ResourceType: "auth", Timestamp: now.Add(-1 * time.Hour),
	}
	_ = store.Create(context.Background(), recent)

	deleted, err := store.Cleanup(context.Background(), now.Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if deleted != 1 {
		t.Errorf("expected 1 deleted, got %d", deleted)
	}

	_, total, _ := store.FindAll(context.Background(), 10, 0)
	if total != 1 {
		t.Errorf("expected 1 remaining, got %d", total)
	}
}
