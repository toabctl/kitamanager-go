package service

import (
	"context"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

// UserOrganizationService handles business logic for user-organization-role operations
type UserOrganizationService struct {
	userOrgStore store.UserOrganizationStorer
	userStore    store.UserStorer
	transactor   store.Transactor
}

// NewUserOrganizationService creates a new UserOrganizationService
func NewUserOrganizationService(userOrgStore store.UserOrganizationStorer, userStore store.UserStorer, transactor store.Transactor) *UserOrganizationService {
	return &UserOrganizationService{
		userOrgStore: userOrgStore,
		userStore:    userStore,
		transactor:   transactor,
	}
}

// AddUserToOrganization adds a user to an organization with a specific role.
// requesterID is the user performing the operation (for authorization check).
func (s *UserOrganizationService) AddUserToOrganization(ctx context.Context, userID, orgID uint, role models.Role, createdBy string, requesterID uint) (*models.UserOrganizationResponse, error) {
	// Validate role
	if !role.IsValid() {
		return nil, apperror.BadRequest("invalid role: must be admin, manager, or member")
	}

	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return nil, classifyStoreError(err, "user")
	}

	// Verify requester has admin access to the organization
	if err := s.verifyRequesterOrgAccess(ctx, requesterID, orgID); err != nil {
		return nil, err
	}

	// Check + create in a single transaction to prevent race conditions
	var uo *models.UserOrganization
	if err := s.transactor.InTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.userOrgStore.Exists(txCtx, userID, orgID)
		if err != nil {
			return apperror.InternalWrap(err, "failed to check existing membership")
		}
		if exists {
			return apperror.BadRequest("user is already a member of this organization")
		}

		uo, err = s.userOrgStore.AddUserToOrg(txCtx, userID, orgID, role, createdBy)
		if err != nil {
			return apperror.InternalWrap(err, "failed to add user to organization")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	resp := uo.ToResponse()
	return &resp, nil
}

// UpdateUserOrganizationRole updates a user's role in an organization.
// requesterID is the user performing the operation (for authorization check).
func (s *UserOrganizationService) UpdateUserOrganizationRole(ctx context.Context, userID, orgID uint, role models.Role, requesterID uint) (*models.UserOrganizationResponse, error) {
	// Validate role
	if !role.IsValid() {
		return nil, apperror.BadRequest("invalid role: must be admin, manager, or member")
	}

	// Verify membership exists
	uo, err := s.userOrgStore.FindByUserAndOrg(ctx, userID, orgID)
	if err != nil {
		return nil, classifyStoreError(err, "user-organization membership")
	}

	// Verify requester has admin access to the organization
	if err := s.verifyRequesterOrgAccess(ctx, requesterID, orgID); err != nil {
		return nil, err
	}

	// Update role
	if err := s.userOrgStore.UpdateRole(ctx, userID, orgID, role); err != nil {
		return nil, apperror.InternalWrap(err, "failed to update role")
	}

	uo.Role = role
	resp := uo.ToResponse()
	return &resp, nil
}

// RemoveUserFromOrganization removes a user from an organization.
// requesterID is the user performing the operation (for authorization check).
func (s *UserOrganizationService) RemoveUserFromOrganization(ctx context.Context, userID, orgID uint, requesterID uint) error {
	// Verify requester has admin access to the organization
	if err := s.verifyRequesterOrgAccess(ctx, requesterID, orgID); err != nil {
		return err
	}

	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return classifyStoreError(err, "user")
	}

	// Check + delete in a single transaction to prevent race conditions
	return s.transactor.InTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.userOrgStore.Exists(txCtx, userID, orgID)
		if err != nil {
			return apperror.InternalWrap(err, "failed to check membership")
		}
		if !exists {
			return apperror.NotFound("user-organization membership")
		}

		if err := s.userOrgStore.RemoveUserFromOrg(txCtx, userID, orgID); err != nil {
			return apperror.InternalWrap(err, "failed to remove user from organization")
		}
		return nil
	})
}

// GetUserMemberships returns all organization memberships for a user
func (s *UserOrganizationService) GetUserMemberships(ctx context.Context, userID uint) (*models.UserMembershipsResponse, error) {
	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return nil, classifyStoreError(err, "user")
	}

	memberships, err := s.userOrgStore.FindByUser(ctx, userID)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to fetch memberships")
	}

	result := make([]models.UserMembership, 0, len(memberships))
	for _, m := range memberships {
		result = append(result, models.UserMembership{
			UserID:         m.UserID,
			OrganizationID: m.OrganizationID,
			Role:           m.Role,
			Organization:   m.Organization,
		})
	}

	return &models.UserMembershipsResponse{Memberships: result}, nil
}

// SetSuperAdmin sets or unsets superadmin status for a user
func (s *UserOrganizationService) SetSuperAdmin(ctx context.Context, userID uint, isSuperAdmin bool) error {
	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return classifyStoreError(err, "user")
	}

	if err := s.userOrgStore.SetSuperAdmin(ctx, userID, isSuperAdmin); err != nil {
		return apperror.InternalWrap(err, "failed to update superadmin status")
	}
	return nil
}

// verifyRequesterOrgAccess checks that the requester is a superadmin or has admin role
// in the given organization. Returns apperror.Forbidden if not authorized.
func (s *UserOrganizationService) verifyRequesterOrgAccess(ctx context.Context, requesterID, orgID uint) error {
	isSuperAdmin, err := s.userOrgStore.IsSuperAdmin(ctx, requesterID)
	if err != nil {
		return apperror.InternalWrap(err, "failed to check superadmin status")
	}
	if isSuperAdmin {
		return nil
	}

	role, err := s.userOrgStore.GetRoleInOrg(ctx, requesterID, orgID)
	if err != nil {
		return apperror.InternalWrap(err, "failed to check organization access")
	}
	if role != models.RoleAdmin {
		return apperror.Forbidden("insufficient permissions for this organization")
	}
	return nil
}
