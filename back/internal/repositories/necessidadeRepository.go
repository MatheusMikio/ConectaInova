package repositories

import (
	query "conectainova/internal/domain/models/query"
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

type INecessidadeRepository interface {
	IBaseRepository[schema.Necessidade]
	GetAllWithLojista(page, size uint) ([]*schema.Necessidade, error)
}

func (r *NecessidadeRepository) GetAllWithLojista(page, size uint) ([]*schema.Necessidade, error) {
	pagination := query.NewPaginationQuery(page, size)
	entities := make([]*schema.Necessidade, 0)
	if uint64(pagination.Page-1) > (1<<31-1)/uint64(pagination.Size) {
		return entities, nil
	}
	offset := int(uint64(pagination.Page-1) * uint64(pagination.Size))
	err := r.DB.Preload("Lojista").Order("id DESC").
		Limit(int(pagination.Size)).Offset(offset).Find(&entities).Error
	return entities, err
}

type NecessidadeRepository struct {
	BaseRepository[schema.Necessidade]
}

func NewNecessidadeRepository(db *gorm.DB) INecessidadeRepository {
	return &NecessidadeRepository{
		BaseRepository: BaseRepository[schema.Necessidade]{DB: db},
	}
}
