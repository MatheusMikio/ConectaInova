package validator

import "conectainova/internal/domain/schema"

func ValidateStartup(entity *schema.Startup) error {
	if err := validateText(entity.Nome, "nome", 150); err != nil {
		return err
	}
	return validateText(entity.Descricao, "descricao", 5000)
}
