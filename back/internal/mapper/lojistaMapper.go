package mapper

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/dto/lojista"
	"conectainova/internal/domain/schema"
)

func LojistaToResponse(entity *schema.Lojista) *lojista.LojistaResponse {
	return &lojista.LojistaResponse{
		ResourceMetadata: common.ResourceMetadata{ID: entity.ID, CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt},
		Nome:             entity.Nome,
		Segmento:         entity.Segmento,
	}
}

func LojistasToResponse(entities []*schema.Lojista) []*lojista.LojistaResponse {
	response := make([]*lojista.LojistaResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, LojistaToResponse(entity))
	}
	return response
}
