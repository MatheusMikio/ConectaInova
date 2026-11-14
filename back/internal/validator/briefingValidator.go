package validator

import (
	"strings"
	"unicode/utf8"

	"conectainova/internal/domain/models"
)

func validateOptionalText(value, property string, max int) error {
	if utf8.RuneCountInString(strings.TrimSpace(value)) > max {
		return models.NewErrorMessage(property, "Texto excede o tamanho permitido.")
	}
	return nil
}

func validateCriteria(values []string) error {
	if len(values) > 20 {
		return models.NewErrorMessage("criterios_sucesso", "Informe no máximo 20 critérios.")
	}
	for _, value := range values {
		if err := validateText(value, "criterios_sucesso", 300); err != nil {
			return err
		}
	}
	return nil
}
