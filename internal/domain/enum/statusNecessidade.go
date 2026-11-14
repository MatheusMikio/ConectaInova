package enum

type StatusNecessidade string

const (
	NecessidadeAberta    StatusNecessidade = "aberta"
	NecessidadeEmAnalise StatusNecessidade = "em_analise"
	NecessidadeResolvida StatusNecessidade = "resolvida"
)

func (s StatusNecessidade) IsValid() bool {
	switch s {
	case NecessidadeAberta,
		NecessidadeEmAnalise,
		NecessidadeResolvida:
		return true
	default:
		return false
	}
}
