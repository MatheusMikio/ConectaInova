package enum

type StatusDesafio string

const (
	DesafioAndamento StatusDesafio = "em_andamento"
	DesafioAberto    StatusDesafio = "aberto"
	DesafioEncerrado StatusDesafio = "encerrado"
)

func (s StatusDesafio) IsValid() bool {
	switch s {
	case DesafioAndamento, DesafioAberto, DesafioEncerrado:
		return true
	default:
		return false
	}
}
