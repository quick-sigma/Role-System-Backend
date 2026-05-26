package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type SQLiteRepository[T any] struct {
	db *gorm.DB
}

func NewSQLiteRepository[T any](db *gorm.DB) *SQLiteRepository[T] {
	return &SQLiteRepository[T]{db: db}
}

func (r *SQLiteRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *SQLiteRepository[T]) GetByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	result := r.db.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("entity not found with id %d", id)
		}
		return nil, result.Error
	}
	return &entity, nil
}

func (r *SQLiteRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *SQLiteRepository[T]) Delete(ctx context.Context, id uint) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *SQLiteRepository[T]) List(ctx context.Context) ([]T, error) {
	var entities []T
	result := r.db.WithContext(ctx).Find(&entities)
	return entities, result.Error
}
