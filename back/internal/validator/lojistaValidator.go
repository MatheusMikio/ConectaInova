package validator

import "conectainova/internal/domain/schema"

func ValidateLojista(entity *schema.Lojista) error {
	if err := validateText(entity.Nome, "nome", 150); err != nil {
		return err
	}
	return validateText(entity.Segmento, "segmento", 150)
}
