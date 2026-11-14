package mapper

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/schema"
)

func SolucaoToResponse(entity *schema.Solucao) *solucao.SolucaoResponse {
	return &solucao.SolucaoResponse{
		ResourceMetadata: common.ResourceMetadata{ID: entity.ID, CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt},
		StartupID:        entity.StartupID,
		Nome:             entity.Nome,
		Descricao:        entity.Descricao,
		Categoria:        entity.Categoria,
		PalavrasChave:    append([]string{}, entity.PalavrasChave...),
	}
}

func SolucoesToResponse(entities []*schema.Solucao) []*solucao.SolucaoResponse {
	response := make([]*solucao.SolucaoResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, SolucaoToResponse(entity))
	}
	return response
}
