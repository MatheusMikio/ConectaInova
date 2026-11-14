package schema

import "gorm.io/gorm"

type Match struct {
	gorm.Model
	NecessidadeID uint        `gorm:"not null;uniqueIndex:idx_match_pair" json:"necessidade_id"`
	Necessidade   Necessidade `gorm:"constraint:OnDelete:RESTRICT" json:"-"`
	SolucaoID     uint        `gorm:"not null;uniqueIndex:idx_match_pair" json:"solucao_id"`
	Solucao       Solucao     `gorm:"constraint:OnDelete:RESTRICT" json:"-"`
	Pontuacao     float64     `gorm:"not null" json:"pontuacao"`
}
