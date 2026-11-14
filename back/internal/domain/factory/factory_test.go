package factory

import (
	"testing"

	"conectainova/internal/domain/dto/desafio"
	"conectainova/internal/domain/dto/lojista"
	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/dto/startup"
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/schema"
	"gorm.io/gorm"
)

func TestCreateFactoriesNormalizeAndSetDefaults(t *testing.T) {
	shop := NewLojista(&lojista.CreateLojistaRequest{Nome: " Loja ", Segmento: " Moda "})
	if shop.Nome != "Loja" || shop.Segmento != "Moda" {
		t.Fatal("lojista was not normalized")
	}
	company := NewStartup(&startup.CreateStartupRequest{Nome: " Startup ", Descricao: " Tecnologia "})
	if company.Nome != "Startup" || company.Descricao != "Tecnologia" {
		t.Fatal("startup was not normalized")
	}
	need := NewNecessidade(&necessidade.CreateNecessidadeRequest{LojistaID: 1, Titulo: " Estoque ", Categoria: enum.CategoriaEstoque, Prioridade: enum.PrioridadeAlta})
	if need.Status != enum.NecessidadeAberta || need.Titulo != "Estoque" || need.LojistaID != 1 {
		t.Fatal("incorrect necessidade defaults")
	}
	challenge := NewDesafio(&desafio.CreateDesafioRequest{NecessidadeID: 1, Titulo: " Inovar "})
	if challenge.Status != enum.DesafioAberto || challenge.Titulo != "Inovar" {
		t.Fatal("incorrect desafio defaults")
	}
	status := enum.DesafioAndamento
	challenge = NewDesafio(&desafio.CreateDesafioRequest{Status: &status})
	if challenge.Status != status {
		t.Fatal("explicit status was ignored")
	}
	keywords := []string{" estoque "}
	solution := NewSolucao(&solucao.CreateSolucaoRequest{Nome: " ERP ", PalavrasChave: keywords})
	if solution.Nome != "ERP" || solution.PalavrasChave[0] != "estoque" {
		t.Fatal("solucao was not normalized")
	}
	if keywords[0] != " estoque " {
		t.Fatal("factory changed request")
	}
	empty := NewSolucao(&solucao.CreateSolucaoRequest{})
	if empty.PalavrasChave == nil {
		t.Fatal("expected empty array")
	}
}

func TestUpdateFactoriesPreserveOmittedFieldsAndOriginalEntity(t *testing.T) {
	name := " Atualizado "
	shop := &schema.Lojista{Model: gorm.Model{ID: 7}, Nome: "Original", Segmento: "Moda"}
	updatedShop := UpdateLojista(shop, &lojista.UpdateLojistaRequest{Nome: &name})
	if shop.Nome != "Original" || updatedShop.Nome != "Atualizado" || updatedShop.Segmento != "Moda" || updatedShop.ID != 7 {
		t.Fatal("incorrect lojista PATCH")
	}
	company := &schema.Startup{Nome: "Original", Descricao: "Descrição"}
	updatedCompany := UpdateStartup(company, &startup.UpdateStartupRequest{Nome: &name})
	if company.Nome != "Original" || updatedCompany.Descricao != "Descrição" {
		t.Fatal("incorrect startup PATCH")
	}
	status := enum.NecessidadeResolvida
	need := &schema.Necessidade{Titulo: "Original", Status: enum.NecessidadeAberta}
	updatedNeed := UpdateNecessidade(need, &necessidade.UpdateNecessidadeRequest{Status: &status})
	if need.Status != enum.NecessidadeAberta || updatedNeed.Status != status || updatedNeed.Titulo != "Original" {
		t.Fatal("incorrect necessidade PATCH")
	}
	challenge := &schema.Desafio{Titulo: "Original", Status: enum.DesafioAberto}
	updatedChallenge := UpdateDesafio(challenge, &desafio.UpdateDesafioRequest{Titulo: &name})
	if challenge.Titulo != "Original" || updatedChallenge.Titulo != "Atualizado" || updatedChallenge.Status != enum.DesafioAberto {
		t.Fatal("incorrect desafio PATCH")
	}
	solution := &schema.Solucao{Nome: "Original", PalavrasChave: []string{"estoque"}}
	updatedSolution := UpdateSolucao(solution, &solucao.UpdateSolucaoRequest{Nome: &name})
	updatedSolution.PalavrasChave[0] = "outra"
	if solution.PalavrasChave[0] != "estoque" {
		t.Fatal("PATCH shares mutable keywords")
	}
	empty := []string{}
	cleared := UpdateSolucao(solution, &solucao.UpdateSolucaoRequest{PalavrasChave: &empty})
	if cleared.PalavrasChave == nil || len(cleared.PalavrasChave) != 0 {
		t.Fatal("explicit empty array was ignored")
	}
}
