package schema

import (
	"conectainova/internal/domain/enum"
	"gorm.io/gorm"
)

type Desafio struct {
	gorm.Model
	NecessidadeID    uint               `gorm:"not null;index" json:"necessidade_id"`
	Necessidade      Necessidade        `gorm:"constraint:OnDelete:RESTRICT" json:"-"`
	Titulo           string             `gorm:"not null" json:"titulo"`
	Descricao        string             `gorm:"not null" json:"descricao"`
	Status           enum.StatusDesafio `gorm:"not null" json:"status"`
	CriteriosSucesso []string           `gorm:"serializer:json;type:jsonb" json:"criterios_sucesso"`
	PrazoDias        int                `gorm:"not null" json:"prazo_dias"`
}
