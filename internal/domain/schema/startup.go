package schema

import "gorm.io/gorm"

type Startup struct {
	gorm.Model
	Nome      string
	Descricao string
}
