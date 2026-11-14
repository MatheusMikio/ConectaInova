package repositories

import (
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

type ILojistaRepository interface {
	IBaseRepository[schema.Lojista]
}

type LojistaRepository struct {
	BaseRepository[schema.Lojista]
}

func NewLojistaRepository(db *gorm.DB) ILojistaRepository {
	return &LojistaRepository{
		BaseRepository: BaseRepository[schema.Lojista]{DB: db},
	}
}
