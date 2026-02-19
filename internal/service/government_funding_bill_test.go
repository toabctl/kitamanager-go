package service

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestGovernmentFundingBillService_ListEmpty(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	ctx := context.Background()

	items, total, err := svc.List(ctx, org.ID, 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
}

func TestGovernmentFundingBillService_ListWithPeriods(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc1@example.com", "password")
	ctx := context.Background()

	// Create test periods directly in DB
	for i := 0; i < 3; i++ {
		month := time.Month(i + 1)
		to := time.Date(2025, month+1, 0, 0, 0, 0, 0, time.UTC)
		p := &models.GovernmentFundingBillPeriod{
			OrganizationID:    org.ID,
			Period:            models.Period{From: time.Date(2025, month, 1, 0, 0, 0, 0, time.UTC), To: &to},
			FileName:          "file.xlsx",
			FileSha256:        "hash",
			FacilityName:      "Kita",
			FacilityTotal:     100000 * (i + 1),
			ContractBooking:   90000 * (i + 1),
			CorrectionBooking: 10000 * (i + 1),
			CreatedBy:         user.ID,
		}
		if err := db.Create(p).Error; err != nil {
			t.Fatalf("setup: Create() error = %v", err)
		}
	}

	items, total, err := svc.List(ctx, org.ID, 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}

	// Verify response fields are populated
	first := items[0]
	if first.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if first.FacilityName != "Kita" {
		t.Errorf("expected facility name 'Kita', got %q", first.FacilityName)
	}
	if first.From == "" {
		t.Error("expected non-empty From")
	}
}

func TestGovernmentFundingBillService_ListOrganizationIsolation(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	user := createTestUser(t, db, "User", "billsvc2@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	if err := db.Create(&models.GovernmentFundingBillPeriod{
		OrganizationID: org1.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "org1.xlsx",
		FileSha256:     "hash1",
		FacilityName:   "Org1 Kita",
		CreatedBy:      user.ID,
	}).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	// Org2 should see nothing
	items, total, err := svc.List(ctx, org2.ID, 10, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 0 {
		t.Errorf("org2 should see 0 bills, got %d", total)
	}
	if len(items) != 0 {
		t.Errorf("org2 should see 0 items, got %d", len(items))
	}
}

func TestGovernmentFundingBillService_GetByID(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc3@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID:    org.ID,
		Period:            models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:          "detail.xlsx",
		FileSha256:        "detailhash",
		FacilityName:      "Kita Detail",
		FacilityTotal:     300000,
		ContractBooking:   280000,
		CorrectionBooking: 20000,
		CreatedBy:         user.ID,
		Children: []models.GovernmentFundingBillChild{
			{
				VoucherNumber: "GB-11111111111-01",
				ChildName:     "Kind, Eins",
				BirthDate:     "05.19",
				District:      2,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "ganztag", Amount: 120000},
					{Key: "ndh", Value: "ndh", Amount: 8000},
					{Key: "qm/mss", Value: "qm/mss", Amount: 5000},
				},
			},
			{
				VoucherNumber: "GB-22222222222-01",
				ChildName:     "Kind, Zwei",
				BirthDate:     "08.20",
				District:      5,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "halbtag", Amount: 90000},
					{Key: "sph", Value: "sph", Amount: 12000},
				},
			},
		},
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	result, err := svc.GetByID(ctx, period.ID, org.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	// Header fields
	if result.ID != period.ID {
		t.Errorf("expected ID %d, got %d", period.ID, result.ID)
	}
	if result.FacilityName != "Kita Detail" {
		t.Errorf("expected facility name 'Kita Detail', got %q", result.FacilityName)
	}
	if result.FacilityTotal != 300000 {
		t.Errorf("expected facility total 300000, got %d", result.FacilityTotal)
	}
	if result.ChildrenCount != 2 {
		t.Errorf("expected children count 2, got %d", result.ChildrenCount)
	}

	// All unmatched (no children in system)
	if result.MatchedCount != 0 {
		t.Errorf("expected matched count 0, got %d", result.MatchedCount)
	}
	if result.UnmatchedCount != 2 {
		t.Errorf("expected unmatched count 2, got %d", result.UnmatchedCount)
	}

	// Surcharges aggregated
	if len(result.Surcharges) != 3 {
		t.Fatalf("expected 3 surcharges, got %d", len(result.Surcharges))
	}
	surchargeMap := map[string]int{}
	for _, s := range result.Surcharges {
		surchargeMap[s.Key] = s.Amount
	}
	if surchargeMap["ndh"] != 8000 {
		t.Errorf("expected ndh surcharge 8000, got %d", surchargeMap["ndh"])
	}
	if surchargeMap["qm/mss"] != 5000 {
		t.Errorf("expected qm/mss surcharge 5000, got %d", surchargeMap["qm/mss"])
	}
	if surchargeMap["sph"] != 12000 {
		t.Errorf("expected sph surcharge 12000, got %d", surchargeMap["sph"])
	}

	// Children detail
	if len(result.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(result.Children))
	}
	child1 := result.Children[0]
	if child1.TotalAmount != 133000 { // 120000 + 8000 + 5000
		t.Errorf("expected child1 total 133000, got %d", child1.TotalAmount)
	}
	if child1.Matched {
		t.Error("expected child1 not matched")
	}
}

func TestGovernmentFundingBillService_GetByIDWithMatching(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc4@example.com", "password")
	section := getDefaultSection(t, db, org.ID)
	ctx := context.Background()

	// Create children with voucher numbers in the system
	child := createTestChild(t, db, "Max", "Mustermann", org.ID)
	voucher := "GB-99999999999-01"
	if err := db.Create(&models.ChildContract{
		ChildID: child.ID,
		BaseContract: models.BaseContract{
			Period:    models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			SectionID: section.ID,
		},
		VoucherNumber: &voucher,
	}).Error; err != nil {
		t.Fatalf("setup: create contract error = %v", err)
	}

	// Create bill period with one matching and one non-matching child
	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "match-test.xlsx",
		FileSha256:     "matchhash",
		FacilityName:   "Kita Match",
		FacilityTotal:  200000,
		CreatedBy:      user.ID,
		Children: []models.GovernmentFundingBillChild{
			{
				VoucherNumber: voucher, // matches
				ChildName:     "Mustermann, Max",
				BirthDate:     "01.20",
				District:      1,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "ganztag", Amount: 100000},
				},
			},
			{
				VoucherNumber: "GB-00000000000-01", // no match
				ChildName:     "Unknown, Child",
				BirthDate:     "06.21",
				District:      3,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "halbtag", Amount: 80000},
				},
			},
		},
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	result, err := svc.GetByID(ctx, period.ID, org.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if result.MatchedCount != 1 {
		t.Errorf("expected matched count 1, got %d", result.MatchedCount)
	}
	if result.UnmatchedCount != 1 {
		t.Errorf("expected unmatched count 1, got %d", result.UnmatchedCount)
	}

	// Find matched child in response
	var matchedChild *models.GovernmentFundingBillChildResponse
	var unmatchedChild *models.GovernmentFundingBillChildResponse
	for i := range result.Children {
		if result.Children[i].VoucherNumber == voucher {
			matchedChild = &result.Children[i]
		} else {
			unmatchedChild = &result.Children[i]
		}
	}

	if matchedChild == nil {
		t.Fatal("expected to find matched child in response")
	}
	if !matchedChild.Matched {
		t.Error("expected matched child to have Matched=true")
	}
	if matchedChild.ChildID == nil || *matchedChild.ChildID != child.ID {
		t.Errorf("expected child_id %d, got %v", child.ID, matchedChild.ChildID)
	}
	if matchedChild.ContractID == nil {
		t.Error("expected contract_id to be set for matched child")
	}

	if unmatchedChild == nil {
		t.Fatal("expected to find unmatched child in response")
	}
	if unmatchedChild.Matched {
		t.Error("expected unmatched child to have Matched=false")
	}
	if unmatchedChild.ChildID != nil {
		t.Error("expected unmatched child_id to be nil")
	}
}

func TestGovernmentFundingBillService_GetByIDWrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	user := createTestUser(t, db, "User", "billsvc5@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org1.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "wrong-org.xlsx",
		FileSha256:     "wrongorghash",
		FacilityName:   "Kita",
		CreatedBy:      user.ID,
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	// Try to access from org2 — should fail
	_, err := svc.GetByID(ctx, period.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error for wrong org, got nil")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestGovernmentFundingBillService_GetByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, 99999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent ID, got nil")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestGovernmentFundingBillService_Delete(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc6@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "delete.xlsx",
		FileSha256:     "deletehash",
		FacilityName:   "Kita Delete",
		CreatedBy:      user.ID,
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	deleted, err := svc.Delete(ctx, period.ID, org.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if deleted.ID != period.ID {
		t.Errorf("expected deleted period ID %d, got %d", period.ID, deleted.ID)
	}

	// Verify it's gone
	_, err = svc.GetByID(ctx, period.ID, org.ID)
	if err == nil {
		t.Error("expected error after delete, got nil")
	}
}

func TestGovernmentFundingBillService_DeleteWrongOrg(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	user := createTestUser(t, db, "User", "billsvc7@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org1.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "protected.xlsx",
		FileSha256:     "protectedhash",
		FacilityName:   "Kita Protected",
		CreatedBy:      user.ID,
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	// Try to delete from org2 — should fail
	_, err := svc.Delete(ctx, period.ID, org2.ID)
	if err == nil {
		t.Fatal("expected error for wrong org delete, got nil")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}

	// Verify it still exists
	_, err = svc.GetByID(ctx, period.ID, org1.ID)
	if err != nil {
		t.Errorf("period should still exist after failed delete: %v", err)
	}
}

func TestGovernmentFundingBillService_DeleteNotFound(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	ctx := context.Background()

	_, err := svc.Delete(ctx, 99999, 1)
	if err == nil {
		t.Fatal("expected error for non-existent ID, got nil")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestGovernmentFundingBillService_GetByIDSurchargesAggregation(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc8@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "surcharge-agg.xlsx",
		FileSha256:     "surchargehash",
		FacilityName:   "Kita Surcharge",
		FacilityTotal:  500000,
		CreatedBy:      user.ID,
		Children: []models.GovernmentFundingBillChild{
			{
				VoucherNumber: "GB-11111111111-01",
				ChildName:     "Kind A",
				BirthDate:     "01.20",
				District:      1,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "ganztag", Amount: 100000},
					{Key: "ndh", Value: "ndh", Amount: 5000},
					{Key: "qm/mss", Value: "qm/mss", Amount: 3000},
				},
			},
			{
				VoucherNumber: "GB-22222222222-01",
				ChildName:     "Kind B",
				BirthDate:     "02.20",
				District:      2,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "ganztag", Amount: 110000},
					{Key: "ndh", Value: "ndh", Amount: 7000},
					{Key: "sph", Value: "sph", Amount: 15000},
				},
			},
			{
				VoucherNumber: "GB-33333333333-01",
				ChildName:     "Kind C",
				BirthDate:     "03.20",
				District:      3,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "halbtag", Amount: 80000},
					// No surcharges for this child
				},
			},
		},
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	result, err := svc.GetByID(ctx, period.ID, org.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	// Verify surcharges are summed across children
	surchargeMap := map[string]int{}
	for _, s := range result.Surcharges {
		surchargeMap[s.Key] = s.Amount
	}

	// ndh: 5000 + 7000 = 12000
	if surchargeMap["ndh"] != 12000 {
		t.Errorf("expected ndh=12000, got %d", surchargeMap["ndh"])
	}
	// qm/mss: 3000 (only from child A)
	if surchargeMap["qm/mss"] != 3000 {
		t.Errorf("expected qm/mss=3000, got %d", surchargeMap["qm/mss"])
	}
	// sph: 15000 (only from child B)
	if surchargeMap["sph"] != 15000 {
		t.Errorf("expected sph=15000, got %d", surchargeMap["sph"])
	}
}

func TestGovernmentFundingBillService_GetByIDNoSurcharges(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc9@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "no-surcharge.xlsx",
		FileSha256:     "nosurchargehash",
		FacilityName:   "Kita No Surcharge",
		FacilityTotal:  100000,
		CreatedBy:      user.ID,
		Children: []models.GovernmentFundingBillChild{
			{
				VoucherNumber: "GB-44444444444-01",
				ChildName:     "Kind, Plain",
				BirthDate:     "04.20",
				District:      1,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "ganztag", Amount: 100000},
				},
			},
		},
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	result, err := svc.GetByID(ctx, period.ID, org.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	// Surcharges should exist but all be 0
	if len(result.Surcharges) != 3 {
		t.Fatalf("expected 3 surcharge entries, got %d", len(result.Surcharges))
	}
	for _, s := range result.Surcharges {
		if s.Amount != 0 {
			t.Errorf("expected surcharge %q to be 0, got %d", s.Key, s.Amount)
		}
	}
}

func TestGovernmentFundingBillService_GetByIDNoChildren(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org := createTestOrganization(t, db, "Test Org")
	user := createTestUser(t, db, "User", "billsvc10@example.com", "password")
	ctx := context.Background()

	to := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org.ID,
		Period:         models.Period{From: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "empty-children.xlsx",
		FileSha256:     "emptychildrenhash",
		FacilityName:   "Kita Empty",
		FacilityTotal:  0,
		CreatedBy:      user.ID,
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	result, err := svc.GetByID(ctx, period.ID, org.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if result.ChildrenCount != 0 {
		t.Errorf("expected children count 0, got %d", result.ChildrenCount)
	}
	if result.MatchedCount != 0 {
		t.Errorf("expected matched count 0, got %d", result.MatchedCount)
	}
	if result.UnmatchedCount != 0 {
		t.Errorf("expected unmatched count 0, got %d", result.UnmatchedCount)
	}
}

func TestComputeFileHash(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		content := []byte("hello world")
		hash1, err := ComputeFileHash(bytes.NewReader(content))
		if err != nil {
			t.Fatalf("ComputeFileHash() error = %v", err)
		}
		hash2, err := ComputeFileHash(bytes.NewReader(content))
		if err != nil {
			t.Fatalf("ComputeFileHash() error = %v", err)
		}
		if hash1 != hash2 {
			t.Error("expected same hash for same content")
		}
	})

	t.Run("different content produces different hash", func(t *testing.T) {
		hash1, _ := ComputeFileHash(bytes.NewReader([]byte("hello")))
		hash2, _ := ComputeFileHash(bytes.NewReader([]byte("world")))
		if hash1 == hash2 {
			t.Error("expected different hashes for different content")
		}
	})

	t.Run("empty content", func(t *testing.T) {
		hash, err := ComputeFileHash(bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatalf("ComputeFileHash() error = %v", err)
		}
		if hash == "" {
			t.Error("expected non-empty hash for empty content")
		}
		// SHA-256 of empty is well-known
		if hash != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
			t.Errorf("unexpected empty hash: %s", hash)
		}
	})

	t.Run("correct sha256 length", func(t *testing.T) {
		hash, _ := ComputeFileHash(bytes.NewReader([]byte("test")))
		if len(hash) != 64 {
			t.Errorf("expected 64 char hex string, got %d chars", len(hash))
		}
	})
}

func TestLastDayOfMonth(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "November 2025 (30 days)",
			input:    time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "February 2024 (leap year)",
			input:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "February 2025 (non-leap year)",
			input:    time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "December 2025 (31 days)",
			input:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "January 2026 (31 days)",
			input:    time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lastDayOfMonth(tt.input)
			if !result.Equal(tt.expected) {
				t.Errorf("lastDayOfMonth(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatToDate(t *testing.T) {
	t.Run("nil returns empty", func(t *testing.T) {
		result := formatToDate(nil)
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("non-nil returns formatted", func(t *testing.T) {
		d := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
		result := formatToDate(&d)
		if result != "2025-11-30" {
			t.Errorf("expected '2025-11-30', got %q", result)
		}
	})
}

func TestGovernmentFundingBillService_GetByIDMatchingCrossTenancy(t *testing.T) {
	db := setupTestDB(t)
	childStore := store.NewChildStore(db)
	billPeriodStore := store.NewGovernmentFundingBillPeriodStore(db)
	svc := NewGovernmentFundingBillService(childStore, billPeriodStore)
	org1 := createTestOrganization(t, db, "Org 1")
	org2 := createTestOrganization(t, db, "Org 2")
	user := createTestUser(t, db, "User", "billsvc11@example.com", "password")
	section2 := getDefaultSection(t, db, org2.ID)
	ctx := context.Background()

	// Create child with voucher in org2
	child := createTestChild(t, db, "Other", "OrgChild", org2.ID)
	voucher := "GB-55555555555-01"
	if err := db.Create(&models.ChildContract{
		ChildID: child.ID,
		BaseContract: models.BaseContract{
			Period:    models.Period{From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			SectionID: section2.ID,
		},
		VoucherNumber: &voucher,
	}).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	// Create bill period in org1 with same voucher
	to := time.Date(2025, 11, 30, 0, 0, 0, 0, time.UTC)
	period := &models.GovernmentFundingBillPeriod{
		OrganizationID: org1.ID,
		Period:         models.Period{From: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), To: &to},
		FileName:       "cross-tenant.xlsx",
		FileSha256:     "crosshash",
		FacilityName:   "Org1 Kita",
		CreatedBy:      user.ID,
		Children: []models.GovernmentFundingBillChild{
			{
				VoucherNumber: voucher,
				ChildName:     "Cross, Tenant",
				BirthDate:     "01.20",
				District:      1,
				Payments: []models.GovernmentFundingBillPayment{
					{Key: "care_type", Value: "ganztag", Amount: 100000},
				},
			},
		},
	}
	if err := db.Create(period).Error; err != nil {
		t.Fatalf("setup error: %v", err)
	}

	result, err := svc.GetByID(ctx, period.ID, org1.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	// Child in org2 should NOT match for org1's bill
	if result.MatchedCount != 0 {
		t.Errorf("expected 0 matches (cross-org should not match), got %d", result.MatchedCount)
	}
}
