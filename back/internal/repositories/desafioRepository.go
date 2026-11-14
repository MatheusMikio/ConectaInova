package repositories

import (
	query "conectainova/internal/domain/models/query"
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

type IDesafioRepository interface {
	IBaseRepository[schema.Desafio]
	GetByNecessidadeID(necessidadeID, page, size uint) ([]*schema.Desafio, error)
}

func (r *DesafioRepository) GetByNecessidadeID(necessidadeID, page, size uint) ([]*schema.Desafio, error) {
	pagination := query.NewPaginationQuery(page, size)
	entities := make([]*schema.Desafio, 0)
	if uint64(pagination.Page-1) > (1<<31-1)/uint64(pagination.Size) {
		return entities, nil
	}
	offset := int(uint64(pagination.Page-1) * uint64(pagination.Size))
	err := r.DB.Where("necessidade_id = ?", necessidadeID).Order("id DESC").
		Limit(int(pagination.Size)).Offset(offset).Find(&entities).Error
	return entities, err
}

type DesafioRepository struct {
	BaseRepository[schema.Desafio]
}

func NewDesafioRepository(db *gorm.DB) IDesafioRepository {
	return &DesafioRepository{
		BaseRepository: BaseRepository[schema.Desafio]{DB: db},
	}
}
