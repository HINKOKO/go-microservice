package database

import (
	"context"
	"errors"

	"github.com/HINKOKO/go-microservice/internal/dberrors"
	"github.com/HINKOKO/go-microservice/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c Client) GetAllVendors(ctx context.Context) ([]models.Vendor, error) {
	var vendors []models.Vendor
	result := c.DB.WithContext(ctx).Find(&vendors)

	return vendors, result.Error
}

func (c Client) AddVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	vendor.VendorID = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&vendor)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}
	return vendor, nil
}

func (c Client) GetVendorById(ctx context.Context, ID string) (*models.Vendor, error) {
	vendor := &models.Vendor{}
	result := c.DB.WithContext(ctx).Where(&models.Vendor{VendorID: ID}).First(&vendor)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "vendor", ID: ID}
		}
		// If not the previous error, return it whatever it is
		return nil, result.Error
	}
	return vendor, nil
}

func (c Client) UpdateVendor(ctx context.Context, vendor *models.Vendor) (*models.Vendor, error) {
	var vendors []models.Vendor
	result := c.DB.WithContext(ctx).
		Model(&vendor).
		Clauses(clause.Returning{}).
		Where(&models.Vendor{VendorID: vendor.VendorID}).
		Updates(models.Vendor{
			Name:    vendor.Name,
			Contact: vendor.Contact,
			Phone:   vendor.Phone,
			Email:   vendor.Email,
			Address: vendor.Address,
		})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, result.Error
		}
	}
	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{
			Entity: "vendor",
			ID:     vendor.VendorID,
		}
	}
	return &vendors[0], nil
}

func (c Client) DeleteVendor(ctx context.Context, ID string) error {
	return c.DB.WithContext(ctx).Delete(&models.Vendor{VendorID: ID}).Error
}
