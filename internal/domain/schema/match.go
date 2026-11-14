package schema

import "gorm.io/gorm"

type Match struct {
	gorm.Model
	NecessidadeID uint
	SolucaoID     uint
	Pontuacao     float64
}
