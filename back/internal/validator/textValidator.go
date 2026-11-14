package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"conectainova/internal/domain/models"
)

func validateText(value, property string, max int) error {
	length := utf8.RuneCountInString(strings.TrimSpace(value))
	if length == 0 || length > max {
		return models.NewErrorMessage(property, fmt.Sprintf("Informe entre 1 e %d caracteres.", max))
	}
	return nil
}
