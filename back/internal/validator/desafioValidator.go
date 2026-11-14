package validator

import (
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
)

func ValidateDesafio(entity *schema.Desafio) error {
	if err := validateText(entity.Titulo, "titulo", 200); err != nil {
		return err
	}
	if err := validateText(entity.Descricao, "descricao", 5000); err != nil {
		return err
	}
	if !entity.Status.IsValid() {
		return models.NewErrorMessage("status", "Status do desafio inválido.")
	}
	if entity.NecessidadeID == 0 {
		return models.NewErrorMessage("necessidade_id", "Informe a necessidade.")
	}
	if entity.PrazoDias < 1 || entity.PrazoDias > 365 {
		return models.NewErrorMessage("prazo_dias", "Informe um prazo entre 1 e 365 dias.")
	}
	return validateCriteria(entity.CriteriosSucesso)
}
