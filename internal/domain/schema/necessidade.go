package schema

import (
	"gorm.io/gorm"
)

type Necessidade struct {
	gorm.Model
	LogistaID uint
	Titulo    string
	Descricao string
	Categoria enum.Categoria
	Proridade enum.Proridade
	Status    enum.Status
}
