package repositories

import (
	query "conectainova/internal/domain/models/query"
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IAvaliacaoRepository interface {
	IBaseRepository[schema.Avaliacao]
	GetFiltered(startupID, page, size uint) ([]*schema.Avaliacao, error)
	CreateOrGet(entity *schema.Avaliacao) (*schema.Avaliacao, bool, error)
}

type AvaliacaoRepository struct {
	BaseRepository[schema.Avaliacao]
}

func NewAvaliacaoRepository(db *gorm.DB) IAvaliacaoRepository {
	return &AvaliacaoRepository{BaseRepository: BaseRepository[schema.Avaliacao]{DB: db}}
}

func (r *AvaliacaoRepository) preload() *gorm.DB {
	return r.DB.Preload("Necessidade.Lojista").Preload("Solucao.Startup")
}

func (r *AvaliacaoRepository) GetByID(id uint) (*schema.Avaliacao, error) {
	var entity schema.Avaliacao
	if err := r.preload().First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *AvaliacaoRepository) GetFiltered(startupID, page, size uint) ([]*schema.Avaliacao, error) {
	pagination := query.NewPaginationQuery(page, size)
	entities := []*schema.Avaliacao{}
	if uint64(pagination.Page-1) > (1<<31-1)/uint64(pagination.Size) {
		return entities, nil
	}
	db := r.preload()
	if startupID > 0 {
		db = db.Joins("JOIN solucoes ON solucoes.id = avaliacoes.solucao_id AND solucoes.deleted_at IS NULL").
			Where("solucoes.startup_id = ?", startupID)
	}
	offset := int(uint64(pagination.Page-1) * uint64(pagination.Size))
	err := db.Order("avaliacoes.id DESC").Limit(int(pagination.Size)).Offset(offset).Find(&entities).Error
	return entities, err
}

func (r *AvaliacaoRepository) CreateOrGet(entity *schema.Avaliacao) (*schema.Avaliacao, bool, error) {
	result := r.DB.Omit(clause.Associations).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "necessidade_id"}, {Name: "solucao_id"}}, DoNothing: true,
	}).Create(entity)
	if result.Error != nil {
		return nil, false, result.Error
	}
	var saved schema.Avaliacao
	err := r.preload().Where("necessidade_id = ? AND solucao_id = ?", entity.NecessidadeID, entity.SolucaoID).First(&saved).Error
	return &saved, result.RowsAffected == 1, err
}
