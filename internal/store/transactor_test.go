package store

import (
	"context"
	"errors"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestTransactor_InTransaction_Commit(t *testing.T) {
	db := setupTestDB(t)
	transactor := NewTransactor(db)

	org := createTestOrganization(t, db, "Test Org")

	// Create a section inside a transaction using DBFromContext
	err := transactor.InTransaction(context.Background(), func(ctx context.Context) error {
		section := &models.Section{
			Name:           "Tx Section",
			OrganizationID: org.ID,
		}
		return DBFromContext(ctx, db).Create(section).Error
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify section was committed
	var count int64
	db.Model(&models.Section{}).Where("name = ?", "Tx Section").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 section, got %d", count)
	}
}

func TestTransactor_InTransaction_Rollback(t *testing.T) {
	db := setupTestDB(t)
	transactor := NewTransactor(db)

	org := createTestOrganization(t, db, "Test Org")

	// Create a section then return an error to trigger rollback
	err := transactor.InTransaction(context.Background(), func(ctx context.Context) error {
		section := &models.Section{
			Name:           "Rolled Back Section",
			OrganizationID: org.ID,
		}
		if err := DBFromContext(ctx, db).Create(section).Error; err != nil {
			return err
		}
		return errors.New("simulated failure")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify section was rolled back
	var count int64
	db.Model(&models.Section{}).Where("name = ?", "Rolled Back Section").Count(&count)
	if count != 0 {
		t.Errorf("expected 0 sections after rollback, got %d", count)
	}
}

func TestTransactor_InTransaction_Nested(t *testing.T) {
	db := setupTestDB(t)
	transactor := NewTransactor(db)

	org := createTestOrganization(t, db, "Test Org")

	// Nested calls should reuse the outer transaction
	err := transactor.InTransaction(context.Background(), func(ctx context.Context) error {
		section1 := &models.Section{
			Name:           "Nested Section 1",
			OrganizationID: org.ID,
		}
		if err := DBFromContext(ctx, db).Create(section1).Error; err != nil {
			return err
		}

		// Nested transaction
		return transactor.InTransaction(ctx, func(innerCtx context.Context) error {
			section2 := &models.Section{
				Name:           "Nested Section 2",
				OrganizationID: org.ID,
			}
			return DBFromContext(innerCtx, db).Create(section2).Error
		})
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Both sections should be committed
	var count int64
	db.Model(&models.Section{}).Where("name LIKE ?", "Nested Section%").Count(&count)
	if count != 2 {
		t.Errorf("expected 2 sections, got %d", count)
	}
}

func TestTransactor_DBFromContext_WithoutTransaction(t *testing.T) {
	db := setupTestDB(t)

	// Without a transaction in context, DBFromContext should return the default DB
	ctx := context.Background()
	result := DBFromContext(ctx, db)
	if result == nil {
		t.Fatal("expected non-nil DB")
	}
}

func TestTransactor_DBFromContext_WithTransaction(t *testing.T) {
	db := setupTestDB(t)
	transactor := NewTransactor(db)

	org := createTestOrganization(t, db, "Test Org")

	// Inside a transaction, DBFromContext should return the tx DB
	err := transactor.InTransaction(context.Background(), func(ctx context.Context) error {
		txDB := DBFromContext(ctx, db)

		// Create using the tx-aware DB
		section := &models.Section{
			Name:           "TxDB Section",
			OrganizationID: org.ID,
		}
		return txDB.Create(section).Error
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var count int64
	db.Model(&models.Section{}).Where("name = ?", "TxDB Section").Count(&count)
	if count != 1 {
		t.Errorf("expected 1 section, got %d", count)
	}
}
