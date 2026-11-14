package validator

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
)

func TestInvalidEntities(t *testing.T) {
	tests := []struct {
		name     string
		validate func() error
		property string
	}{
		{"lojista name", func() error { return ValidateLojista(&schema.Lojista{Nome: " "}) }, "nome"},
		{"lojista segment", func() error { return ValidateLojista(&schema.Lojista{Nome: "Loja"}) }, "segmento"},
		{"startup description", func() error { return ValidateStartup(&schema.Startup{Nome: "Startup"}) }, "descricao"},
		{"necessidade category", func() error {
			return ValidateNecessidade(&schema.Necessidade{Titulo: "Título", Descricao: "Descrição"})
		}, "categoria"},
		{"necessidade priority", func() error {
			return ValidateNecessidade(&schema.Necessidade{Titulo: "Título", Descricao: "Descrição", Categoria: enum.CategoriaEstoque})
		}, "prioridade"},
		{"necessidade status", func() error {
			return ValidateNecessidade(&schema.Necessidade{Titulo: "Título", Descricao: "Descrição", Categoria: enum.CategoriaEstoque, Prioridade: enum.PrioridadeAlta})
		}, "status"},
		{"necessidade owner", func() error {
			return ValidateNecessidade(&schema.Necessidade{Titulo: "Título", Descricao: "Descrição", Categoria: enum.CategoriaEstoque, Prioridade: enum.PrioridadeAlta, Status: enum.NecessidadeAberta})
		}, "lojista_id"},
		{"solucao keywords", func() error {
			return ValidateSolucao(&schema.Solucao{Nome: "ERP", Descricao: "Descrição", Categoria: enum.CategoriaEstoque, PalavrasChave: []string{" "}})
		}, "palavras_chave"},
		{"solucao startup", func() error {
			return ValidateSolucao(&schema.Solucao{Nome: "ERP", Descricao: "Descrição", Categoria: enum.CategoriaEstoque})
		}, "startup_id"},
		{"desafio status", func() error { return ValidateDesafio(&schema.Desafio{Titulo: "Título", Descricao: "Descrição"}) }, "status"},
		{"desafio need", func() error {
			return ValidateDesafio(&schema.Desafio{Titulo: "Título", Descricao: "Descrição", Status: enum.DesafioAberto})
		}, "necessidade_id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var validation *models.ErrorMessage
			if err := tt.validate(); !errors.As(err, &validation) || validation.Property != tt.property {
				t.Fatalf("expected %s validation, got %v", tt.property, err)
			}
		})
	}
}

func TestTextLengthCountsCharacters(t *testing.T) {
	if err := validateText(strings.Repeat("á", 150), "nome", 150); err != nil {
		t.Fatal("counted bytes instead of characters")
	}
	if err := validateText(strings.Repeat("á", 151), "nome", 150); err == nil {
		t.Fatal("accepted overly long text")
	}
}

func TestValidatorsDoNotMutateEntities(t *testing.T) {
	entity := &schema.Solucao{Nome: " ERP ", Descricao: "Descrição", Categoria: enum.CategoriaEstoque, StartupID: 1, PalavrasChave: []string{" estoque "}}
	before := *entity
	before.PalavrasChave = append([]string{}, entity.PalavrasChave...)
	if err := ValidateSolucao(entity); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, *entity) {
		t.Fatal("validator mutated entity")
	}
}
