package necessidade

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/enum"
)

type CreateNecessidadeRequest struct {
	LojistaID         uint                    `json:"lojista_id"`
	Titulo            string                  `json:"titulo"`
	Descricao         string                  `json:"descricao"`
	Categoria         enum.Categoria          `json:"categoria"`
	Prioridade        enum.Prioridade         `json:"prioridade"`
	Status            *enum.StatusNecessidade `json:"status,omitempty"`
	Impacto           string                  `json:"impacto"`
	ResultadoEsperado string                  `json:"resultado_esperado"`
	Restricoes        string                  `json:"restricoes"`
	CanalVenda        string                  `json:"canal_venda"`
	CriteriosSucesso  []string                `json:"criterios_sucesso"`
}

type UpdateNecessidadeRequest struct {
	LojistaID         *uint                   `json:"lojista_id"`
	Titulo            *string                 `json:"titulo"`
	Descricao         *string                 `json:"descricao"`
	Categoria         *enum.Categoria         `json:"categoria"`
	Prioridade        *enum.Prioridade        `json:"prioridade"`
	Status            *enum.StatusNecessidade `json:"status"`
	Impacto           *string                 `json:"impacto"`
	ResultadoEsperado *string                 `json:"resultado_esperado"`
	Restricoes        *string                 `json:"restricoes"`
	CanalVenda        *string                 `json:"canal_venda"`
	CriteriosSucesso  *[]string               `json:"criterios_sucesso"`
}

type NecessidadeResponse struct {
	common.ResourceMetadata
	LojistaID         uint                   `json:"lojista_id"`
	Titulo            string                 `json:"titulo"`
	Descricao         string                 `json:"descricao"`
	Categoria         enum.Categoria         `json:"categoria"`
	Prioridade        enum.Prioridade        `json:"prioridade"`
	Status            enum.StatusNecessidade `json:"status"`
	Impacto           string                 `json:"impacto"`
	ResultadoEsperado string                 `json:"resultado_esperado"`
	Restricoes        string                 `json:"restricoes"`
	CanalVenda        string                 `json:"canal_venda"`
	CriteriosSucesso  []string               `json:"criterios_sucesso"`
}

type NecessidadeListResponse struct {
	NecessidadeResponse
	LojistaNome string `json:"lojista_nome"`
}
