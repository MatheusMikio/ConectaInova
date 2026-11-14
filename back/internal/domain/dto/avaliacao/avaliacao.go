package avaliacao

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/dto/startup"
	"conectainova/internal/domain/enum"
)

type CreateAvaliacaoRequest struct {
	NecessidadeID uint   `json:"necessidade_id"`
	SolucaoID     uint   `json:"solucao_id"`
	Mensagem      string `json:"mensagem"`
}

type UpdateAvaliacaoRequest struct {
	Status   *enum.StatusAvaliacao `json:"status"`
	Resposta *string               `json:"resposta"`
}

type AvaliacaoResponse struct {
	common.ResourceMetadata
	Necessidade necessidade.NecessidadeListResponse `json:"necessidade"`
	Solucao     solucao.SolucaoResponse             `json:"solucao"`
	Startup     startup.StartupResponse             `json:"startup"`
	Mensagem    string                              `json:"mensagem"`
	Resposta    string                              `json:"resposta"`
	Associacao  string                              `json:"associacao"`
	Regiao      string                              `json:"regiao"`
	Status      enum.StatusAvaliacao                `json:"status"`
}
