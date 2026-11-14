package enum

type StatusDesafio string

const (
	DesafioAndamento StatusDesafio = "Em andamento"
	DesafioAberto    StatusDesafio = "Aberto"
	DesafioEncerrado StatusDesafio = "Encerrado"
)

func (s StatusDesafio) IsValid() bool {
	switch s {
	case DesafioAndamento, DesafioAberto, DesafioEncerrado:
		return true
	default:
		return false
	}
}
