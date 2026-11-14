package enum

type Categoria string

const (
	CategoriaEstoque     Categoria = "estoque"
	CategoriaVendas      Categoria = "vendas"
	CategoriaAtendimento Categoria = "atendimento"
	CategoriaFinanceiro  Categoria = "financeiro"
	CategoriaMarketing   Categoria = "marketing"
	CategoriaLogistica   Categoria = "logistica"
)

func (c Categoria) IsValid() bool {
	switch c {
	case CategoriaEstoque, CategoriaVendas, CategoriaAtendimento, CategoriaFinanceiro, CategoriaMarketing, CategoriaLogistica:
		return true
	default:
		return false
	}
}
