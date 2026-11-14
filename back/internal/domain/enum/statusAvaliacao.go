package enum

type StatusAvaliacao string

const (
	AvaliacaoSolicitada     StatusAvaliacao = "solicitada"
	AvaliacaoEmAnalise      StatusAvaliacao = "em_analise"
	AvaliacaoPilotoProposto StatusAvaliacao = "piloto_proposto"
	AvaliacaoNaoAtendida    StatusAvaliacao = "nao_atendida"
	AvaliacaoConcluida      StatusAvaliacao = "concluida"
)

func (s StatusAvaliacao) IsValid() bool {
	switch s {
	case AvaliacaoSolicitada, AvaliacaoEmAnalise, AvaliacaoPilotoProposto, AvaliacaoNaoAtendida, AvaliacaoConcluida:
		return true
	}
	return false
}
