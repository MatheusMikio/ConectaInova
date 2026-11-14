package schema

import (
	"conectainova/internal/domain/enum"
	"gorm.io/gorm"
)

type Necessidade struct {
	gorm.Model
	LojistaID         uint                   `gorm:"not null;index" json:"lojista_id"`
	Lojista           Lojista                `gorm:"constraint:OnDelete:RESTRICT" json:"-"`
	Titulo            string                 `gorm:"not null" json:"titulo"`
	Descricao         string                 `gorm:"not null" json:"descricao"`
	Categoria         enum.Categoria         `gorm:"not null" json:"categoria"`
	Prioridade        enum.Prioridade        `gorm:"not null" json:"prioridade"`
	Status            enum.StatusNecessidade `gorm:"not null" json:"status"`
	Impacto           string                 `gorm:"not null" json:"impacto"`
	ResultadoEsperado string                 `gorm:"not null" json:"resultado_esperado"`
	Restricoes        string                 `gorm:"not null" json:"restricoes"`
	CanalVenda        string                 `gorm:"not null" json:"canal_venda"`
	CriteriosSucesso  []string               `gorm:"serializer:json;type:jsonb" json:"criterios_sucesso"`
}

func (Necessidade) TableName() string { return "necessidades" }
