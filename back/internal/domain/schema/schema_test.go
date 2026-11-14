package schema

import (
	"sync"
	"testing"

	gormschema "gorm.io/gorm/schema"
)

func TestEntityMappings(t *testing.T) {
	entities := []struct {
		entity any
		table  string
	}{
		{&Lojista{}, "lojistas"}, {&Startup{}, "startups"},
		{&Necessidade{}, "necessidades"}, {&Solucao{}, "solucoes"},
		{&Desafio{}, "desafios"}, {&Match{}, "matches"},
		{&Avaliacao{}, "avaliacoes"},
	}
	cache := &sync.Map{}
	for _, tt := range entities {
		t.Run(tt.table, func(t *testing.T) {
			mapping, err := gormschema.Parse(tt.entity, cache, gormschema.NamingStrategy{})
			if err != nil {
				t.Fatal(err)
			}
			if mapping.Table != tt.table {
				t.Fatalf("table = %s, want %s", mapping.Table, tt.table)
			}
			for name, relation := range mapping.Relationships.Relations {
				if relation.Type != gormschema.BelongsTo {
					t.Errorf("%s: expected belongs-to", name)
				}
				for _, ref := range relation.References {
					if ref.PrimaryKey == nil || ref.ForeignKey == nil {
						t.Errorf("%s: unresolved foreign key", name)
					}
				}
			}
		})
	}
}
