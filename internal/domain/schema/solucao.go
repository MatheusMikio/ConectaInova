package schema

import "gorm.io/gorm"

type Solucao struct {
	gorm.Model
	StartupID     uint
	Nome          string
	Descricao     string
	Categoria     enum.Categoria
	PalavrasChave []string
}
/