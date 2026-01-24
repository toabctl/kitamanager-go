package repository

import (
	"github.com/eenemeene/kitamanager-go/internal/models"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) FindAll() ([]models.Organization, error) {
	var organizations []models.Organization
	if err := r.db.Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *OrganizationRepository) FindByID(id uint) (*models.Organization, error) {
	var organization models.Organization
	if err := r.db.Preload("Users").Preload("Groups").First(&organization, id).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (r *OrganizationRepository) Create(organization *models.Organization) error {
	return r.db.Create(organization).Error
}

func (r *OrganizationRepository) Update(organization *models.Organization) error {
	return r.db.Save(organization).Error
}

func (r *OrganizationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Organization{}, id).Error
}
