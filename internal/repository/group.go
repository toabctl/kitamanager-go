package repository

import (
	"github.com/eenemeene/kitamanager-go/internal/models"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) FindAll() ([]models.Group, error) {
	var groups []models.Group
	if err := r.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GroupRepository) FindByID(id uint) (*models.Group, error) {
	var group models.Group
	if err := r.db.Preload("Users").Preload("Organizations").First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) Create(group *models.Group) error {
	return r.db.Create(group).Error
}

func (r *GroupRepository) Update(group *models.Group) error {
	return r.db.Save(group).Error
}

func (r *GroupRepository) Delete(id uint) error {
	return r.db.Delete(&models.Group{}, id).Error
}

func (r *GroupRepository) AddToOrganization(groupID, orgID uint) error {
	group := &models.Group{ID: groupID}
	org := &models.Organization{ID: orgID}
	return r.db.Model(group).Association("Organizations").Append(org)
}

func (r *GroupRepository) RemoveFromOrganization(groupID, orgID uint) error {
	group := &models.Group{ID: groupID}
	org := &models.Organization{ID: orgID}
	return r.db.Model(group).Association("Organizations").Delete(org)
}
