package factory

import (
	"conectainova/internal/domain/dto/desafio"
	"conectainova/internal/domain/dto/necessidade"
	"testing"
)

func TestBriefingFieldsAndCriteriaAreNormalizedAndCopied(t *testing.T) {
	input := &necessidade.CreateNecessidadeRequest{
		Impacto: " Perda de vendas ", ResultadoEsperado: " Vender sem rupturas ",
		Restricoes: " Operação simples ", CanalVenda: "hibrida", CriteriosSucesso: []string{" Receber alertas "},
	}
	entity := NewNecessidade(input)
	if entity.Impacto != "Perda de vendas" || entity.ResultadoEsperado != "Vender sem rupturas" || entity.CriteriosSucesso[0] != "Receber alertas" {
		t.Fatal("briefing fields were lost")
	}
	updated := UpdateNecessidade(entity, &necessidade.UpdateNecessidadeRequest{})
	updated.CriteriosSucesso[0] = "Outro"
	if entity.CriteriosSucesso[0] != "Receber alertas" || input.CriteriosSucesso[0] != " Receber alertas " {
		t.Fatal("shared mutable criteria")
	}
	days := 14
	challenge := NewDesafio(&desafio.CreateDesafioRequest{CriteriosSucesso: []string{" Medir resultado "}, PrazoDias: &days})
	if challenge.PrazoDias != 14 || challenge.CriteriosSucesso[0] != "Medir resultado" {
		t.Fatal("challenge details were lost")
	}
	copy := UpdateDesafio(challenge, &desafio.UpdateDesafioRequest{})
	copy.CriteriosSucesso[0] = "Outro"
	if challenge.CriteriosSucesso[0] != "Medir resultado" {
		t.Fatal("challenge criteria were not copied")
	}
}
