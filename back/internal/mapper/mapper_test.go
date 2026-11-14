package mapper

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

func TestEntityMappers(t *testing.T) {
	now := time.Now().UTC()
	model := gorm.Model{ID: 7, CreatedAt: now, UpdatedAt: now}
	shop := LojistaToResponse(&schema.Lojista{Model: model, Nome: "Loja", Segmento: "Moda"})
	if shop.ID != 7 || shop.Nome != "Loja" || shop.Segmento != "Moda" || !shop.CreatedAt.Equal(now) {
		t.Fatal("incorrect lojista mapping")
	}
	company := StartupToResponse(&schema.Startup{Model: model, Nome: "Startup", Descricao: "Tecnologia"})
	if company.ID != 7 || company.Descricao != "Tecnologia" {
		t.Fatal("incorrect startup mapping")
	}
	need := NecessidadeToResponse(&schema.Necessidade{Model: model, LojistaID: 1, Titulo: "Estoque", Categoria: enum.CategoriaEstoque, Prioridade: enum.PrioridadeAlta, Status: enum.NecessidadeAberta})
	if need.LojistaID != 1 || need.Status != enum.NecessidadeAberta || need.Prioridade != enum.PrioridadeAlta {
		t.Fatal("incorrect necessidade mapping")
	}
	challenge := DesafioToResponse(&schema.Desafio{Model: model, NecessidadeID: 2, Status: enum.DesafioAberto})
	if challenge.NecessidadeID != 2 || challenge.Status != enum.DesafioAberto {
		t.Fatal("incorrect desafio mapping")
	}
	solution := &schema.Solucao{Model: model, StartupID: 3, Nome: "ERP", Categoria: enum.CategoriaEstoque, PalavrasChave: []string{"estoque"}, Startup: schema.Startup{Model: gorm.Model{ID: 3}, Nome: "Startup"}}
	response := SolucaoToResponse(solution)
	if response.StartupID != 3 || response.Categoria != enum.CategoriaEstoque {
		t.Fatal("incorrect solucao mapping")
	}
	response.PalavrasChave[0] = "outra"
	if solution.PalavrasChave[0] != "estoque" {
		t.Fatal("response shares mutable entity fields")
	}
	match := MatchToResponse(solution, 85)
	if match.Pontuacao != 85 || match.Solucao.ID != 7 || match.Startup.ID != 3 {
		t.Fatal("incorrect match mapping")
	}
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "DeletedAt") || strings.Contains(string(data), `"startup":`) {
		t.Fatal("mapper leaked persistence details")
	}
}

func TestListMappers(t *testing.T) {
	emptyLists := []any{
		LojistasToResponse(nil), NecessidadesToResponse(nil),
		StartupsToResponse(nil), SolucoesToResponse(nil), DesafiosToResponse(nil),
	}
	for _, list := range emptyLists {
		data, err := json.Marshal(list)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "[]" {
			t.Fatalf("expected [], got %s", data)
		}
	}
	if got := LojistasToResponse([]*schema.Lojista{{Nome: "A"}, {Nome: "B"}}); len(got) != 2 || got[1].Nome != "B" {
		t.Fatal("incorrect lojista list")
	}
	if len(NecessidadesToResponse([]*schema.Necessidade{{}})) != 1 {
		t.Fatal("incorrect necessidade list")
	}
	if len(StartupsToResponse([]*schema.Startup{{}})) != 1 {
		t.Fatal("incorrect startup list")
	}
	if len(SolucoesToResponse([]*schema.Solucao{{}})) != 1 {
		t.Fatal("incorrect solucao list")
	}
	if len(DesafiosToResponse([]*schema.Desafio{{}})) != 1 {
		t.Fatal("incorrect desafio list")
	}
}
