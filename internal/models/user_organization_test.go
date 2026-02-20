package models

import (
	"testing"
	"time"
)

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{"valid admin", RoleAdmin, true},
		{"valid manager", RoleManager, true},
		{"valid member", RoleMember, true},
		{"valid staff", RoleStaff, true},
		{"empty string", Role(""), false},
		{"superadmin is not valid", Role("superadmin"), false},
		{"case sensitive - Admin", Role("Admin"), false},
		{"case sensitive - ADMIN", Role("ADMIN"), false},
		{"unknown role", Role("unknown"), false},
		{"whitespace", Role(" "), false},
		{"whitespace admin", Role(" admin"), false},
		{"admin with trailing space", Role("admin "), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.IsValid()
			if got != tt.expected {
				t.Errorf("Role(%q).IsValid() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestRole_Precedence(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected int
	}{
		{"admin has precedence 3", RoleAdmin, 3},
		{"manager has precedence 2", RoleManager, 2},
		{"member has precedence 1", RoleMember, 1},
		{"staff has precedence 0", RoleStaff, 0},
		{"invalid role has precedence 0", Role("invalid"), 0},
		{"empty string has precedence 0", Role(""), 0},
		{"superadmin has precedence 0", Role("superadmin"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.Precedence()
			if got != tt.expected {
				t.Errorf("Role(%q).Precedence() = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestRole_Precedence_Ordering(t *testing.T) {
	// Verify that admin > manager > member > staff
	if RoleAdmin.Precedence() <= RoleManager.Precedence() {
		t.Error("admin should have higher precedence than manager")
	}
	if RoleManager.Precedence() <= RoleMember.Precedence() {
		t.Error("manager should have higher precedence than member")
	}
	if RoleMember.Precedence() <= RoleStaff.Precedence() {
		t.Error("member should have higher precedence than staff")
	}
}

func TestUserOrganization_ToResponse(t *testing.T) {
	now := time.Now()
	org := &Organization{
		ID:   1,
		Name: "Test Org",
	}

	uo := &UserOrganization{
		UserID:         1,
		OrganizationID: 1,
		Role:           RoleAdmin,
		CreatedAt:      now,
		CreatedBy:      "admin@example.com",
		Organization:   org,
	}

	resp := uo.ToResponse()

	if resp.UserID != 1 {
		t.Errorf("ToResponse().UserID = %d, want 1", resp.UserID)
	}
	if resp.OrganizationID != 1 {
		t.Errorf("ToResponse().OrganizationID = %d, want 1", resp.OrganizationID)
	}
	if resp.Role != RoleAdmin {
		t.Errorf("ToResponse().Role = %v, want %v", resp.Role, RoleAdmin)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("ToResponse().CreatedAt = %v, want %v", resp.CreatedAt, now)
	}
	if resp.CreatedBy != "admin@example.com" {
		t.Errorf("ToResponse().CreatedBy = %v, want admin@example.com", resp.CreatedBy)
	}
	if resp.Organization != org {
		t.Error("ToResponse().Organization should reference the same organization")
	}
}

func TestUserOrganization_ToResponse_NilOrganization(t *testing.T) {
	uo := &UserOrganization{
		UserID:         1,
		OrganizationID: 1,
		Role:           RoleMember,
		CreatedBy:      "test@example.com",
		Organization:   nil,
	}

	resp := uo.ToResponse()

	if resp.Organization != nil {
		t.Error("ToResponse().Organization should be nil when UserOrganization.Organization is nil")
	}
	if resp.UserID != 1 {
		t.Errorf("ToResponse().UserID = %d, want 1", resp.UserID)
	}
}

func TestUserOrganization_ToResponse_ZeroValues(t *testing.T) {
	uo := &UserOrganization{}

	resp := uo.ToResponse()

	if resp.UserID != 0 {
		t.Errorf("ToResponse().UserID = %d, want 0", resp.UserID)
	}
	if resp.OrganizationID != 0 {
		t.Errorf("ToResponse().OrganizationID = %d, want 0", resp.OrganizationID)
	}
	if resp.Role != "" {
		t.Errorf("ToResponse().Role = %v, want empty", resp.Role)
	}
	if resp.CreatedBy != "" {
		t.Errorf("ToResponse().CreatedBy = %v, want empty", resp.CreatedBy)
	}
}

func TestUserOrganization_TableName(t *testing.T) {
	uo := UserOrganization{}
	if uo.TableName() != "user_organizations" {
		t.Errorf("TableName() = %v, want user_organizations", uo.TableName())
	}
}

func TestUserOrganizationRoleUpdateRequest(t *testing.T) {
	req := UserOrganizationRoleUpdateRequest{
		Role: RoleAdmin,
	}

	if req.Role != RoleAdmin {
		t.Errorf("Role = %v, want %v", req.Role, RoleAdmin)
	}
}

func TestUserMembership(t *testing.T) {
	org := &Organization{ID: 1, Name: "Test Org"}

	membership := UserMembership{
		UserID:         1,
		OrganizationID: 1,
		Role:           RoleAdmin,
		Organization:   org,
	}

	if membership.UserID != 1 {
		t.Errorf("UserID = %d, want 1", membership.UserID)
	}
	if membership.Role != RoleAdmin {
		t.Errorf("Role = %v, want %v", membership.Role, RoleAdmin)
	}
	if membership.Organization != org {
		t.Error("Organization should reference the same org")
	}
}

func TestUserMembershipsResponse(t *testing.T) {
	memberships := []UserMembership{
		{UserID: 1, OrganizationID: 1, Role: RoleAdmin},
		{UserID: 1, OrganizationID: 2, Role: RoleMember},
	}

	resp := UserMembershipsResponse{Memberships: memberships}

	if len(resp.Memberships) != 2 {
		t.Errorf("len(Memberships) = %d, want 2", len(resp.Memberships))
	}
}

func TestUserMembershipsResponse_Empty(t *testing.T) {
	resp := UserMembershipsResponse{Memberships: []UserMembership{}}

	if len(resp.Memberships) != 0 {
		t.Errorf("len(Memberships) = %d, want 0", len(resp.Memberships))
	}
}
