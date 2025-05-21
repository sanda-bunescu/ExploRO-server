package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BaseRepositoryInterface[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetByID(ctx context.Context, id any) (*T, error)
	GetAll(ctx context.Context) ([]T, error)
	Update(ctx context.Context, entity *T) error
	SoftDelete(ctx context.Context, entity *T) error
	AddRange(ctx context.Context, entities []*T) error
	SoftDeleteRange(ctx context.Context, entities []*T) error
}

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		DB: db,
	}
}

var _ BaseRepositoryInterface[any] = (*BaseRepository[any])(nil)

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(&entity).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id any) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at is NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) GetAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.DB.WithContext(ctx).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) SoftDelete(ctx context.Context, entity *T) error {
	var err = r.DB.WithContext(ctx).FirstOrInit(&entity).UpdateColumn("deleted_at", time.Now()).Error
	if err != nil {
		return fmt.Errorf("failed to soft delete entity: %w", err)
	}

	return nil
}

func (r *BaseRepository[T]) AddRange(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return nil
	}

	return r.DB.WithContext(ctx).Create(entities).Error
}

func (r *BaseRepository[T]) SoftDeleteRange(ctx context.Context, entities []*T) error {
	for _, entity := range entities {
		err := r.SoftDelete(ctx, entity)
		if err != nil {
			return fmt.Errorf("failed to soft delete entities: %w", err)
		}
	}

	return nil
}
