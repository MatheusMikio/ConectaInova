package desafio

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/enum"
)

type CreateDesafioRequest struct {
	NecessidadeID    uint                `json:"necessidade_id"`
	Titulo           string              `json:"titulo"`
	Descricao        string              `json:"descricao"`
	Status           *enum.StatusDesafio `json:"status,omitempty"`
	CriteriosSucesso []string            `json:"criterios_sucesso"`
	PrazoDias        *int                `json:"prazo_dias,omitempty"`
}

type UpdateDesafioRequest struct {
	NecessidadeID    *uint               `json:"necessidade_id"`
	Titulo           *string             `json:"titulo"`
	Descricao        *string             `json:"descricao"`
	Status           *enum.StatusDesafio `json:"status"`
	CriteriosSucesso *[]string           `json:"criterios_sucesso"`
	PrazoDias        *int                `json:"prazo_dias"`
}

type DesafioResponse struct {
	common.ResourceMetadata
	NecessidadeID    uint               `json:"necessidade_id"`
	Titulo           string             `json:"titulo"`
	Descricao        string             `json:"descricao"`
	Status           enum.StatusDesafio `json:"status"`
	CriteriosSucesso []string           `json:"criterios_sucesso"`
	PrazoDias        int                `json:"prazo_dias"`
}
