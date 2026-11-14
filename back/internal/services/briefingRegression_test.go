package services

import (
	"testing"

	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/schema"
)

func TestInvalidBriefingNeverReachesPersistence(t *testing.T) {
	blank := "   "
	criteria := []string{}
	for _, request := range []*necessidade.UpdateNecessidadeRequest{
		{ResultadoEsperado: &blank},
		{CriteriosSucesso: &criteria},
	} {
		repo := &necessidadeRepoStub{entity: schema.Necessidade{
			ResultadoEsperado: "Preservar o objetivo", CriteriosSucesso: []string{"Preservar o critério"},
		}}
		service := NewNecessidadeService(repo, nil)
		if _, err := service.UpdateNecessidade(1, request); err == nil {
			t.Fatal("expected an invalid briefing error")
		}
		if repo.writes != 0 || repo.entity.ResultadoEsperado != "Preservar o objetivo" || len(repo.entity.CriteriosSucesso) != 1 {
			t.Fatal("invalid request changed the entity or reached persistence")
		}
	}
	repo := &necessidadeRepoStub{}
	service := NewNecessidadeService(repo, nil)
	if _, err := service.CreateNecessidade(&necessidade.CreateNecessidadeRequest{
		ResultadoEsperado: blank, CriteriosSucesso: []string{"Critério"},
	}); err == nil || repo.writes != 0 {
		t.Fatal("invalid create reached persistence")
	}
}
