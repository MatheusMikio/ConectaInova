package mapper

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/schema"
)

func NecessidadeToResponse(entity *schema.Necessidade) *necessidade.NecessidadeResponse {
	return &necessidade.NecessidadeResponse{
		ResourceMetadata:  common.ResourceMetadata{ID: entity.ID, CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt},
		LojistaID:         entity.LojistaID,
		Titulo:            entity.Titulo,
		Descricao:         entity.Descricao,
		Categoria:         entity.Categoria,
		Prioridade:        entity.Prioridade,
		Status:            entity.Status,
		Impacto:           entity.Impacto,
		ResultadoEsperado: entity.ResultadoEsperado,
		Restricoes:        entity.Restricoes,
		CanalVenda:        entity.CanalVenda,
		CriteriosSucesso:  append([]string{}, entity.CriteriosSucesso...),
	}
}

func NecessidadesToResponse(entities []*schema.Necessidade) []*necessidade.NecessidadeResponse {
	response := make([]*necessidade.NecessidadeResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, NecessidadeToResponse(entity))
	}
	return response
}

func NecessidadeToListResponse(entity *schema.Necessidade) *necessidade.NecessidadeListResponse {
	return &necessidade.NecessidadeListResponse{
		NecessidadeResponse: *NecessidadeToResponse(entity),
		LojistaNome:         entity.Lojista.Nome,
	}
}

func NecessidadesToListResponse(entities []*schema.Necessidade) []*necessidade.NecessidadeListResponse {
	response := make([]*necessidade.NecessidadeListResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, NecessidadeToListResponse(entity))
	}
	return response
}
