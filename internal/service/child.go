package service

import (
	"context"
	"errors"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

// ChildService handles business logic for child operations
type ChildService struct {
	store store.ChildStorer
}

// NewChildService creates a new child service
func NewChildService(store store.ChildStorer) *ChildService {
	return &ChildService{store: store}
}

// List returns a paginated list of children
func (s *ChildService) List(ctx context.Context, limit, offset int) ([]models.Child, int64, error) {
	children, total, err := s.store.FindAll(limit, offset)
	if err != nil {
		return nil, 0, apperror.Internal("failed to fetch children")
	}
	return children, total, nil
}

// ListByOrganization returns a paginated list of children for an organization
func (s *ChildService) ListByOrganization(ctx context.Context, orgID uint, limit, offset int) ([]models.Child, int64, error) {
	children, total, err := s.store.FindByOrganization(orgID, limit, offset)
	if err != nil {
		return nil, 0, apperror.Internal("failed to fetch children")
	}
	return children, total, nil
}

// GetByID returns a child by ID
func (s *ChildService) GetByID(ctx context.Context, id uint) (*models.Child, error) {
	child, err := s.store.FindByID(id)
	if err != nil {
		return nil, apperror.NotFound("child")
	}
	return child, nil
}

// Create creates a new child
func (s *ChildService) Create(ctx context.Context, req *models.ChildCreate) (*models.Child, error) {
	child := &models.Child{
		Person: models.Person{
			OrganizationID: req.OrganizationID,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Birthdate:      req.Birthdate,
		},
	}

	if err := s.store.Create(child); err != nil {
		return nil, apperror.Internal("failed to create child")
	}

	return child, nil
}

// Update updates an existing child
func (s *ChildService) Update(ctx context.Context, id uint, req *models.ChildUpdate) (*models.Child, error) {
	child, err := s.store.FindByID(id)
	if err != nil {
		return nil, apperror.NotFound("child")
	}

	if req.FirstName != nil {
		child.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		child.LastName = *req.LastName
	}
	if req.Birthdate != nil {
		child.Birthdate = *req.Birthdate
	}

	if err := s.store.Update(child); err != nil {
		return nil, apperror.Internal("failed to update child")
	}

	return child, nil
}

// Delete deletes a child
func (s *ChildService) Delete(ctx context.Context, id uint) error {
	if err := s.store.Delete(id); err != nil {
		return apperror.Internal("failed to delete child")
	}
	return nil
}

// ListContracts returns contract history for a child
func (s *ChildService) ListContracts(ctx context.Context, childID uint) ([]models.ChildContract, error) {
	// Verify child exists
	_, err := s.store.FindByID(childID)
	if err != nil {
		return nil, apperror.NotFound("child")
	}

	contracts, err := s.store.Contracts().GetHistory(childID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch contracts")
	}
	return contracts, nil
}

// GetCurrentContract returns the current active contract for a child
func (s *ChildService) GetCurrentContract(ctx context.Context, childID uint) (*models.ChildContract, error) {
	contract, err := s.store.Contracts().GetCurrentContract(childID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch contract")
	}
	if contract == nil {
		return nil, apperror.NotFound("active contract")
	}
	return contract, nil
}

// CreateContract creates a new contract for a child
func (s *ChildService) CreateContract(ctx context.Context, childID uint, req *models.ChildContractCreate) (*models.ChildContract, error) {
	// Verify child exists
	_, err := s.store.FindByID(childID)
	if err != nil {
		return nil, apperror.NotFound("child")
	}

	// Validate no overlap
	if err := s.store.Contracts().ValidateNoOverlap(childID, req.From, req.To, nil); err != nil {
		if errors.Is(err, store.ErrContractOverlap) {
			return nil, apperror.Conflict(err.Error())
		}
		return nil, apperror.Internal("failed to validate contract")
	}

	contract := &models.ChildContract{
		ChildID: childID,
		Period: models.Period{
			From: req.From,
			To:   req.To,
		},
		CareHoursPerWeek: req.CareHoursPerWeek,
		GroupID:          req.GroupID,
		MealsIncluded:    req.MealsIncluded,
		SpecialNeeds:     req.SpecialNeeds,
	}

	if err := s.store.CreateContract(contract); err != nil {
		return nil, apperror.Internal("failed to create contract")
	}

	return contract, nil
}

// DeleteContract deletes a contract
func (s *ChildService) DeleteContract(ctx context.Context, contractID uint) error {
	if err := s.store.DeleteContract(contractID); err != nil {
		return apperror.Internal("failed to delete contract")
	}
	return nil
}
