package validator

import (
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
)

func ValidateNecessidade(entity *schema.Necessidade) error {
	if err := validateOptionalText(entity.Impacto, "impacto", 2000); err != nil {
		return err
	}
	if err := validateOptionalText(entity.ResultadoEsperado, "resultado_esperado", 2000); err != nil {
		return err
	}
	if err := validateOptionalText(entity.Restricoes, "restricoes", 2000); err != nil {
		return err
	}
	if err := validateCriteria(entity.CriteriosSucesso); err != nil {
		return err
	}
	switch entity.CanalVenda {
	case "", "fisica", "digital", "hibrida":
	default:
		return models.NewErrorMessage("canal_venda", "Canal de venda inválido.")
	}
	if err := validateText(entity.Titulo, "titulo", 200); err != nil {
		return err
	}
	if err := validateText(entity.Descricao, "descricao", 5000); err != nil {
		return err
	}
	if !entity.Categoria.IsValid() {
		return models.NewErrorMessage("categoria", "Categoria inválida.")
	}
	if !entity.Prioridade.IsValid() {
		return models.NewErrorMessage("prioridade", "Prioridade inválida.")
	}
	if !entity.Status.IsValid() {
		return models.NewErrorMessage("status", "Status da necessidade inválido.")
	}
	if entity.LojistaID == 0 {
		return models.NewErrorMessage("lojista_id", "Informe o lojista.")
	}
	return nil
}
