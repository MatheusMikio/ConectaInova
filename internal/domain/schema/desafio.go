package schema

import "gorm.io/gorm"

type Desafio struct {
	gorm.Model
	NecessidadeID uint
	Titulo        string
	Descricao     string
	Status        enum.StatusDesafio
}
