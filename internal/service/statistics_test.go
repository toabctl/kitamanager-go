package service

import (
	"context"
	"math"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func createStatisticsService(db *gorm.DB) *StatisticsService {
	childStore := store.NewChildStore(db)
	employeeStore := store.NewEmployeeStore(db)
	orgStore := store.NewOrganizationStore(db)
	fundingStore := store.NewGovernmentFundingStore(db)
	payPlanStore := store.NewPayPlanStore(db)
	costStore := store.NewCostStore(db)
	return NewStatisticsService(childStore, employeeStore, orgStore, fundingStore, payPlanStore, costStore)
}

func createTestEmployeeContractWithCategory(t *testing.T, db *gorm.DB, employeeID uint, payplanID uint, from time.Time, to *time.Time, weeklyHours float64, staffCategory string, sectionID uint) *models.EmployeeContract {
	t.Helper()
	contract := &models.EmployeeContract{
		EmployeeID: employeeID,
		BaseContract: models.BaseContract{
			Period:    models.Period{From: from, To: to},
			SectionID: sectionID,
		},
		StaffCategory: staffCategory,
		WeeklyHours:   weeklyHours,
		PayPlanID:     payplanID,
	}
	if err := db.Create(contract).Error; err != nil {
		t.Fatalf("failed to create test employee contract: %v", err)
	}
	return contract
}

func createTestChildContract(t *testing.T, db *gorm.DB, childID uint, from time.Time, to *time.Time, sectionID uint, properties models.ContractProperties) *models.ChildContract {
	t.Helper()
	contract := &models.ChildContract{
		ChildID: childID,
		BaseContract: models.BaseContract{
			Period:     models.Period{From: from, To: to},
			SectionID:  sectionID,
			Properties: properties,
		},
	}
	if err := db.Create(contract).Error; err != nil {
		t.Fatalf("failed to create test child contract: %v", err)
	}
	return contract
}

func createTestFundingPropertyWithRequirement(t *testing.T, db *gorm.DB, periodID uint, key, value string, requirement float64, minAge, maxAge int) *models.GovernmentFundingProperty {
	t.Helper()
	var minAgePtr, maxAgePtr *int
	if minAge >= 0 {
		minAgePtr = &minAge
	}
	if maxAge >= 0 {
		maxAgePtr = &maxAge
	}
	prop := &models.GovernmentFundingProperty{
		PeriodID:    periodID,
		Key:         key,
		Value:       value,
		Payment:     10000,
		Requirement: requirement,
		MinAge:      minAgePtr,
		MaxAge:      maxAgePtr,
	}
	if err := db.Create(prop).Error; err != nil {
		t.Fatalf("failed to create test funding property: %v", err)
	}
	return prop
}

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestStatisticsService_GetStaffingHours_Basic(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	// Create org with state "berlin"
	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	// Create government funding with period covering 2024
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)

	// Funding property: care_type=ganztag, requirement=0.25, ages 0-6
	createTestFundingPropertyWithRequirement(t, db, period.ID, "care_type", "ganztag", 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// 2 children with contracts from 2024-01-01, ongoing
	child1 := createTestChild(t, db, "Child", "One", org.ID)
	child2 := createTestChild(t, db, "Child", "Two", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}
	createTestChildContract(t, db, child1.ID, contractFrom, nil, section.ID, props)
	createTestChildContract(t, db, child2.ID, contractFrom, nil, section.ID, props)

	// 2 employees with qualified contracts, 30 hours each
	payplan := createTestPayPlan(t, db, "TV-L", org.ID)
	emp1 := createTestEmployee(t, db, "Emp", "One", org.ID)
	emp2 := createTestEmployee(t, db, "Emp", "Two", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp1.ID, payplan.ID, contractFrom, nil, 30.0, "qualified", section.ID)
	createTestEmployeeContractWithCategory(t, db, emp2.ID, payplan.ID, contractFrom, nil, 30.0, "qualified", section.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 6 data points: Jan, Feb, Mar, Apr, May, Jun
	if len(result.DataPoints) != 6 {
		t.Fatalf("expected 6 data points, got %d", len(result.DataPoints))
	}

	for i, dp := range result.DataPoints {
		// RequiredHours = 2 children * 0.25 requirement * 39.0 full-time hours = 19.5
		if !almostEqual(dp.RequiredHours, 19.5, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want 19.5", i, dp.RequiredHours)
		}
		// AvailableHours = 2 employees * 30.0 hours = 60.0
		if !almostEqual(dp.AvailableHours, 60.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 60.0", i, dp.AvailableHours)
		}
		if dp.ChildCount != 2 {
			t.Errorf("data point %d: ChildCount = %d, want 2", i, dp.ChildCount)
		}
		if dp.StaffCount != 2 {
			t.Errorf("data point %d: StaffCount = %d, want 2", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetStaffingHours_NoChildren(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyWithRequirement(t, db, period.ID, "care_type", "ganztag", 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Only employees, no children
	payplan := createTestPayPlan(t, db, "TV-L", org.ID)
	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 30.0, "qualified", section.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, dp := range result.DataPoints {
		if !almostEqual(dp.RequiredHours, 0.0, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want 0.0", i, dp.RequiredHours)
		}
		if !almostEqual(dp.AvailableHours, 30.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 30.0", i, dp.AvailableHours)
		}
		if dp.ChildCount != 0 {
			t.Errorf("data point %d: ChildCount = %d, want 0", i, dp.ChildCount)
		}
		if dp.StaffCount != 1 {
			t.Errorf("data point %d: StaffCount = %d, want 1", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetStaffingHours_NoEmployees(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyWithRequirement(t, db, period.ID, "care_type", "ganztag", 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Only children, no employees
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, props)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, dp := range result.DataPoints {
		if !almostEqual(dp.AvailableHours, 0.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 0.0", i, dp.AvailableHours)
		}
		// RequiredHours should be > 0 since child has contract with matching funding
		if !almostEqual(dp.RequiredHours, 0.25*39.0, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want %v", i, dp.RequiredHours, 0.25*39.0)
		}
		if dp.ChildCount != 1 {
			t.Errorf("data point %d: ChildCount = %d, want 1", i, dp.ChildCount)
		}
		if dp.StaffCount != 0 {
			t.Errorf("data point %d: StaffCount = %d, want 0", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetStaffingHours_Empty(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should still return data points (3: Jan, Feb, Mar)
	if len(result.DataPoints) != 3 {
		t.Fatalf("expected 3 data points, got %d", len(result.DataPoints))
	}

	for i, dp := range result.DataPoints {
		if !almostEqual(dp.RequiredHours, 0.0, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want 0.0", i, dp.RequiredHours)
		}
		if !almostEqual(dp.AvailableHours, 0.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 0.0", i, dp.AvailableHours)
		}
		if dp.ChildCount != 0 {
			t.Errorf("data point %d: ChildCount = %d, want 0", i, dp.ChildCount)
		}
		if dp.StaffCount != 0 {
			t.Errorf("data point %d: StaffCount = %d, want 0", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetStaffingHours_SectionFilter(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyWithRequirement(t, db, period.ID, "care_type", "ganztag", 0.25, 0, 6)

	section1 := createTestSection(t, db, "Krippe", org.ID, false)
	section2 := createTestSection(t, db, "Elementar", org.ID, false)

	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}

	// Child in section 1
	child1 := createTestChild(t, db, "Child", "One", org.ID)
	createTestChildContract(t, db, child1.ID, contractFrom, nil, section1.ID, props)

	// Child in section 2
	child2 := createTestChild(t, db, "Child", "Two", org.ID)
	createTestChildContract(t, db, child2.ID, contractFrom, nil, section2.ID, props)

	// Employee in section 1
	payplan := createTestPayPlan(t, db, "TV-L", org.ID)
	emp1 := createTestEmployee(t, db, "Emp", "One", org.ID)
	empContract1 := &models.EmployeeContract{
		EmployeeID: emp1.ID,
		BaseContract: models.BaseContract{
			Period:    models.Period{From: contractFrom, To: nil},
			SectionID: section1.ID,
		},
		StaffCategory: "qualified",
		WeeklyHours:   30.0,
		PayPlanID:     payplan.ID,
	}
	if err := db.Create(empContract1).Error; err != nil {
		t.Fatalf("failed to create employee contract: %v", err)
	}

	// Employee in section 2
	emp2 := createTestEmployee(t, db, "Emp", "Two", org.ID)
	empContract2 := &models.EmployeeContract{
		EmployeeID: emp2.ID,
		BaseContract: models.BaseContract{
			Period:    models.Period{From: contractFrom, To: nil},
			SectionID: section2.ID,
		},
		StaffCategory: "qualified",
		WeeklyHours:   25.0,
		PayPlanID:     payplan.ID,
	}
	if err := db.Create(empContract2).Error; err != nil {
		t.Fatalf("failed to create employee contract: %v", err)
	}

	// Filter by section 1
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, &section1.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, dp := range result.DataPoints {
		// Only 1 child in section 1
		if dp.ChildCount != 1 {
			t.Errorf("data point %d: ChildCount = %d, want 1", i, dp.ChildCount)
		}
		// Only 1 employee in section 1, 30 hours
		if dp.StaffCount != 1 {
			t.Errorf("data point %d: StaffCount = %d, want 1", i, dp.StaffCount)
		}
		if !almostEqual(dp.AvailableHours, 30.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 30.0", i, dp.AvailableHours)
		}
		// 1 child * 0.25 * 39.0 = 9.75
		if !almostEqual(dp.RequiredHours, 9.75, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want 9.75", i, dp.RequiredHours)
		}
	}
}

func TestStatisticsService_GetStaffingHours_CustomDateRange(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	// Range from 2024-03 to 2024-08 should give 6 data points
	from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 6 data points: Mar, Apr, May, Jun, Jul, Aug
	if len(result.DataPoints) != 6 {
		t.Fatalf("expected 6 data points, got %d", len(result.DataPoints))
	}

	// Verify first and last dates
	if result.DataPoints[0].Date != "2024-03-01" {
		t.Errorf("first data point date = %v, want 2024-03-01", result.DataPoints[0].Date)
	}
	if result.DataPoints[5].Date != "2024-08-01" {
		t.Errorf("last data point date = %v, want 2024-08-01", result.DataPoints[5].Date)
	}
}

func TestStatisticsService_GetStaffingHours_DefaultDateRange(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	// Pass nil for from/to to use default range: 12 months back to 6 months forward = 19 data points
	result, err := svc.GetStaffingHours(ctx, org.ID, nil, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Default is 12 months back + current month + 6 months forward = 19 data points
	if len(result.DataPoints) != 19 {
		t.Errorf("expected 19 data points, got %d", len(result.DataPoints))
	}
}

func TestStatisticsService_GetStaffingHours_ContractStartsMidRange(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyWithRequirement(t, db, period.ID, "care_type", "ganztag", 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Child contract starts 2024-03-01 (mid-range)
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, props)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 6 {
		t.Fatalf("expected 6 data points, got %d", len(result.DataPoints))
	}

	// Jan and Feb: child not yet active
	for i := 0; i < 2; i++ {
		if result.DataPoints[i].ChildCount != 0 {
			t.Errorf("data point %d: ChildCount = %d, want 0 (contract not started)", i, result.DataPoints[i].ChildCount)
		}
		if !almostEqual(result.DataPoints[i].RequiredHours, 0.0, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want 0.0", i, result.DataPoints[i].RequiredHours)
		}
	}

	// Mar through Jun: child active
	for i := 2; i < 6; i++ {
		if result.DataPoints[i].ChildCount != 1 {
			t.Errorf("data point %d: ChildCount = %d, want 1", i, result.DataPoints[i].ChildCount)
		}
		if !almostEqual(result.DataPoints[i].RequiredHours, 0.25*39.0, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want %v", i, result.DataPoints[i].RequiredHours, 0.25*39.0)
		}
	}
}

func TestStatisticsService_GetStaffingHours_OngoingContracts(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyWithRequirement(t, db, period.ID, "care_type", "ganztag", 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Contracts with To = nil (ongoing)
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, props)

	payplan := createTestPayPlan(t, db, "TV-L", org.ID)
	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 35.0, "qualified", section.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// All data points should have the child and employee active
	for i, dp := range result.DataPoints {
		if dp.ChildCount != 1 {
			t.Errorf("data point %d: ChildCount = %d, want 1", i, dp.ChildCount)
		}
		if dp.StaffCount != 1 {
			t.Errorf("data point %d: StaffCount = %d, want 1", i, dp.StaffCount)
		}
		if !almostEqual(dp.AvailableHours, 35.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 35.0", i, dp.AvailableHours)
		}
	}
}

func TestStatisticsService_GetStaffingHours_NonPedagogicalExcluded(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	payplan := createTestPayPlan(t, db, "TV-L", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Qualified employee (should be counted)
	emp1 := createTestEmployee(t, db, "Emp", "Qualified", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp1.ID, payplan.ID, contractFrom, nil, 30.0, "qualified", section.ID)

	// Supplementary employee (should be counted)
	emp2 := createTestEmployee(t, db, "Emp", "Supplementary", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp2.ID, payplan.ID, contractFrom, nil, 20.0, "supplementary", section.ID)

	// Non-pedagogical employee (should NOT be counted)
	emp3 := createTestEmployee(t, db, "Emp", "Kitchen", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp3.ID, payplan.ID, contractFrom, nil, 40.0, "non_pedagogical", section.ID)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, dp := range result.DataPoints {
		// Only qualified + supplementary = 30 + 20 = 50 hours
		if !almostEqual(dp.AvailableHours, 50.0, 0.01) {
			t.Errorf("data point %d: AvailableHours = %v, want 50.0", i, dp.AvailableHours)
		}
		// Only 2 pedagogical staff counted
		if dp.StaffCount != 2 {
			t.Errorf("data point %d: StaffCount = %d, want 2", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetStaffingHours_NoFundingForState(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	// Org with state "hamburg" - no funding exists for this state
	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "hamburg")

	section := getDefaultSection(t, db, org.ID)

	// Child with contract
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, props)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, dp := range result.DataPoints {
		// No funding for state, so required hours = 0
		if !almostEqual(dp.RequiredHours, 0.0, 0.01) {
			t.Errorf("data point %d: RequiredHours = %v, want 0.0", i, dp.RequiredHours)
		}
		// Child still counted
		if dp.ChildCount != 1 {
			t.Errorf("data point %d: ChildCount = %d, want 1", i, dp.ChildCount)
		}
	}
}

// ============================================================
// GetFinancials tests
// ============================================================

func TestStatisticsService_GetFinancials_Basic(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	// Government funding: care_type=ganztag -> 166847 cents payment
	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	fundingPeriod := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyFull(t, db, fundingPeriod.ID, "care_type", "ganztag", 166847, 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// 1 child with ganztag contract
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "ganztag"})

	// Pay plan: S8a step 3, 350000 cents/month, 39h full-time, 22% employer contrib
	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 2200)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 3, 350000, nil)

	// 1 employee: S8a step 3, 39h (full-time)
	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 3})

	// 1 cost: rent 150000 cents/month
	cost := createTestCost(t, db, "Rent", org.ID)
	createTestCostEntry(t, db, cost.ID, contractFrom, nil, 150000, "Monthly rent")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 3 {
		t.Fatalf("expected 3 data points, got %d", len(result.DataPoints))
	}

	for i, dp := range result.DataPoints {
		// Income: 1 child * 166847 cents
		if dp.FundingIncome != 166847 {
			t.Errorf("dp %d: FundingIncome = %d, want 166847", i, dp.FundingIncome)
		}
		// Gross salary: 350000 * 39/39 = 350000
		if dp.GrossSalary != 350000 {
			t.Errorf("dp %d: GrossSalary = %d, want 350000", i, dp.GrossSalary)
		}
		// Employer costs: 350000 * 2200/10000 = 77000
		if dp.EmployerCosts != 77000 {
			t.Errorf("dp %d: EmployerCosts = %d, want 77000", i, dp.EmployerCosts)
		}
		// Operating cost: 150000
		if dp.OperatingCost != 150000 {
			t.Errorf("dp %d: OperatingCost = %d, want 150000", i, dp.OperatingCost)
		}
		// Totals
		if dp.TotalIncome != 166847 {
			t.Errorf("dp %d: TotalIncome = %d, want 166847", i, dp.TotalIncome)
		}
		expectedExpenses := 350000 + 77000 + 150000
		if dp.TotalExpenses != expectedExpenses {
			t.Errorf("dp %d: TotalExpenses = %d, want %d", i, dp.TotalExpenses, expectedExpenses)
		}
		expectedBalance := 166847 - expectedExpenses
		if dp.Balance != expectedBalance {
			t.Errorf("dp %d: Balance = %d, want %d", i, dp.Balance, expectedBalance)
		}
		if dp.ChildCount != 1 {
			t.Errorf("dp %d: ChildCount = %d, want 1", i, dp.ChildCount)
		}
		if dp.StaffCount != 1 {
			t.Errorf("dp %d: StaffCount = %d, want 1", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetFinancials_Empty(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 3 {
		t.Fatalf("expected 3 data points, got %d", len(result.DataPoints))
	}

	for i, dp := range result.DataPoints {
		if dp.FundingIncome != 0 {
			t.Errorf("dp %d: FundingIncome = %d, want 0", i, dp.FundingIncome)
		}
		if dp.GrossSalary != 0 {
			t.Errorf("dp %d: GrossSalary = %d, want 0", i, dp.GrossSalary)
		}
		if dp.EmployerCosts != 0 {
			t.Errorf("dp %d: EmployerCosts = %d, want 0", i, dp.EmployerCosts)
		}
		if dp.OperatingCost != 0 {
			t.Errorf("dp %d: OperatingCost = %d, want 0", i, dp.OperatingCost)
		}
		if dp.TotalIncome != 0 {
			t.Errorf("dp %d: TotalIncome = %d, want 0", i, dp.TotalIncome)
		}
		if dp.TotalExpenses != 0 {
			t.Errorf("dp %d: TotalExpenses = %d, want 0", i, dp.TotalExpenses)
		}
		if dp.Balance != 0 {
			t.Errorf("dp %d: Balance = %d, want 0", i, dp.Balance)
		}
		if dp.ChildCount != 0 {
			t.Errorf("dp %d: ChildCount = %d, want 0", i, dp.ChildCount)
		}
		if dp.StaffCount != 0 {
			t.Errorf("dp %d: StaffCount = %d, want 0", i, dp.StaffCount)
		}
	}
}

func TestStatisticsService_GetFinancials_ProRataSalary(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	// Pay plan: 39h full-time, 350000 cents/month at S8a step 3, no employer contrib
	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 0)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 3, 350000, nil)

	// Employee works 30h/week (part-time)
	emp := createTestEmployee(t, db, "Emp", "PartTime", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 30.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 3})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 1 {
		t.Fatalf("expected 1 data point, got %d", len(result.DataPoints))
	}

	dp := result.DataPoints[0]
	// Pro-rated: 350000 * 30/39 = 269231 (rounded)
	expected := int(math.Round(350000.0 * 30.0 / 39.0))
	if dp.GrossSalary != expected {
		t.Errorf("GrossSalary = %d, want %d", dp.GrossSalary, expected)
	}
}

func TestStatisticsService_GetFinancials_EmployerContribution(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	// Pay plan: 39h, 400000 cents/month, 22.50% employer contribution (2250 hundredths)
	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 2250)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S11b", 5, 400000, nil)

	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp.ID).Updates(map[string]interface{}{"grade": "S11b", "step": 5})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	if dp.GrossSalary != 400000 {
		t.Errorf("GrossSalary = %d, want 400000", dp.GrossSalary)
	}
	// Employer: 400000 * 2250/10000 = 90000
	expectedContrib := int(math.Round(400000.0 * 2250.0 / 10000.0))
	if dp.EmployerCosts != expectedContrib {
		t.Errorf("EmployerCosts = %d, want %d", dp.EmployerCosts, expectedContrib)
	}
}

func TestStatisticsService_GetFinancials_MissingPayPlanEntry(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	// Pay plan with S8a step 3, but employee is S9 step 1 (no matching entry)
	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 2200)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 3, 350000, nil)

	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp.ID).Updates(map[string]interface{}{"grade": "S9", "step": 1})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	// No matching entry -> salary should be 0, but staff still counted
	if dp.GrossSalary != 0 {
		t.Errorf("GrossSalary = %d, want 0 (no matching pay plan entry)", dp.GrossSalary)
	}
	if dp.EmployerCosts != 0 {
		t.Errorf("EmployerCosts = %d, want 0", dp.EmployerCosts)
	}
	if dp.StaffCount != 1 {
		t.Errorf("StaffCount = %d, want 1 (employee still counted)", dp.StaffCount)
	}
}

func TestStatisticsService_GetFinancials_NoPayPlanPeriodForDate(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	// Pay plan period covers only 2025, but we query 2024
	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 2200)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 3, 350000, nil)

	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 3})

	from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	if dp.GrossSalary != 0 {
		t.Errorf("GrossSalary = %d, want 0 (no pay plan period for 2024)", dp.GrossSalary)
	}
	if dp.StaffCount != 1 {
		t.Errorf("StaffCount = %d, want 1", dp.StaffCount)
	}
}

func TestStatisticsService_GetFinancials_NoFundingForState(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	// Org with state "hamburg" - no funding exists
	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "hamburg")

	section := getDefaultSection(t, db, org.ID)

	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "ganztag"})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	if dp.FundingIncome != 0 {
		t.Errorf("FundingIncome = %d, want 0 (no funding for state)", dp.FundingIncome)
	}
	if dp.ChildCount != 1 {
		t.Errorf("ChildCount = %d, want 1 (child still counted)", dp.ChildCount)
	}
}

func TestStatisticsService_GetFinancials_AllStaffIncluded(t *testing.T) {
	// Financials should include ALL staff categories (not just pedagogical)
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 0)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 3, 300000, nil)

	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Qualified staff
	emp1 := createTestEmployee(t, db, "Emp", "Qualified", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp1.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp1.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 3})

	// Non-pedagogical staff (kitchen, admin, etc.)
	emp2 := createTestEmployee(t, db, "Emp", "Kitchen", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp2.ID, payplan.ID, contractFrom, nil, 39.0, "non_pedagogical", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp2.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 3})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	// Both employees should contribute to salary
	if dp.GrossSalary != 600000 {
		t.Errorf("GrossSalary = %d, want 600000 (2 employees * 300000)", dp.GrossSalary)
	}
	if dp.StaffCount != 2 {
		t.Errorf("StaffCount = %d, want 2 (all staff categories)", dp.StaffCount)
	}
}

func TestStatisticsService_GetFinancials_MultipleCosts(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Multiple cost categories
	rent := createTestCost(t, db, "Rent", org.ID)
	createTestCostEntry(t, db, rent.ID, contractFrom, nil, 200000, "")

	insurance := createTestCost(t, db, "Insurance", org.ID)
	createTestCostEntry(t, db, insurance.ID, contractFrom, nil, 50000, "")

	utilities := createTestCost(t, db, "Utilities", org.ID)
	createTestCostEntry(t, db, utilities.ID, contractFrom, nil, 30000, "")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	expectedCost := 200000 + 50000 + 30000
	if dp.OperatingCost != expectedCost {
		t.Errorf("OperatingCost = %d, want %d", dp.OperatingCost, expectedCost)
	}
}

func TestStatisticsService_GetFinancials_CostEntryNotActiveOnDate(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	// Cost entry covers only Jul-Dec 2024
	cost := createTestCost(t, db, "Rent", org.ID)
	entryFrom := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	entryTo := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	createTestCostEntry(t, db, cost.ID, entryFrom, &entryTo, 150000, "")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 12 {
		t.Fatalf("expected 12 data points, got %d", len(result.DataPoints))
	}

	// Jan-Jun: no cost
	for i := 0; i < 6; i++ {
		if result.DataPoints[i].OperatingCost != 0 {
			t.Errorf("dp %d: OperatingCost = %d, want 0 (entry not active yet)", i, result.DataPoints[i].OperatingCost)
		}
	}
	// Jul-Dec: cost active
	for i := 6; i < 12; i++ {
		if result.DataPoints[i].OperatingCost != 150000 {
			t.Errorf("dp %d: OperatingCost = %d, want 150000", i, result.DataPoints[i].OperatingCost)
		}
	}
}

func TestStatisticsService_GetFinancials_ContractStartsMidRange(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	fundingPeriod := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyFull(t, db, fundingPeriod.ID, "care_type", "ganztag", 100000, 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Child contract starts March (mid-range)
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "ganztag"})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Jan, Feb: no child income
	for i := 0; i < 2; i++ {
		if result.DataPoints[i].FundingIncome != 0 {
			t.Errorf("dp %d: FundingIncome = %d, want 0 (contract not started)", i, result.DataPoints[i].FundingIncome)
		}
		if result.DataPoints[i].ChildCount != 0 {
			t.Errorf("dp %d: ChildCount = %d, want 0", i, result.DataPoints[i].ChildCount)
		}
	}
	// Mar-Jun: child active
	for i := 2; i < 6; i++ {
		if result.DataPoints[i].FundingIncome != 100000 {
			t.Errorf("dp %d: FundingIncome = %d, want 100000", i, result.DataPoints[i].FundingIncome)
		}
		if result.DataPoints[i].ChildCount != 1 {
			t.Errorf("dp %d: ChildCount = %d, want 1", i, result.DataPoints[i].ChildCount)
		}
	}
}

func TestStatisticsService_GetFinancials_MultipleChildren(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	fundingPeriod := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	createTestFundingPropertyFull(t, db, fundingPeriod.ID, "care_type", "ganztag", 80000, 0.25, 0, 6)
	createTestFundingPropertyFull(t, db, fundingPeriod.ID, "care_type", "halbtag", 40000, 0.12, 0, 6)

	section := getDefaultSection(t, db, org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// 2 ganztag children
	child1 := createTestChild(t, db, "Child", "One", org.ID)
	createTestChildContract(t, db, child1.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "ganztag"})
	child2 := createTestChild(t, db, "Child", "Two", org.ID)
	createTestChildContract(t, db, child2.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "ganztag"})

	// 1 halbtag child
	child3 := createTestChild(t, db, "Child", "Three", org.ID)
	createTestChildContract(t, db, child3.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "halbtag"})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	expectedIncome := 2*80000 + 1*40000
	if dp.FundingIncome != expectedIncome {
		t.Errorf("FundingIncome = %d, want %d", dp.FundingIncome, expectedIncome)
	}
	if dp.ChildCount != 3 {
		t.Errorf("ChildCount = %d, want 3", dp.ChildCount)
	}
}

func TestStatisticsService_GetFinancials_MultipleEmployees(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	payplan := createTestPayPlan(t, db, "TVöD-SuE", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, 39.0, 2000)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 3, 300000, nil)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S11b", 5, 450000, nil)

	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Employee 1: S8a step 3, 39h
	emp1 := createTestEmployee(t, db, "Emp", "One", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp1.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp1.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 3})

	// Employee 2: S11b step 5, 20h (part-time)
	emp2 := createTestEmployee(t, db, "Emp", "Two", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp2.ID, payplan.ID, contractFrom, nil, 20.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp2.ID).Updates(map[string]interface{}{"grade": "S11b", "step": 5})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	// Emp1: 300000 * 39/39 = 300000
	// Emp2: 450000 * 20/39 = 230769 (rounded)
	expectedEmp2 := int(math.Round(450000.0 * 20.0 / 39.0))
	expectedGross := 300000 + expectedEmp2
	if dp.GrossSalary != expectedGross {
		t.Errorf("GrossSalary = %d, want %d", dp.GrossSalary, expectedGross)
	}
	// Employer: each gross * 2000/10000
	expectedContrib := int(math.Round(300000.0*2000.0/10000.0)) + int(math.Round(float64(expectedEmp2)*2000.0/10000.0))
	if dp.EmployerCosts != expectedContrib {
		t.Errorf("EmployerCosts = %d, want %d", dp.EmployerCosts, expectedContrib)
	}
	if dp.StaffCount != 2 {
		t.Errorf("StaffCount = %d, want 2", dp.StaffCount)
	}
}

func TestStatisticsService_GetFinancials_BalanceNegative(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// No children (no income), but has costs
	cost := createTestCost(t, db, "Rent", org.ID)
	createTestCostEntry(t, db, cost.ID, contractFrom, nil, 500000, "")

	// Employee with salary
	payplan := createTestPayPlan(t, db, "TVöD", org.ID)
	ppPeriod := createTestPayPlanPeriodWithContrib(t, db, payplan.ID, contractFrom, nil, 39.0, 0)
	createTestPayPlanEntry(t, db, ppPeriod.ID, "S8a", 1, 250000, nil)

	emp := createTestEmployee(t, db, "Emp", "One", org.ID)
	createTestEmployeeContractWithCategory(t, db, emp.ID, payplan.ID, contractFrom, nil, 39.0, "qualified", section.ID)
	db.Model(&models.EmployeeContract{}).Where("employee_id = ?", emp.ID).Updates(map[string]interface{}{"grade": "S8a", "step": 1})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	if dp.TotalIncome != 0 {
		t.Errorf("TotalIncome = %d, want 0", dp.TotalIncome)
	}
	expectedExpenses := 250000 + 500000
	if dp.TotalExpenses != expectedExpenses {
		t.Errorf("TotalExpenses = %d, want %d", dp.TotalExpenses, expectedExpenses)
	}
	if dp.Balance >= 0 {
		t.Errorf("Balance = %d, want negative", dp.Balance)
	}
	if dp.Balance != -expectedExpenses {
		t.Errorf("Balance = %d, want %d", dp.Balance, -expectedExpenses)
	}
}

func TestStatisticsService_GetFinancials_DefaultDateRange(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	// nil from/to -> default 12 months back + 6 months forward = 19 data points
	result, err := svc.GetFinancials(ctx, org.ID, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 19 {
		t.Errorf("expected 19 data points (default range), got %d", len(result.DataPoints))
	}
}

func TestStatisticsService_GetFinancials_UnmatchedChildProperty(t *testing.T) {
	// Child has a property not matching any funding property -> no income from it
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")
	toDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	fundingPeriod := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate, 39.0)
	// Funding only covers "ganztag"
	createTestFundingPropertyFull(t, db, fundingPeriod.ID, "care_type", "ganztag", 100000, 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Child has "halbtag" - won't match
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, models.ContractProperties{"care_type": "halbtag"})

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	if dp.FundingIncome != 0 {
		t.Errorf("FundingIncome = %d, want 0 (no matching funding property)", dp.FundingIncome)
	}
	if dp.ChildCount != 1 {
		t.Errorf("ChildCount = %d, want 1 (child still counted)", dp.ChildCount)
	}
}

func TestStatisticsService_GetFinancials_EmployeeNoPayPlan(t *testing.T) {
	// Employee with PayPlanID=0 should not cause errors
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	section := getDefaultSection(t, db, org.ID)

	// Create employee contract with PayPlanID=0 (no pay plan assigned)
	emp := createTestEmployee(t, db, "Emp", "NoPlan", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	contract := &models.EmployeeContract{
		EmployeeID: emp.ID,
		BaseContract: models.BaseContract{
			Period:    models.Period{From: contractFrom},
			SectionID: section.ID,
		},
		StaffCategory: "qualified",
		WeeklyHours:   39.0,
		PayPlanID:     0,
	}
	if err := db.Create(contract).Error; err != nil {
		// PayPlanID 0 may fail FK constraint - that's fine, skip
		t.Skip("PayPlanID=0 not allowed by FK constraint")
	}

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	if dp.GrossSalary != 0 {
		t.Errorf("GrossSalary = %d, want 0 (no pay plan)", dp.GrossSalary)
	}
}

func TestStatisticsService_GetFinancials_OneCostEntryPerCategory(t *testing.T) {
	// Only one active entry per cost category should be counted (break after first match)
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	cost := createTestCost(t, db, "Rent", org.ID)

	// Two overlapping entries for same cost (shouldn't happen in practice, but test the break)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createTestCostEntry(t, db, cost.ID, contractFrom, nil, 100000, "Old rate")
	// Second entry also active - only first should be counted due to break
	createTestCostEntry(t, db, cost.ID, contractFrom, nil, 200000, "New rate")

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetFinancials(ctx, org.ID, &from, &to)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dp := result.DataPoints[0]
	// Should NOT be 300000 (both added) - only one should be used
	if dp.OperatingCost == 300000 {
		t.Errorf("OperatingCost = 300000, but should only count one entry per cost category")
	}
	// Should be either 100000 or 200000 (whichever DB returns first)
	if dp.OperatingCost != 100000 && dp.OperatingCost != 200000 {
		t.Errorf("OperatingCost = %d, want 100000 or 200000", dp.OperatingCost)
	}
}

func TestStatisticsService_GetStaffingHours_FundingPeriodChange(t *testing.T) {
	db := setupTestDB(t)
	svc := createStatisticsService(db)
	ctx := context.Background()

	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	funding := createTestGovernmentFunding(t, db, "Berlin Funding")

	// Period 1: Jan-Jun with 39 hours
	toDate1 := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	period1 := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), &toDate1, 39.0)
	createTestFundingPropertyWithRequirement(t, db, period1.ID, "care_type", "ganztag", 0.25, 0, 6)

	// Period 2: Jul-Dec with 40 hours
	toDate2 := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	period2 := createTestFundingPeriod(t, db, funding.ID, time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), &toDate2, 40.0)
	createTestFundingPropertyWithRequirement(t, db, period2.ID, "care_type", "ganztag", 0.25, 0, 6)

	section := getDefaultSection(t, db, org.ID)

	// Child with ongoing contract
	child := createTestChild(t, db, "Child", "One", org.ID)
	contractFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	props := models.ContractProperties{"care_type": "ganztag"}
	createTestChildContract(t, db, child.ID, contractFrom, nil, section.ID, props)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	result, err := svc.GetStaffingHours(ctx, org.ID, &from, &to, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.DataPoints) != 12 {
		t.Fatalf("expected 12 data points, got %d", len(result.DataPoints))
	}

	// Jan-Jun: 1 child * 0.25 * 39.0 = 9.75
	expectedFirst := 0.25 * 39.0
	for i := 0; i < 6; i++ {
		if !almostEqual(result.DataPoints[i].RequiredHours, expectedFirst, 0.01) {
			t.Errorf("data point %d (period 1): RequiredHours = %v, want %v", i, result.DataPoints[i].RequiredHours, expectedFirst)
		}
	}

	// Jul-Dec: 1 child * 0.25 * 40.0 = 10.0
	expectedSecond := 0.25 * 40.0
	for i := 6; i < 12; i++ {
		if !almostEqual(result.DataPoints[i].RequiredHours, expectedSecond, 0.01) {
			t.Errorf("data point %d (period 2): RequiredHours = %v, want %v", i, result.DataPoints[i].RequiredHours, expectedSecond)
		}
	}
}
