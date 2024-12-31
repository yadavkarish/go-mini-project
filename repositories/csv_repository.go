package repository

import (
	"context"
	"csv-microservice/models"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Repository interface to define the data access contract.
type RepositoryInterface interface {
	InsertRecord(ctx context.Context, record interface{}) error
	DeleteRecord(ctx context.Context, id int) error
	QueryRecords(ctx context.Context, queryParams map[string]interface{}, offset, limit int) ([]models.User, error)
	AddRecord(record models.User) error
	BulkInsert(records []models.User) error
}

// Repository implementation
type Repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{Db: db}
}

func (r *Repository) InsertRecord(ctx context.Context, record interface{}) error {
	return r.Db.WithContext(ctx).Create(record).Error
}

func (r *Repository) DeleteRecord(ctx context.Context, id int) error {
	// return r.Db.WithContext(ctx).Delete(&CSV{}, id).Error
	return errors.ErrUnsupported
}

func (r *Repository) QueryRecords(ctx context.Context, queryParams map[string]interface{}, offset, limit int) ([]models.User, error) {
	var results []models.User
	query := r.Db.WithContext(ctx)

	// Apply filters dynamically
	for key, value := range queryParams {
		if key == "name" {
			query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(value.(string))+"%")
		} else {
			query = query.Where(key+" = ?", value)
		}
	}

	// Apply pagination
	err := query.Offset(offset).Limit(limit).Find(&results).Error
	return results, err
}

// unused method
func (r *Repository) AddRecord(record models.User) error {
	if r.Db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if err := r.Db.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to insert record into database: %w", err)
	}
	return nil
}

// BulkInsert inserts multiple records in a single transaction.
func (r *Repository) BulkInsert(records []models.User) error {
	if len(records) == 0 {
		return nil
	}

	tx := r.Db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(&records).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
