package service

import (
	"context"
	"log/slog"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

// UserGroupService handles business logic for user-group-role operations
type UserGroupService struct {
	userGroupStore store.UserGroupStorer
	userStore      store.UserStorer
	groupStore     store.GroupStorer
}

// NewUserGroupService creates a new UserGroupService
func NewUserGroupService(userGroupStore store.UserGroupStorer, userStore store.UserStorer, groupStore store.GroupStorer) *UserGroupService {
	return &UserGroupService{
		userGroupStore: userGroupStore,
		userStore:      userStore,
		groupStore:     groupStore,
	}
}

// AddUserToGroup adds a user to a group with a specific role
func (s *UserGroupService) AddUserToGroup(ctx context.Context, userID, groupID uint, role models.Role, createdBy string) (*models.UserGroupResponse, error) {
	// Validate role
	if !role.IsValid() {
		return nil, apperror.BadRequest("invalid role: must be admin, manager, or member")
	}

	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.NotFound("user")
	}

	// Verify group exists
	group, err := s.groupStore.FindByID(ctx, groupID)
	if err != nil {
		return nil, apperror.NotFound("group")
	}

	// Check if already a member
	exists, err := s.userGroupStore.Exists(ctx, userID, groupID)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to check existing membership")
	}
	if exists {
		return nil, apperror.BadRequest("user is already a member of this group")
	}

	// Create membership
	ug, err := s.userGroupStore.AddUserToGroup(ctx, userID, groupID, role, createdBy)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to add user to group")
	}

	// Load group data for response
	ug.Group = group
	resp := ug.ToResponse()
	return &resp, nil
}

// UpdateUserGroupRole updates a user's role in a group
func (s *UserGroupService) UpdateUserGroupRole(ctx context.Context, userID, groupID uint, role models.Role) (*models.UserGroupResponse, error) {
	// Validate role
	if !role.IsValid() {
		return nil, apperror.BadRequest("invalid role: must be admin, manager, or member")
	}

	// Verify membership exists
	ug, err := s.userGroupStore.FindByUserAndGroup(ctx, userID, groupID)
	if err != nil {
		return nil, apperror.NotFound("user-group membership")
	}

	// Update role
	if err := s.userGroupStore.UpdateRole(ctx, userID, groupID, role); err != nil {
		return nil, apperror.InternalWrap(err, "failed to update role")
	}

	// Load group for response
	group, err := s.groupStore.FindByID(ctx, groupID)
	if err != nil {
		slog.Warn("failed to load group for response", "group_id", groupID, "error", err)
	}
	ug.Role = role
	ug.Group = group
	resp := ug.ToResponse()
	return &resp, nil
}

// RemoveUserFromGroup removes a user from a group
func (s *UserGroupService) RemoveUserFromGroup(ctx context.Context, userID, groupID uint) error {
	// Check if membership exists
	exists, err := s.userGroupStore.Exists(ctx, userID, groupID)
	if err != nil {
		return apperror.InternalWrap(err, "failed to check membership")
	}
	if !exists {
		return apperror.NotFound("user-group membership")
	}

	if err := s.userGroupStore.RemoveUserFromGroup(ctx, userID, groupID); err != nil {
		return apperror.InternalWrap(err, "failed to remove user from group")
	}
	return nil
}

// GetUserMemberships returns all group memberships for a user with effective org roles
func (s *UserGroupService) GetUserMemberships(ctx context.Context, userID uint) (*models.UserMembershipsResponse, error) {
	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.NotFound("user")
	}

	memberships, err := s.userGroupStore.FindByUser(ctx, userID)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to fetch memberships")
	}

	// Calculate effective org roles
	orgRoles, err := s.userGroupStore.GetUserOrganizationsWithRoles(ctx, userID)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to calculate effective roles")
	}

	result := make([]models.UserMembership, 0, len(memberships))
	for _, m := range memberships {
		var org *models.Organization
		var effectiveRole models.Role

		if m.Group != nil {
			org = m.Group.Organization
			effectiveRole = orgRoles[m.Group.OrganizationID]
		}

		result = append(result, models.UserMembership{
			UserID:           m.UserID,
			GroupID:          m.GroupID,
			Role:             m.Role,
			Group:            m.Group,
			EffectiveOrgRole: effectiveRole,
			Organization:     org,
		})
	}

	return &models.UserMembershipsResponse{Memberships: result}, nil
}

// SetSuperAdmin sets or unsets superadmin status for a user
func (s *UserGroupService) SetSuperAdmin(ctx context.Context, userID uint, isSuperAdmin bool) error {
	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return apperror.NotFound("user")
	}

	if err := s.userGroupStore.SetSuperAdmin(ctx, userID, isSuperAdmin); err != nil {
		return apperror.InternalWrap(err, "failed to update superadmin status")
	}
	return nil
}

// AddUserToOrganization adds a user to an organization's default group with member role
func (s *UserGroupService) AddUserToOrganization(ctx context.Context, userID, orgID uint, createdBy string) (*models.UserGroupResponse, error) {
	// Find default group for organization
	defaultGroup, err := s.groupStore.FindDefaultGroup(ctx, orgID)
	if err != nil {
		return nil, apperror.NotFound("organization or default group")
	}

	return s.AddUserToGroup(ctx, userID, defaultGroup.ID, models.RoleMember, createdBy)
}

// RemoveUserFromOrganization removes a user from all groups in an organization
func (s *UserGroupService) RemoveUserFromOrganization(ctx context.Context, userID, orgID uint) error {
	// Verify user exists
	_, err := s.userStore.FindByID(ctx, userID)
	if err != nil {
		return apperror.NotFound("user")
	}

	if err := s.userGroupStore.RemoveUserFromOrg(ctx, userID, orgID); err != nil {
		return apperror.InternalWrap(err, "failed to remove user from organization")
	}
	return nil
}
