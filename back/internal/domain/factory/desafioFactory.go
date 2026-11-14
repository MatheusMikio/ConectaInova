package factory

import (
	"strings"

	"conectainova/internal/domain/dto/desafio"
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/schema"
)

func NewDesafio(request *desafio.CreateDesafioRequest) *schema.Desafio {
	entity := &schema.Desafio{
		NecessidadeID:    request.NecessidadeID,
		Titulo:           strings.TrimSpace(request.Titulo),
		Descricao:        strings.TrimSpace(request.Descricao),
		Status:           enum.DesafioAberto,
		CriteriosSucesso: normalizeLines(request.CriteriosSucesso),
		PrazoDias:        7,
	}
	if request.Status != nil {
		entity.Status = *request.Status
	}
	if request.PrazoDias != nil {
		entity.PrazoDias = *request.PrazoDias
	}
	return entity
}

func UpdateDesafio(entity *schema.Desafio, request *desafio.UpdateDesafioRequest) *schema.Desafio {
	updated := *entity
	updated.CriteriosSucesso = append([]string{}, entity.CriteriosSucesso...)
	if request.CriteriosSucesso != nil {
		updated.CriteriosSucesso = normalizeLines(*request.CriteriosSucesso)
	}
	if request.PrazoDias != nil {
		updated.PrazoDias = *request.PrazoDias
	}
	if request.NecessidadeID != nil {
		updated.NecessidadeID = *request.NecessidadeID
	}
	if request.Titulo != nil {
		updated.Titulo = strings.TrimSpace(*request.Titulo)
	}
	if request.Descricao != nil {
		updated.Descricao = strings.TrimSpace(*request.Descricao)
	}
	if request.Status != nil {
		updated.Status = *request.Status
	}
	return &updated
}
