package mapper

import (
	"conectainova/internal/domain/dto/matching"
	"conectainova/internal/domain/schema"
)

func MatchToResponse(solucao *schema.Solucao, pontuacao float64) *matching.MatchResponse {
	return &matching.MatchResponse{
		Solucao:   *SolucaoToResponse(solucao),
		Startup:   *StartupToResponse(&solucao.Startup),
		Pontuacao: pontuacao,
	}
}
