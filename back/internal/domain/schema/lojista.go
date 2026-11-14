package schema

import "gorm.io/gorm"

type Lojista struct {
	gorm.Model
	Nome     string `gorm:"not null" json:"nome"`
	Segmento string `gorm:"not null" json:"segmento"`
}

func (Lojista) TableName() string { return "lojistas" }
