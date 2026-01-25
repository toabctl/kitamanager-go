package service

import (
	"context"
	"errors"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

// EmployeeService handles business logic for employee operations
type EmployeeService struct {
	store store.EmployeeStorer
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(store store.EmployeeStorer) *EmployeeService {
	return &EmployeeService{store: store}
}

// List returns a paginated list of employees
func (s *EmployeeService) List(ctx context.Context, limit, offset int) ([]models.Employee, int64, error) {
	employees, total, err := s.store.FindAll(limit, offset)
	if err != nil {
		return nil, 0, apperror.Internal("failed to fetch employees")
	}
	return employees, total, nil
}

// ListByOrganization returns a paginated list of employees for an organization
func (s *EmployeeService) ListByOrganization(ctx context.Context, orgID uint, limit, offset int) ([]models.Employee, int64, error) {
	employees, total, err := s.store.FindByOrganization(orgID, limit, offset)
	if err != nil {
		return nil, 0, apperror.Internal("failed to fetch employees")
	}
	return employees, total, nil
}

// GetByID returns an employee by ID
func (s *EmployeeService) GetByID(ctx context.Context, id uint) (*models.Employee, error) {
	employee, err := s.store.FindByID(id)
	if err != nil {
		return nil, apperror.NotFound("employee")
	}
	return employee, nil
}

// Create creates a new employee
func (s *EmployeeService) Create(ctx context.Context, req *models.EmployeeCreate) (*models.Employee, error) {
	employee := &models.Employee{
		Person: models.Person{
			OrganizationID: req.OrganizationID,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Birthdate:      req.Birthdate,
		},
	}

	if err := s.store.Create(employee); err != nil {
		return nil, apperror.Internal("failed to create employee")
	}

	return employee, nil
}

// Update updates an existing employee
func (s *EmployeeService) Update(ctx context.Context, id uint, req *models.EmployeeUpdate) (*models.Employee, error) {
	employee, err := s.store.FindByID(id)
	if err != nil {
		return nil, apperror.NotFound("employee")
	}

	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.Birthdate != nil {
		employee.Birthdate = *req.Birthdate
	}

	if err := s.store.Update(employee); err != nil {
		return nil, apperror.Internal("failed to update employee")
	}

	return employee, nil
}

// Delete deletes an employee
func (s *EmployeeService) Delete(ctx context.Context, id uint) error {
	if err := s.store.Delete(id); err != nil {
		return apperror.Internal("failed to delete employee")
	}
	return nil
}

// ListContracts returns contract history for an employee
func (s *EmployeeService) ListContracts(ctx context.Context, employeeID uint) ([]models.EmployeeContract, error) {
	// Verify employee exists
	_, err := s.store.FindByID(employeeID)
	if err != nil {
		return nil, apperror.NotFound("employee")
	}

	contracts, err := s.store.Contracts().GetHistory(employeeID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch contracts")
	}
	return contracts, nil
}

// GetCurrentContract returns the current active contract for an employee
func (s *EmployeeService) GetCurrentContract(ctx context.Context, employeeID uint) (*models.EmployeeContract, error) {
	contract, err := s.store.Contracts().GetCurrentContract(employeeID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch contract")
	}
	if contract == nil {
		return nil, apperror.NotFound("active contract")
	}
	return contract, nil
}

// CreateContract creates a new contract for an employee
func (s *EmployeeService) CreateContract(ctx context.Context, employeeID uint, req *models.EmployeeContractCreate) (*models.EmployeeContract, error) {
	// Verify employee exists
	_, err := s.store.FindByID(employeeID)
	if err != nil {
		return nil, apperror.NotFound("employee")
	}

	// Validate no overlap
	if err := s.store.Contracts().ValidateNoOverlap(employeeID, req.From, req.To, nil); err != nil {
		if errors.Is(err, store.ErrContractOverlap) {
			return nil, apperror.Conflict(err.Error())
		}
		return nil, apperror.Internal("failed to validate contract")
	}

	contract := &models.EmployeeContract{
		EmployeeID: employeeID,
		Period: models.Period{
			From: req.From,
			To:   req.To,
		},
		Position:    req.Position,
		WeeklyHours: req.WeeklyHours,
		Salary:      req.Salary,
	}

	if err := s.store.CreateContract(contract); err != nil {
		return nil, apperror.Internal("failed to create contract")
	}

	return contract, nil
}

// DeleteContract deletes a contract
func (s *EmployeeService) DeleteContract(ctx context.Context, contractID uint) error {
	if err := s.store.DeleteContract(contractID); err != nil {
		return apperror.Internal("failed to delete contract")
	}
	return nil
}
