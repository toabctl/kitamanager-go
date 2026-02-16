package service

import (
	"context"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
	"github.com/eenemeene/kitamanager-go/internal/validation"
)

// EmployeeService handles business logic for employee operations
type EmployeeService struct {
	store        store.EmployeeStorer
	payPlanStore store.PayPlanStorer
	sectionStore store.SectionStorer
	transactor   store.Transactor
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(store store.EmployeeStorer, payPlanStore store.PayPlanStorer, sectionStore store.SectionStorer, transactor store.Transactor) *EmployeeService {
	return &EmployeeService{store: store, payPlanStore: payPlanStore, sectionStore: sectionStore, transactor: transactor}
}

// List returns a paginated list of employees
func (s *EmployeeService) List(ctx context.Context, limit, offset int) ([]models.EmployeeResponse, int64, error) {
	return personList(ctx, s.store.FindAll, (*models.Employee).ToResponse, "employees", limit, offset)
}

// ListByOrganization returns a paginated list of employees for an organization
func (s *EmployeeService) ListByOrganization(ctx context.Context, orgID uint, limit, offset int) ([]models.EmployeeResponse, int64, error) {
	return s.ListByOrganizationAndSection(ctx, orgID, models.EmployeeListFilter{}, limit, offset)
}

// ListByOrganizationAndSection returns a paginated list of employees for an organization,
// optionally filtered by section, active contract date, name search, and/or staff category.
func (s *EmployeeService) ListByOrganizationAndSection(ctx context.Context, orgID uint, filter models.EmployeeListFilter, limit, offset int) ([]models.EmployeeResponse, int64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, apperror.BadRequest(err.Error())
	}

	employees, total, err := s.store.FindByOrganizationAndSection(ctx, orgID, filter.SectionID, filter.ActiveOn, filter.Search, filter.StaffCategory, limit, offset)
	if err != nil {
		return nil, 0, apperror.InternalWrap(err, "failed to fetch employees")
	}

	return toResponseList(employees, (*models.Employee).ToResponse), total, nil
}

// GetByID returns an employee by ID, validating it belongs to the specified organization
func (s *EmployeeService) GetByID(ctx context.Context, id, orgID uint) (*models.EmployeeResponse, error) {
	return personGetByID(ctx, s.store.FindByID, (*models.Employee).ToResponse, id, orgID, "employee")
}

// Create creates a new employee
func (s *EmployeeService) Create(ctx context.Context, orgID uint, req *models.EmployeeCreateRequest) (*models.EmployeeResponse, error) {
	return personCreate(ctx,
		&validation.PersonCreateFields{FirstName: req.FirstName, LastName: req.LastName, Gender: req.Gender, Birthdate: req.Birthdate},
		func(p models.Person) *models.Employee { return &models.Employee{Person: p} },
		s.store.Create, (*models.Employee).ToResponse, orgID, "employee")
}

// Update updates an existing employee, validating it belongs to the specified organization
func (s *EmployeeService) Update(ctx context.Context, id, orgID uint, req *models.EmployeeUpdateRequest) (*models.EmployeeResponse, error) {
	return personUpdate(ctx, s.store.FindByID, func(e *models.Employee) *models.Person { return &e.Person },
		s.store.Update, (*models.Employee).ToResponse, id, orgID,
		personUpdateFields{FirstName: req.FirstName, LastName: req.LastName, Gender: req.Gender, Birthdate: req.Birthdate},
		"employee")
}

// Delete deletes an employee and its contracts, validating it belongs to the specified organization.
// The ownership check and deletion run in a single transaction.
func (s *EmployeeService) Delete(ctx context.Context, id, orgID uint) error {
	return personDelete(ctx, s.transactor, s.store.FindByID, s.store.Delete, id, orgID, "employee")
}
