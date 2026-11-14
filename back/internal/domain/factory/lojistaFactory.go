package factory

import (
	"strings"

	"conectainova/internal/domain/dto/lojista"
	"conectainova/internal/domain/schema"
)

func NewLojista(request *lojista.CreateLojistaRequest) *schema.Lojista {
	return &schema.Lojista{
		Nome:     strings.TrimSpace(request.Nome),
		Segmento: strings.TrimSpace(request.Segmento),
	}
}

func UpdateLojista(entity *schema.Lojista, request *lojista.UpdateLojistaRequest) *schema.Lojista {
	updated := *entity
	if request.Nome != nil {
		updated.Nome = strings.TrimSpace(*request.Nome)
	}
	if request.Segmento != nil {
		updated.Segmento = strings.TrimSpace(*request.Segmento)
	}
	return &updated
}
