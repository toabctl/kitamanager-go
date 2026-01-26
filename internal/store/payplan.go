package store

import (
	"gorm.io/gorm"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

type PayplanStore struct {
	db *gorm.DB
}

func NewPayplanStore(db *gorm.DB) *PayplanStore {
	return &PayplanStore{db: db}
}

// Payplan CRUD

func (s *PayplanStore) FindAll(limit, offset int) ([]models.Payplan, int64, error) {
	var payplans []models.Payplan
	var total int64

	if err := s.db.Model(&models.Payplan{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Limit(limit).Offset(offset).Find(&payplans).Error; err != nil {
		return nil, 0, err
	}

	return payplans, total, nil
}

func (s *PayplanStore) FindByID(id uint) (*models.Payplan, error) {
	var payplan models.Payplan
	if err := s.db.First(&payplan, id).Error; err != nil {
		return nil, err
	}
	return &payplan, nil
}

func (s *PayplanStore) FindByIDWithDetails(id uint) (*models.Payplan, error) {
	var payplan models.Payplan
	if err := s.db.
		Preload("Periods", func(db *gorm.DB) *gorm.DB {
			return db.Order("from_date ASC")
		}).
		Preload("Periods.Entries", func(db *gorm.DB) *gorm.DB {
			return db.Order("min_age ASC")
		}).
		Preload("Periods.Entries.Properties", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		First(&payplan, id).Error; err != nil {
		return nil, err
	}
	return &payplan, nil
}

func (s *PayplanStore) Create(payplan *models.Payplan) error {
	return s.db.Create(payplan).Error
}

func (s *PayplanStore) Update(payplan *models.Payplan) error {
	return s.db.Save(payplan).Error
}

func (s *PayplanStore) Delete(id uint) error {
	return s.db.Delete(&models.Payplan{}, id).Error
}

// Period CRUD

func (s *PayplanStore) FindPeriodByID(id uint) (*models.PayplanPeriod, error) {
	var period models.PayplanPeriod
	if err := s.db.
		Preload("Entries", func(db *gorm.DB) *gorm.DB {
			return db.Order("min_age ASC")
		}).
		Preload("Entries.Properties", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		First(&period, id).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (s *PayplanStore) FindPeriodsByPayplanID(payplanID uint) ([]models.PayplanPeriod, error) {
	var periods []models.PayplanPeriod
	if err := s.db.Where("payplan_id = ?", payplanID).Order("from_date ASC").Find(&periods).Error; err != nil {
		return nil, err
	}
	return periods, nil
}

func (s *PayplanStore) CreatePeriod(period *models.PayplanPeriod) error {
	return s.db.Create(period).Error
}

func (s *PayplanStore) UpdatePeriod(period *models.PayplanPeriod) error {
	return s.db.Save(period).Error
}

func (s *PayplanStore) DeletePeriod(id uint) error {
	return s.db.Delete(&models.PayplanPeriod{}, id).Error
}

// Entry CRUD

func (s *PayplanStore) FindEntryByID(id uint) (*models.PayplanEntry, error) {
	var entry models.PayplanEntry
	if err := s.db.
		Preload("Properties", func(db *gorm.DB) *gorm.DB {
			return db.Order("name ASC")
		}).
		First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (s *PayplanStore) CreateEntry(entry *models.PayplanEntry) error {
	return s.db.Create(entry).Error
}

func (s *PayplanStore) UpdateEntry(entry *models.PayplanEntry) error {
	return s.db.Save(entry).Error
}

func (s *PayplanStore) DeleteEntry(id uint) error {
	return s.db.Delete(&models.PayplanEntry{}, id).Error
}

// Property CRUD

func (s *PayplanStore) FindPropertyByID(id uint) (*models.PayplanProperty, error) {
	var property models.PayplanProperty
	if err := s.db.First(&property, id).Error; err != nil {
		return nil, err
	}
	return &property, nil
}

func (s *PayplanStore) CreateProperty(property *models.PayplanProperty) error {
	return s.db.Create(property).Error
}

func (s *PayplanStore) UpdateProperty(property *models.PayplanProperty) error {
	return s.db.Save(property).Error
}

func (s *PayplanStore) DeleteProperty(id uint) error {
	return s.db.Delete(&models.PayplanProperty{}, id).Error
}

// Organization payplan assignment

func (s *PayplanStore) AssignPayplanToOrg(orgID, payplanID uint) error {
	return s.db.Model(&models.Organization{}).Where("id = ?", orgID).Update("payplan_id", payplanID).Error
}

func (s *PayplanStore) RemovePayplanFromOrg(orgID uint) error {
	return s.db.Model(&models.Organization{}).Where("id = ?", orgID).Update("payplan_id", nil).Error
}
