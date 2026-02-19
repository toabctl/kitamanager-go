package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

type GovernmentFundingBillPeriodStore struct {
	db *gorm.DB
}

func NewGovernmentFundingBillPeriodStore(db *gorm.DB) *GovernmentFundingBillPeriodStore {
	return &GovernmentFundingBillPeriodStore{db: db}
}

func (s *GovernmentFundingBillPeriodStore) Create(ctx context.Context, period *models.GovernmentFundingBillPeriod) error {
	return DBFromContext(ctx, s.db).Create(period).Error
}

func (s *GovernmentFundingBillPeriodStore) FindByID(ctx context.Context, id uint) (*models.GovernmentFundingBillPeriod, error) {
	var period models.GovernmentFundingBillPeriod
	if err := DBFromContext(ctx, s.db).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Preload("Children.Payments").
		First(&period, id).Error; err != nil {
		return nil, WrapNotFound(err)
	}
	return &period, nil
}

func (s *GovernmentFundingBillPeriodStore) FindByOrganization(ctx context.Context, orgID uint, limit, offset int) ([]models.GovernmentFundingBillPeriod, int64, error) {
	var periods []models.GovernmentFundingBillPeriod
	var total int64

	db := DBFromContext(ctx, s.db).Where("organization_id = ?", orgID)

	if err := db.Model(&models.GovernmentFundingBillPeriod{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("from_date DESC").Limit(limit).Offset(offset).Find(&periods).Error; err != nil {
		return nil, 0, err
	}

	return periods, total, nil
}

func (s *GovernmentFundingBillPeriodStore) Delete(ctx context.Context, id uint) error {
	return DBFromContext(ctx, s.db).Delete(&models.GovernmentFundingBillPeriod{}, id).Error
}
