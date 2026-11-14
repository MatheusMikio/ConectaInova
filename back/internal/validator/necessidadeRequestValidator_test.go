package validator

import (
	"testing"

	"conectainova/internal/domain/dto/necessidade"
)

func TestExplicitEmptyBriefingIsRejected(t *testing.T) {
	blank := " \t\n "
	empty := ""
	emptyCriteria := []string{}
	blankCriteria := []string{"  "}
	tests := []*necessidade.UpdateNecessidadeRequest{
		{ResultadoEsperado: &blank},
		{ResultadoEsperado: &empty},
		{CriteriosSucesso: &emptyCriteria},
		{CriteriosSucesso: &blankCriteria},
		{ResultadoEsperado: &empty, CriteriosSucesso: &emptyCriteria},
	}
	for _, request := range tests {
		if err := ValidateNecessidadeUpdate(request); err == nil {
			t.Fatalf("accepted an explicit empty briefing: %+v", request)
		}
	}
}

func TestBriefingFieldsCanBeOmittedForLegacyUpdates(t *testing.T) {
	if err := ValidateNecessidadeCreate(&necessidade.CreateNecessidadeRequest{}); err != nil {
		t.Fatal("legacy creation is no longer compatible:", err)
	}
	title := "Título atualizado"
	if err := ValidateNecessidadeUpdate(&necessidade.UpdateNecessidadeRequest{Titulo: &title}); err != nil {
		t.Fatal("unrelated PATCH requires briefing fields:", err)
	}
	if err := ValidateNecessidadeCreate(&necessidade.CreateNecessidadeRequest{
		ResultadoEsperado: "   ", CriteriosSucesso: []string{"Critério"},
	}); err == nil {
		t.Fatal("accepted a newly submitted blank goal")
	}
}
