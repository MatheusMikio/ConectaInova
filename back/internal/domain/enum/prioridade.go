package enum

type Prioridade string

const (
	PrioridadeBaixa Prioridade = "baixa"
	PrioridadeMedia Prioridade = "media"
	PrioridadeAlta  Prioridade = "alta"
)

func (p Prioridade) IsValid() bool {
	switch p {
	case PrioridadeBaixa, PrioridadeMedia, PrioridadeAlta:
		return true
	default:
		return false
	}
}
