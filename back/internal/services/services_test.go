package services

import (
	"errors"
	"testing"

	desafiodto "conectainova/internal/domain/dto/desafio"
	lojistadto "conectainova/internal/domain/dto/lojista"
	necessidadedto "conectainova/internal/domain/dto/necessidade"
	solucaodto "conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
	"conectainova/internal/repositories"
	"gorm.io/gorm"
)

type lojistaRepoStub struct {
	repositories.ILojistaRepository
	entity schema.Lojista
	err    error
	writes int
}

func (r *lojistaRepoStub) GetByID(id uint) (*schema.Lojista, error) {
	copy := r.entity
	return &copy, r.err
}
func (r *lojistaRepoStub) Create(entity *schema.Lojista) error {
	r.writes++
	entity.ID = 1
	r.entity = *entity
	return r.err
}
func (r *lojistaRepoStub) Update(entity *schema.Lojista) error {
	r.writes++
	r.entity = *entity
	return r.err
}

type necessidadeRepoStub struct {
	repositories.INecessidadeRepository
	entity schema.Necessidade
	err    error
	writes int
}

func (r *necessidadeRepoStub) GetByID(id uint) (*schema.Necessidade, error) {
	copy := r.entity
	return &copy, r.err
}
func (r *necessidadeRepoStub) Create(entity *schema.Necessidade) error {
	r.writes++
	entity.ID = 1
	r.entity = *entity
	return r.err
}
func (r *necessidadeRepoStub) Update(entity *schema.Necessidade) error {
	r.writes++
	r.entity = *entity
	return r.err
}

type startupRepoStub struct {
	repositories.IStartupRepository
	err error
}

func (r *startupRepoStub) GetByID(id uint) (*schema.Startup, error) {
	return &schema.Startup{}, r.err
}

type solucaoRepoStub struct {
	repositories.ISolucaoRepository
	entities []*schema.Solucao
	entity   schema.Solucao
	err      error
	writes   int
}

func (r *solucaoRepoStub) GetByID(id uint) (*schema.Solucao, error) {
	copy := r.entity
	copy.PalavrasChave = append([]string{}, copy.PalavrasChave...)
	return &copy, r.err
}
func (r *solucaoRepoStub) GetByCategoria(categoria enum.Categoria) ([]*schema.Solucao, error) {
	return r.entities, r.err
}
func (r *solucaoRepoStub) Update(entity *schema.Solucao) error {
	r.writes++
	r.entity = *entity
	return r.err
}

type desafioRepoStub struct {
	repositories.IDesafioRepository
	entity schema.Desafio
}

func (r *desafioRepoStub) Create(entity *schema.Desafio) error {
	entity.ID = 1
	r.entity = *entity
	return nil
}

func TestLojistaCreateAndPartialUpdate(t *testing.T) {
	repo := &lojistaRepoStub{}
	service := NewLojistaService(repo)
	created, err := service.CreateLojista(&lojistadto.CreateLojistaRequest{Nome: " Loja ", Segmento: "Moda"})
	if err != nil {
		t.Fatal(err)
	}
	if created.Nome != "Loja" || created.ID != 1 {
		t.Fatalf("unexpected result: %+v", created)
	}
	name := "Novo nome"
	updated, err := service.UpdateLojista(1, &lojistadto.UpdateLojistaRequest{Nome: &name})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Nome != name || updated.Segmento != "Moda" {
		t.Fatalf("PATCH erased fields: %+v", updated)
	}
	blank := " "
	_, err = service.UpdateLojista(1, &lojistadto.UpdateLojistaRequest{Nome: &blank})
	var validation *models.ErrorMessage
	if !errors.As(err, &validation) || validation.Property != "nome" {
		t.Fatalf("expected nome validation, got %v", err)
	}
	if repo.writes != 2 {
		t.Fatal("invalid update persisted")
	}
}

func TestNecessidadeDefaultsAndReferences(t *testing.T) {
	repo := &necessidadeRepoStub{}
	lojistas := &lojistaRepoStub{}
	service := NewNecessidadeService(repo, lojistas)
	request := &necessidadedto.CreateNecessidadeRequest{
		LojistaID: 1, Titulo: "Estoque", Descricao: "Automatizar estoque",
		Categoria: enum.CategoriaEstoque, Prioridade: enum.PrioridadeAlta,
	}
	created, err := service.CreateNecessidade(request)
	if err != nil {
		t.Fatal(err)
	}
	if created.Status != enum.NecessidadeAberta {
		t.Fatalf("status = %s", created.Status)
	}
	status := enum.NecessidadeEmAnalise
	updated, err := service.UpdateNecessidade(created.ID, &necessidadedto.UpdateNecessidadeRequest{Status: &status})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Titulo != request.Titulo || updated.Status != status {
		t.Fatalf("unexpected PATCH: %+v", updated)
	}
	request.Categoria = "invalida"
	if _, err := service.CreateNecessidade(request); err == nil {
		t.Fatal("accepted invalid category")
	}
	request.Categoria = enum.CategoriaEstoque
	lojistas.err = gorm.ErrRecordNotFound
	_, err = service.CreateNecessidade(request)
	var validation *models.ErrorMessage
	if !errors.As(err, &validation) || validation.Property != "lojista_id" {
		t.Fatalf("expected invalid reference, got %v", err)
	}
	lojistas.err = errors.New("database unavailable")
	_, err = service.CreateNecessidade(request)
	if !errors.Is(err, lojistas.err) {
		t.Fatal("database failure was hidden")
	}
}

func TestSolucaoPatchCanClearKeywords(t *testing.T) {
	repo := &solucaoRepoStub{entity: schema.Solucao{
		Model: gorm.Model{ID: 1}, StartupID: 1, Nome: "Gestão", Descricao: "Controlar estoque",
		Categoria: enum.CategoriaEstoque, PalavrasChave: []string{"estoque"},
	}}
	service := NewSolucaoService(repo, &startupRepoStub{})
	empty := []string{}
	result, err := service.UpdateSolucao(1, &solucaodto.UpdateSolucaoRequest{PalavrasChave: &empty})
	if err != nil {
		t.Fatal(err)
	}
	if result.PalavrasChave == nil || len(result.PalavrasChave) != 0 {
		t.Fatal("expected empty JSON array")
	}
	if result.Nome != "Gestão" {
		t.Fatal("PATCH erased name")
	}
}

func TestDesafioDefaults(t *testing.T) {
	service := NewDesafioService(&desafioRepoStub{}, &necessidadeRepoStub{})
	result, err := service.CreateDesafio(&desafiodto.CreateDesafioRequest{
		NecessidadeID: 1, Titulo: "Inovar", Descricao: "Criar uma nova solução",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != enum.DesafioAberto {
		t.Fatal("expected default aberto")
	}
}

func TestScore(t *testing.T) {
	need := schema.Necessidade{
		Titulo: "Previsão de vendas", Descricao: "Controlar estoque e reposição",
		Categoria: enum.CategoriaEstoque,
	}
	tests := []struct {
		name     string
		category enum.Categoria
		keywords []string
		want     float64
	}{
		{"other category", enum.CategoriaVendas, []string{"estoque"}, 0},
		{"no keywords", enum.CategoriaEstoque, nil, 70},
		{"accent and case", enum.CategoriaEstoque, []string{"PREVISAO", "reposição"}, 100},
		{"partial", enum.CategoriaEstoque, []string{"estoque", "caixa"}, 85},
		{"whole words", enum.CategoriaEstoque, []string{"venda"}, 70},
		{"duplicates", enum.CategoriaEstoque, []string{"Estoque", "estoque", ""}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Score(need, schema.Solucao{Categoria: tt.category, PalavrasChave: tt.keywords})
			if got != tt.want {
				t.Fatalf("score = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchingSortEmptyAndMissing(t *testing.T) {
	necessidades := &necessidadeRepoStub{entity: schema.Necessidade{Categoria: enum.CategoriaEstoque, Titulo: "Estoque"}}
	startup := schema.Startup{Model: gorm.Model{ID: 1}, Nome: "Startup"}
	solucoes := &solucaoRepoStub{entities: []*schema.Solucao{
		{Model: gorm.Model{ID: 2}, Startup: startup, Categoria: enum.CategoriaEstoque},
		{Model: gorm.Model{ID: 3}, Startup: startup, Categoria: enum.CategoriaEstoque, PalavrasChave: []string{"estoque"}},
		{Model: gorm.Model{ID: 1}, Startup: startup, Categoria: enum.CategoriaEstoque},
		{Model: gorm.Model{ID: 4}, Categoria: enum.CategoriaEstoque},
	}}
	service := NewMatchingService(necessidades, solucoes)
	result, err := service.GetMatches(1)
	if err != nil {
		t.Fatal(err)
	}
	if result.SugerirDesafio || len(result.Matches) != 3 {
		t.Fatalf("unexpected matches: %+v", result)
	}
	for i, want := range []uint{3, 1, 2} {
		if result.Matches[i].Solucao.ID != want {
			t.Fatal("incorrect ordering")
		}
	}
	solucoes.entities = nil
	result, err = service.GetMatches(1)
	if err != nil || !result.SugerirDesafio || result.Matches == nil {
		t.Fatal("expected challenge suggestion and empty array")
	}
	necessidades.err = gorm.ErrRecordNotFound
	if _, err := service.GetMatches(999); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatal("expected missing need")
	}
}
