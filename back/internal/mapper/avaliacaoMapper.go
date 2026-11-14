package mapper

import (
	"conectainova/internal/domain/dto/avaliacao"
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/schema"
)

func AvaliacaoToResponse(entity *schema.Avaliacao) *avaliacao.AvaliacaoResponse {
	return &avaliacao.AvaliacaoResponse{
		ResourceMetadata: common.ResourceMetadata{ID: entity.ID, CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt},
		Necessidade:      *NecessidadeToListResponse(&entity.Necessidade),
		Solucao:          *SolucaoToResponse(&entity.Solucao), Startup: *StartupToResponse(&entity.Solucao.Startup),
		Mensagem: entity.Mensagem, Resposta: entity.Resposta, Associacao: entity.Associacao, Regiao: entity.Regiao, Status: entity.Status,
	}
}

func AvaliacoesToResponse(entities []*schema.Avaliacao) []*avaliacao.AvaliacaoResponse {
	response := make([]*avaliacao.AvaliacaoResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, AvaliacaoToResponse(entity))
	}
	return response
}
