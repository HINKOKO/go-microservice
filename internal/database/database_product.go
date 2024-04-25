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

func (c Client) GetAllProducts(ctx context.Context, vendorID string) ([]models.Product, error) {
	var products []models.Product
	result := c.DB.WithContext(ctx).Where(models.Product{VendorID: vendorID}).Find(&products)

	return products, result.Error
}

func (c Client) AddProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product.ProductID = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}
	return product, nil
}

func (c Client) GetProductById(ctx context.Context, ID string) (*models.Product, error) {
	product := &models.Product{}
	result := c.DB.WithContext(ctx).Where(&models.Product{ProductID: ID}).First(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "product", ID: ID}
		}
		// If not the previous error, return it whatever it is
		return nil, result.Error
	}
	return product, nil
}

func (c Client) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	var products []models.Product
	result := c.DB.WithContext(ctx).
		Model(&product).
		Clauses(clause.Returning{}).
		Where(&models.Product{ProductID: product.ProductID}).
		Updates(models.Product{
			Name:     product.Name,
			Price:    product.Price,
			VendorID: product.VendorID,
		})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, result.Error
		}
	}
	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{
			Entity: "product",
			ID:     product.ProductID,
		}
	}
	return &products[0], nil
}

func (c Client) DeleteProduct(ctx context.Context, ID string) error {
	return c.DB.WithContext(ctx).Delete(&models.Product{ProductID: ID}).Error
}
