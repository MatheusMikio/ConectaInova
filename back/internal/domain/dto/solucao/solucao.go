package solucao

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/enum"
)

type CreateSolucaoRequest struct {
	StartupID     uint           `json:"startup_id"`
	Nome          string         `json:"nome"`
	Descricao     string         `json:"descricao"`
	Categoria     enum.Categoria `json:"categoria"`
	PalavrasChave []string       `json:"palavras_chave"`
}

type UpdateSolucaoRequest struct {
	StartupID     *uint           `json:"startup_id"`
	Nome          *string         `json:"nome"`
	Descricao     *string         `json:"descricao"`
	Categoria     *enum.Categoria `json:"categoria"`
	PalavrasChave *[]string       `json:"palavras_chave"`
}

type SolucaoResponse struct {
	common.ResourceMetadata
	StartupID     uint           `json:"startup_id"`
	Nome          string         `json:"nome"`
	Descricao     string         `json:"descricao"`
	Categoria     enum.Categoria `json:"categoria"`
	PalavrasChave []string       `json:"palavras_chave"`
}
