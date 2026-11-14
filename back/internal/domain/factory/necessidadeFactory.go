package factory

import (
	"strings"

	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/schema"
)

func NewNecessidade(request *necessidade.CreateNecessidadeRequest) *schema.Necessidade {
	entity := &schema.Necessidade{
		LojistaID:         request.LojistaID,
		Titulo:            strings.TrimSpace(request.Titulo),
		Descricao:         strings.TrimSpace(request.Descricao),
		Categoria:         request.Categoria,
		Prioridade:        request.Prioridade,
		Status:            enum.NecessidadeAberta,
		Impacto:           strings.TrimSpace(request.Impacto),
		ResultadoEsperado: strings.TrimSpace(request.ResultadoEsperado),
		Restricoes:        strings.TrimSpace(request.Restricoes),
		CanalVenda:        request.CanalVenda,
		CriteriosSucesso:  normalizeLines(request.CriteriosSucesso),
	}
	if request.Status != nil {
		entity.Status = *request.Status
	}
	return entity
}

func UpdateNecessidade(entity *schema.Necessidade, request *necessidade.UpdateNecessidadeRequest) *schema.Necessidade {
	updated := *entity
	updated.CriteriosSucesso = append([]string{}, entity.CriteriosSucesso...)
	if request.Impacto != nil {
		updated.Impacto = strings.TrimSpace(*request.Impacto)
	}
	if request.ResultadoEsperado != nil {
		updated.ResultadoEsperado = strings.TrimSpace(*request.ResultadoEsperado)
	}
	if request.Restricoes != nil {
		updated.Restricoes = strings.TrimSpace(*request.Restricoes)
	}
	if request.CanalVenda != nil {
		updated.CanalVenda = *request.CanalVenda
	}
	if request.CriteriosSucesso != nil {
		updated.CriteriosSucesso = normalizeLines(*request.CriteriosSucesso)
	}
	if request.LojistaID != nil {
		updated.LojistaID = *request.LojistaID
	}
	if request.Titulo != nil {
		updated.Titulo = strings.TrimSpace(*request.Titulo)
	}
	if request.Descricao != nil {
		updated.Descricao = strings.TrimSpace(*request.Descricao)
	}
	if request.Categoria != nil {
		updated.Categoria = *request.Categoria
	}
	if request.Prioridade != nil {
		updated.Prioridade = *request.Prioridade
	}
	if request.Status != nil {
		updated.Status = *request.Status
	}
	return &updated
}
