package repository

import (
	"context"
	"csv-microservice/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository interface to define the data access contract.
type RepositoryInterface interface {
	InsertRecord(ctx context.Context, record interface{}) error
	DeleteRecord(ctx context.Context, id int) error
	QueryRecords(ctx context.Context, queryParams map[string]interface{}) ([]interface{}, error)
	AddRecord(record models.CSV) error
	BulkInsert(records []models.CSV) error
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

func (r *Repository) QueryRecords(ctx context.Context, queryParams map[string]interface{}) ([]interface{}, error) {
	var results []interface{}
	err := r.Db.WithContext(ctx).Where(queryParams).Find(&results).Error
	return results, err
}

// AddRecord adds a new record to the database
//
//	func (r *Repository) AddRecord(record models.CSV) error {
//		if err := r.Db.Create(&record).Error; err != nil {
//			utils.LogError("Failed to insert record into database", err)
//			return err
//		}
//		return nil
//	}
func (r *Repository) AddRecord(record models.CSV) error {
	if r.Db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if err := r.Db.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to insert record into database: %w", err)
	}
	return nil
}

// BulkInsert inserts multiple records in a single transaction.
func (r *Repository) BulkInsert(records []models.CSV) error {
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
