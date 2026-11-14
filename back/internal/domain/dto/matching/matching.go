package matching

import (
	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/dto/startup"
)

type MatchResponse struct {
	Solucao             solucao.SolucaoResponse `json:"solucao"`
	Startup             startup.StartupResponse `json:"startup"`
	Pontuacao           float64                 `json:"pontuacao"`
	Motivos             []string                `json:"motivos"`
	PalavrasEncontradas []string                `json:"palavras_encontradas"`
	PontosAValidar      []string                `json:"pontos_a_validar"`
}

type MatchingResponse struct {
	NecessidadeID  uint            `json:"necessidade_id"`
	Matches        []MatchResponse `json:"matches"`
	SugerirDesafio bool            `json:"sugerir_desafio"`
}
