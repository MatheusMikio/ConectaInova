package validator

import (
	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/models"
)

// Older records may have no briefing. Omitted fields remain supported, but an
// explicit briefing submission must contain a goal and at least one criterion.
func ValidateNecessidadeCreate(request *necessidade.CreateNecessidadeRequest) error {
	if request.ResultadoEsperado == "" && request.CriteriosSucesso == nil {
		return nil
	}
	if err := validateText(request.ResultadoEsperado, "resultado_esperado", 2000); err != nil {
		return err
	}
	return validateRequiredCriteria(request.CriteriosSucesso)
}

func ValidateNecessidadeUpdate(request *necessidade.UpdateNecessidadeRequest) error {
	if request.ResultadoEsperado != nil {
		if err := validateText(*request.ResultadoEsperado, "resultado_esperado", 2000); err != nil {
			return err
		}
	}
	if request.CriteriosSucesso != nil {
		if err := validateRequiredCriteria(*request.CriteriosSucesso); err != nil {
			return err
		}
	}
	return nil
}

func validateRequiredCriteria(criteria []string) error {
	if len(criteria) == 0 {
		return models.NewErrorMessage("criterios_sucesso", "Defina pelo menos um critério de sucesso.")
	}
	return validateCriteria(criteria)
}
