package factory

import (
	"strings"

	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/schema"
)

func NewSolucao(request *solucao.CreateSolucaoRequest) *schema.Solucao {
	return &schema.Solucao{
		StartupID:     request.StartupID,
		Nome:          strings.TrimSpace(request.Nome),
		Descricao:     strings.TrimSpace(request.Descricao),
		Categoria:     request.Categoria,
		PalavrasChave: normalizeKeywords(request.PalavrasChave),
	}
}

func UpdateSolucao(entity *schema.Solucao, request *solucao.UpdateSolucaoRequest) *schema.Solucao {
	updated := *entity
	if request.StartupID != nil {
		updated.StartupID = *request.StartupID
	}
	if request.Nome != nil {
		updated.Nome = strings.TrimSpace(*request.Nome)
	}
	if request.Descricao != nil {
		updated.Descricao = strings.TrimSpace(*request.Descricao)
	}
	if request.Categoria != nil {
		updated.Categoria = *request.Categoria
	}
	updated.PalavrasChave = append([]string{}, entity.PalavrasChave...)
	if request.PalavrasChave != nil {
		updated.PalavrasChave = normalizeKeywords(*request.PalavrasChave)
	}
	return &updated
}

func normalizeKeywords(keywords []string) []string {
	result := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		result = append(result, strings.TrimSpace(keyword))
	}
	return result
}
