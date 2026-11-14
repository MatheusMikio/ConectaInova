package factory

import (
	"conectainova/internal/domain/dto/avaliacao"
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/schema"
	"strings"
)

func NewAvaliacao(request *avaliacao.CreateAvaliacaoRequest, associacao, regiao string) *schema.Avaliacao {
	return &schema.Avaliacao{NecessidadeID: request.NecessidadeID, SolucaoID: request.SolucaoID,
		Mensagem: strings.TrimSpace(request.Mensagem), Associacao: associacao, Regiao: regiao, Status: enum.AvaliacaoSolicitada}
}

func UpdateAvaliacao(entity *schema.Avaliacao, request *avaliacao.UpdateAvaliacaoRequest) *schema.Avaliacao {
	updated := *entity
	if request.Status != nil {
		updated.Status = *request.Status
	}
	if request.Resposta != nil {
		updated.Resposta = strings.TrimSpace(*request.Resposta)
	}
	return &updated
}
