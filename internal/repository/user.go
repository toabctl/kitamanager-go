package repository

import (
	"github.com/eenemeene/kitamanager-go/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Organizations").Preload("Groups").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) AddToGroup(userID, groupID uint) error {
	user := &models.User{ID: userID}
	group := &models.Group{ID: groupID}
	return r.db.Model(user).Association("Groups").Append(group)
}

func (r *UserRepository) RemoveFromGroup(userID, groupID uint) error {
	user := &models.User{ID: userID}
	group := &models.Group{ID: groupID}
	return r.db.Model(user).Association("Groups").Delete(group)
}

func (r *UserRepository) AddToOrganization(userID, orgID uint) error {
	user := &models.User{ID: userID}
	org := &models.Organization{ID: orgID}
	return r.db.Model(user).Association("Organizations").Append(org)
}

func (r *UserRepository) RemoveFromOrganization(userID, orgID uint) error {
	user := &models.User{ID: userID}
	org := &models.Organization{ID: orgID}
	return r.db.Model(user).Association("Organizations").Delete(org)
}
