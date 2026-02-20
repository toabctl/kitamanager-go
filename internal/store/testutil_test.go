package store

import (
	"testing"

	"gorm.io/gorm"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/testutil"
)

// setupTestDB creates a PostgreSQL test database using testcontainers.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testutil.SetupTestDB(t)
}

// createTestOrganization creates an organization for testing.
func createTestOrganization(t *testing.T, db *gorm.DB, name string) *models.Organization {
	t.Helper()
	return testutil.CreateTestOrganization(t, db, name)
}

// createTestUser creates a user for testing.
func createTestUser(t *testing.T, db *gorm.DB, name, email string) *models.User {
	t.Helper()
	return testutil.CreateTestUser(t, db, name, email, "hashedpassword")
}

// createTestUserOrganization creates a user-organization membership for testing.
func createTestUserOrganization(t *testing.T, db *gorm.DB, userID, orgID uint, role models.Role) *models.UserOrganization {
	t.Helper()
	return testutil.CreateTestUserOrganization(t, db, userID, orgID, role)
}

// createTestPayPlan creates a pay plan for the given organization.
func createTestPayPlan(t *testing.T, db *gorm.DB, orgID uint) *models.PayPlan {
	t.Helper()
	return testutil.CreateTestPayPlan(t, db, "Test Pay Plan", orgID)
}

// getDefaultSectionID returns the ID of the default section for an organization.
func getDefaultSectionID(t *testing.T, db *gorm.DB, orgID uint) uint {
	t.Helper()
	var section models.Section
	if err := db.Where("organization_id = ? AND is_default = ?", orgID, true).First(&section).Error; err != nil {
		t.Fatalf("failed to get default section for org %d: %v", orgID, err)
	}
	return section.ID
}
