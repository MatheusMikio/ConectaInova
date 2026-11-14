package services

import (
	"errors"
	"strings"
	"testing"

	dto "conectainova/internal/domain/dto/avaliacao"
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
	"conectainova/internal/repositories"
	"gorm.io/gorm"
)

type avaliacaoRepoStub struct {
	repositories.IAvaliacaoRepository
	saved    *schema.Avaliacao
	need     schema.Necessidade
	solution schema.Solucao
	created  int
}

func (r *avaliacaoRepoStub) CreateOrGet(entity *schema.Avaliacao) (*schema.Avaliacao, bool, error) {
	if r.saved != nil {
		return r.saved, false, nil
	}
	entity.ID = 1
	entity.Necessidade, entity.Solucao = r.need, r.solution
	r.saved = entity
	r.created++
	return entity, true, nil
}

func (r *avaliacaoRepoStub) GetByID(id uint) (*schema.Avaliacao, error) {
	if r.saved == nil {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *r.saved
	return &copy, nil
}

func (r *avaliacaoRepoStub) Update(entity *schema.Avaliacao) error { r.saved = entity; return nil }

func TestEvaluationRequiresQualifiedBriefing(t *testing.T) {
	need := &necessidadeRepoStub{entity: schema.Necessidade{Categoria: enum.CategoriaEstoque}}
	solution := &solucaoRepoStub{entity: schema.Solucao{Categoria: enum.CategoriaEstoque}}
	repo := &avaliacaoRepoStub{}
	service := NewAvaliacaoService(repo, need, solution, "Associação A", "Região")
	request := &dto.CreateAvaliacaoRequest{NecessidadeID: 1, SolucaoID: 2}
	_, _, err := service.CreateAvaliacao(request)
	var validation *models.ErrorMessage
	if !errors.As(err, &validation) || validation.Property != "briefing" {
		t.Fatalf("expected briefing error: %v", err)
	}
	if repo.created != 0 {
		t.Fatal("created evaluation without sufficient context")
	}
	need.entity.ResultadoEsperado = "Diminuir falta de produtos"
	need.entity.CriteriosSucesso = []string{"Receber alertas"}
	solution.entity.Categoria = enum.CategoriaVendas
	if _, _, err = service.CreateAvaliacao(request); err == nil {
		t.Fatal("accepted incompatible solution")
	}
	solution.entity.Categoria = enum.CategoriaEstoque
	repo.need, repo.solution = need.entity, solution.entity
	result, created, err := service.CreateAvaliacao(request)
	if err != nil || !created || result.Status != enum.AvaliacaoSolicitada {
		t.Fatalf("unexpected result: %+v %v", result, err)
	}
	if result.Associacao != "Associação A" || result.Regiao != "Região" {
		t.Fatal("missing association context")
	}
	_, created, err = service.CreateAvaliacao(request)
	if err != nil || created || repo.created != 1 {
		t.Fatal("duplicate evaluation")
	}
}

func TestEvaluationNeedsAnExplicitNextStep(t *testing.T) {
	repo := &avaliacaoRepoStub{saved: &schema.Avaliacao{
		Model: gorm.Model{ID: 1}, NecessidadeID: 1, SolucaoID: 2, Status: enum.AvaliacaoSolicitada,
	}}
	service := NewAvaliacaoService(repo, nil, nil, "", "")
	proposed := enum.AvaliacaoPilotoProposto
	if _, err := service.UpdateAvaliacao(1, &dto.UpdateAvaliacaoRequest{Status: &proposed}); err == nil {
		t.Fatal("accepted pilot proposal without next steps")
	}
	if repo.saved.Status != enum.AvaliacaoSolicitada {
		t.Fatal("invalid update mutated entity")
	}
	next := "Agendar avaliação de implantação e confirmar condições com a loja."
	result, err := service.UpdateAvaliacao(1, &dto.UpdateAvaliacaoRequest{Status: &proposed, Resposta: &next})
	if err != nil || result.Status != proposed || result.Resposta != next {
		t.Fatal("did not persist next steps")
	}
}

func TestMatchingExplanationReportsEvidenceAndLimits(t *testing.T) {
	need := schema.Necessidade{Titulo: "Controlar estoque", Descricao: "Reposição de produtos", Categoria: enum.CategoriaEstoque}
	reasons, found := explainMatch(need, schema.Solucao{PalavrasChave: []string{"estoque", "REPOSICAO", "estoque", "caixa"}})
	if len(reasons) != 2 || len(found) != 2 {
		t.Fatalf("wrong evidence: %v %v", reasons, found)
	}
	if strings.Contains(strings.Join(found, ","), "caixa") {
		t.Fatal("invented a matching keyword")
	}
	reasons, found = explainMatch(need, schema.Solucao{PalavrasChave: []string{"caixa"}})
	if len(found) != 0 || !strings.Contains(reasons[1], "apenas a categoria") {
		t.Fatal("did not disclose category-only matching")
	}
}
