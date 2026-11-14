package factory

import (
	"strings"

	"conectainova/internal/domain/dto/startup"
	"conectainova/internal/domain/schema"
)

func NewStartup(request *startup.CreateStartupRequest) *schema.Startup {
	return &schema.Startup{
		Nome:      strings.TrimSpace(request.Nome),
		Descricao: strings.TrimSpace(request.Descricao),
	}
}

func UpdateStartup(entity *schema.Startup, request *startup.UpdateStartupRequest) *schema.Startup {
	updated := *entity
	if request.Nome != nil {
		updated.Nome = strings.TrimSpace(*request.Nome)
	}
	if request.Descricao != nil {
		updated.Descricao = strings.TrimSpace(*request.Descricao)
	}
	return &updated
}
