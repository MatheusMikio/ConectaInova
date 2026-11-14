package schema

import "gorm.io/gorm"

type Startup struct {
	gorm.Model
	Nome      string `gorm:"not null" json:"nome"`
	Descricao string `gorm:"not null" json:"descricao"`
}
