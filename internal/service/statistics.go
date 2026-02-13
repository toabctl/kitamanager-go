package service

import (
	"context"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
	"github.com/eenemeene/kitamanager-go/internal/validation"
)

// StatisticsService handles cross-resource statistics calculations
type StatisticsService struct {
	childStore    store.ChildStorer
	employeeStore store.EmployeeStorer
	orgStore      store.OrganizationStorer
	fundingStore  store.GovernmentFundingStorer
}

// NewStatisticsService creates a new statistics service
func NewStatisticsService(childStore store.ChildStorer, employeeStore store.EmployeeStorer, orgStore store.OrganizationStorer, fundingStore store.GovernmentFundingStorer) *StatisticsService {
	return &StatisticsService{
		childStore:    childStore,
		employeeStore: employeeStore,
		orgStore:      orgStore,
		fundingStore:  fundingStore,
	}
}

// pedagogicalCategories lists staff categories counted toward staffing requirements
var pedagogicalCategories = []string{
	string(models.StaffCategoryQualified),
	string(models.StaffCategorySupplementary),
}

// GetStaffingHours calculates monthly staffing hours data points
func (s *StatisticsService) GetStaffingHours(ctx context.Context, orgID uint, from, to *time.Time, sectionID *uint) (*models.StaffingHoursResponse, error) {
	// Default date range: 12 months back to 6 months forward, snapped to 1st of month
	now := time.Now()
	var rangeStart, rangeEnd time.Time
	if from != nil {
		rangeStart = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		rangeStart = time.Date(now.Year(), now.Month()-12, 1, 0, 0, 0, 0, time.UTC)
	}
	if to != nil {
		rangeEnd = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		rangeEnd = time.Date(now.Year(), now.Month()+6, 1, 0, 0, 0, 0, time.UTC)
	}

	// Fetch organization for state
	org, err := s.orgStore.FindByID(ctx, orgID)
	if err != nil {
		return nil, apperror.NotFound("organization")
	}

	// Fetch government funding with all periods and properties
	var fundingPeriods []models.GovernmentFundingPeriod
	funding, err := s.fundingStore.FindByStateWithDetails(ctx, org.State, 0, nil)
	if err == nil {
		fundingPeriods = funding.Periods
	}
	// If no funding found, fundingPeriods stays nil — required hours will be 0

	// Fetch children with contracts in range
	children, err := s.childStore.FindByOrganizationInDateRange(ctx, orgID, rangeStart, rangeEnd, sectionID)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to fetch children")
	}

	// Fetch employee contracts in range (pedagogical staff only)
	employeeContracts, err := s.employeeStore.FindContractsByOrganizationInDateRange(ctx, orgID, rangeStart, rangeEnd, pedagogicalCategories, sectionID)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to fetch employee contracts")
	}

	// Generate data points for each month
	var dataPoints []models.StaffingHoursDataPoint
	for date := rangeStart; !date.After(rangeEnd); date = date.AddDate(0, 1, 0) {
		dp := models.StaffingHoursDataPoint{
			Date: date.Format(models.DateFormat),
		}

		// Calculate required hours from children
		requiredHours := 0.0
		childCount := 0
		period := findPeriodForDate(fundingPeriods, date)
		for i := range children {
			child := &children[i]
			// Check if child has a contract active on this date
			for j := range child.Contracts {
				contract := &child.Contracts[j]
				if !contract.From.After(date) && (contract.To == nil || !contract.To.Before(date)) {
					childCount++
					if period != nil {
						age := validation.CalculateAgeOnDate(child.Birthdate, date)
						requirement := sumChildRequirement(age, contract.Properties, period)
						requiredHours += requirement * period.FullTimeWeeklyHours
					}
					break // Only count each child once per month
				}
			}
		}

		// Calculate available hours from employee contracts
		availableHours := 0.0
		staffCount := 0
		employeeSeen := make(map[uint]bool)
		for i := range employeeContracts {
			ec := &employeeContracts[i]
			if !ec.From.After(date) && (ec.To == nil || !ec.To.Before(date)) {
				availableHours += ec.WeeklyHours
				if !employeeSeen[ec.EmployeeID] {
					employeeSeen[ec.EmployeeID] = true
					staffCount++
				}
			}
		}

		dp.RequiredHours = requiredHours
		dp.AvailableHours = availableHours
		dp.ChildCount = childCount
		dp.StaffCount = staffCount
		dataPoints = append(dataPoints, dp)
	}

	return &models.StaffingHoursResponse{
		DataPoints: dataPoints,
	}, nil
}
