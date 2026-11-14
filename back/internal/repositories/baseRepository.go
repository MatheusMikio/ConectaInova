package repositories

import (
	query "conectainova/internal/domain/models/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IBaseRepository[T any] interface {
	GetAll(page, size uint) ([]*T, error)
	GetByID(id uint) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(entity *T) error
}

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func (br *BaseRepository[T]) GetAll(page, size uint) ([]*T, error) {
	pagination := query.NewPaginationQuery(page, size)
	entities := make([]*T, 0)
	// Check before multiplying: a uint64 overflow can wrap a huge page back to 0.
	if uint64(pagination.Page-1) > (1<<31-1)/uint64(pagination.Size) {
		return entities, nil
	}
	offset := uint64(pagination.Page-1) * uint64(pagination.Size)

	err := br.DB.Order("id DESC").Limit(int(pagination.Size)).Offset(int(offset)).Find(&entities).Error
	return entities, err
}

func (br *BaseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	if err := br.DB.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (br *BaseRepository[T]) Create(entity *T) error {
	return br.DB.Omit(clause.Associations).Create(entity).Error
}

func (br *BaseRepository[T]) Update(entity *T) error {
	return br.DB.Omit(clause.Associations).Save(entity).Error
}

func (br *BaseRepository[T]) Delete(entity *T) error {
	return br.DB.Delete(entity).Error
}
