package schema

import (
	"conectainova/internal/domain/enum"
	"gorm.io/gorm"
)

type Avaliacao struct {
	gorm.Model
	NecessidadeID uint        `gorm:"not null;uniqueIndex:idx_avaliacao_par"`
	Necessidade   Necessidade `gorm:"constraint:OnDelete:RESTRICT"`
	SolucaoID     uint        `gorm:"not null;uniqueIndex:idx_avaliacao_par"`
	Solucao       Solucao     `gorm:"constraint:OnDelete:RESTRICT"`
	Mensagem      string
	Resposta      string
	Associacao    string
	Regiao        string
	Status        enum.StatusAvaliacao `gorm:"not null"`
}

func (Avaliacao) TableName() string { return "avaliacoes" }
