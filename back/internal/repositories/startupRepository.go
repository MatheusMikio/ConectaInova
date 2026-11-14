package repositories

import (
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

type IStartupRepository interface {
	IBaseRepository[schema.Startup]
}

type StartupRepository struct {
	BaseRepository[schema.Startup]
}

func NewStartupRepository(db *gorm.DB) IStartupRepository {
	return &StartupRepository{
		BaseRepository: BaseRepository[schema.Startup]{DB: db},
	}
}
