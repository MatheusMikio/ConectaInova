package services

import (
	"errors"
	"testing"

	"conectainova/internal/domain/schema"
	"conectainova/internal/repositories"
	"gorm.io/gorm"
)

type problemListRepoStub struct {
	repositories.INecessidadeRepository
	entities []*schema.Necessidade
	err      error
	page     uint
	size     uint
}

func (r *problemListRepoStub) GetAllWithLojista(page, size uint) ([]*schema.Necessidade, error) {
	r.page, r.size = page, size
	return r.entities, r.err
}

type linkedChallengesRepoStub struct {
	repositories.IDesafioRepository
	entities []*schema.Desafio
	err      error
	needID   uint
	page     uint
	size     uint
	calls    int
}

func (r *linkedChallengesRepoStub) GetByNecessidadeID(id, page, size uint) ([]*schema.Desafio, error) {
	r.needID, r.page, r.size = id, page, size
	r.calls++
	return r.entities, r.err
}

func TestPublicProblemListingIncludesDifferentStores(t *testing.T) {
	repo := &problemListRepoStub{entities: []*schema.Necessidade{
		{Model: gorm.Model{ID: 2}, LojistaID: 20, Lojista: schema.Lojista{Nome: "Loja B"}, Titulo: "Problema B"},
		{Model: gorm.Model{ID: 1}, LojistaID: 10, Lojista: schema.Lojista{Nome: "Loja A"}, Titulo: "Problema A"},
	}}
	// Listing is public and does not require an active store or any identity provider.
	service := NewNecessidadeService(repo, nil)
	result, err := service.GetNecessidades(2, 12)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 || result[0].LojistaNome != "Loja B" || result[1].LojistaNome != "Loja A" {
		t.Fatalf("missing stores in public listing: %+v", result)
	}
	if repo.page != 2 || repo.size != 12 {
		t.Fatal("pagination was not forwarded")
	}
	repo.entities = nil
	result, err = service.GetNecessidades(1, 12)
	if err != nil || result == nil || len(result) != 0 {
		t.Fatal("empty listing must return an empty array")
	}
	repo.err = errors.New("database unavailable")
	if _, err := service.GetNecessidades(1, 12); !errors.Is(err, repo.err) {
		t.Fatal("repository error was hidden")
	}
}

func TestLinkedChallengesValidateProblemAndForwardFilter(t *testing.T) {
	needs := &necessidadeRepoStub{}
	repo := &linkedChallengesRepoStub{entities: []*schema.Desafio{
		{Model: gorm.Model{ID: 5}, NecessidadeID: 3, Titulo: "Desafio"},
	}}
	service := NewDesafioService(repo, needs)
	result, err := service.GetDesafiosByNecessidade(3, 2, 12)
	if err != nil {
		t.Fatal(err)
	}
	if repo.needID != 3 || repo.page != 2 || repo.size != 12 || len(result) != 1 || result[0].NecessidadeID != 3 {
		t.Fatal("incorrect problem filter or pagination")
	}
	repo.entities = nil
	result, err = service.GetDesafiosByNecessidade(3, 1, 12)
	if err != nil || result == nil || len(result) != 0 {
		t.Fatal("expected empty array for a problem without challenges")
	}
	calls := repo.calls
	needs.err = gorm.ErrRecordNotFound
	if _, err := service.GetDesafiosByNecessidade(999, 1, 12); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatal("expected missing problem")
	}
	if calls != repo.calls {
		t.Fatal("queried challenges for a nonexistent problem")
	}
}
