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

func (c Client) GetAllServices(ctx context.Context) ([]models.Service, error) {
	var services []models.Service
	result := c.DB.WithContext(ctx).Find(&services)

	return services, result.Error
}

func (c Client) AddService(ctx context.Context, service *models.Service) (*models.Service, error) {
	service.ServiceID = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&service)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &dberrors.ConflictError{}
		}
		return nil, result.Error
	}
	return service, nil
}

func (c Client) GetServiceById(ctx context.Context, ID string) (*models.Service, error) {
	service := &models.Service{}
	result := c.DB.WithContext(ctx).Where(&models.Service{ServiceID: ID}).First(&service)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &dberrors.NotFoundError{Entity: "service", ID: ID}
		}
		// If not the previous error, return it whatever it is
		return nil, result.Error
	}
	return service, nil
}

func (c Client) UpdateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	var services []models.Service
	result := c.DB.WithContext(ctx).
		Model(&service).
		Clauses(clause.Returning{}).
		Where(&models.Service{ServiceID: service.ServiceID}).
		Updates(models.Service{
			Name:  service.Name,
			Price: service.Price,
		})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, result.Error
		}
	}
	if result.RowsAffected == 0 {
		return nil, &dberrors.NotFoundError{
			Entity: "service",
			ID:     service.ServiceID,
		}
	}
	return &services[0], nil
}

func (c Client) DeleteService(ctx context.Context, ID string) error {
	return c.DB.WithContext(ctx).Delete(&models.Service{ServiceID: ID}).Error
}
