package config

import (
	"os"
	"strings"
)

func GetAssociacaoNome() string {
	if value := strings.TrimSpace(os.Getenv("ASSOCIACAO_NOME")); value != "" {
		return value
	}
	return "Rede de associações comerciais"
}

func GetAssociacaoRegiao() string { return strings.TrimSpace(os.Getenv("ASSOCIACAO_REGIAO")) }
