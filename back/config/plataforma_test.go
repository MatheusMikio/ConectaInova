package config

import "testing"

func TestAssociationCanBeChangedPerDeployment(t *testing.T) {
	t.Setenv("ASSOCIACAO_NOME", " Associação Comercial A ")
	t.Setenv("ASSOCIACAO_REGIAO", " Região A ")
	if GetAssociacaoNome() != "Associação Comercial A" || GetAssociacaoRegiao() != "Região A" {
		t.Fatal("association configuration not applied")
	}
	t.Setenv("ASSOCIACAO_NOME", "Associação B")
	t.Setenv("ASSOCIACAO_REGIAO", "")
	if GetAssociacaoNome() != "Associação B" || GetAssociacaoRegiao() != "" {
		t.Fatal("deployment configuration was hardcoded")
	}
}
