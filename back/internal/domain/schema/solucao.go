package schema

import (
	"conectainova/internal/domain/enum"
	"gorm.io/gorm"
)

type Solucao struct {
	gorm.Model
	StartupID     uint           `gorm:"not null;index" json:"startup_id"`
	Startup       Startup        `gorm:"constraint:OnDelete:RESTRICT" json:"-"`
	Nome          string         `gorm:"not null" json:"nome"`
	Descricao     string         `gorm:"not null" json:"descricao"`
	Categoria     enum.Categoria `gorm:"not null" json:"categoria"`
	PalavrasChave []string       `gorm:"serializer:json;type:jsonb" json:"palavras_chave"`
}

func (Solucao) TableName() string { return "solucoes" }
