package enum

type Prioridade string

const (
	PrioridadeBaixa Prioridade = "Baixa"
	PrioridadeMedia Prioridade = "Media"
	PrioridadeAlta  Prioridade = "Alta"
)

func (p Prioridade) IsValid() bool {
	switch p {
	case PrioridadeBaixa, PrioridadeMedia, PrioridadeAlta:
		return true
	default:
		return false
	}
}
