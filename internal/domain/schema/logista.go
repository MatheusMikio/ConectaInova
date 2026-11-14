package schema

import "gorm.io/gorm"

type Logista struct {
	gorm.Model
	Nome     string
	Segmento string
}
