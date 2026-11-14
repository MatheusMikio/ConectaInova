package repositories

import (
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

type ISolucaoRepository interface {
	IBaseRepository[schema.Solucao]
	GetByCategoria(categoria enum.Categoria) ([]*schema.Solucao, error)
}

type SolucaoRepository struct {
	BaseRepository[schema.Solucao]
}

func NewSolucaoRepository(db *gorm.DB) ISolucaoRepository {
	return &SolucaoRepository{
		BaseRepository: BaseRepository[schema.Solucao]{DB: db},
	}
}

func (s *SolucaoRepository) GetByCategoria(categoria enum.Categoria) ([]*schema.Solucao, error) {
	entities := make([]*schema.Solucao, 0)
	err := s.DB.Preload("Startup").Where("categoria = ?", categoria).Order("id").Find(&entities).Error
	return entities, err
}
