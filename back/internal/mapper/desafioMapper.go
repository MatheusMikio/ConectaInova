package mapper

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/dto/desafio"
	"conectainova/internal/domain/schema"
)

func DesafioToResponse(entity *schema.Desafio) *desafio.DesafioResponse {
	return &desafio.DesafioResponse{
		ResourceMetadata: common.ResourceMetadata{ID: entity.ID, CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt},
		NecessidadeID:    entity.NecessidadeID,
		Titulo:           entity.Titulo,
		Descricao:        entity.Descricao,
		Status:           entity.Status,
		CriteriosSucesso: append([]string{}, entity.CriteriosSucesso...),
		PrazoDias:        entity.PrazoDias,
	}
}

func DesafiosToResponse(entities []*schema.Desafio) []*desafio.DesafioResponse {
	response := make([]*desafio.DesafioResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, DesafioToResponse(entity))
	}
	return response
}
