package validator

import (
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
)

func ValidateSolucao(entity *schema.Solucao) error {
	if err := validateText(entity.Nome, "nome", 150); err != nil {
		return err
	}
	if err := validateText(entity.Descricao, "descricao", 5000); err != nil {
		return err
	}
	if !entity.Categoria.IsValid() {
		return models.NewErrorMessage("categoria", "Categoria inválida.")
	}
	if len(entity.PalavrasChave) > 30 {
		return models.NewErrorMessage("palavras_chave", "Informe no máximo 30 palavras-chave.")
	}
	for _, keyword := range entity.PalavrasChave {
		if err := validateText(keyword, "palavras_chave", 100); err != nil {
			return err
		}
	}
	if entity.StartupID == 0 {
		return models.NewErrorMessage("startup_id", "Informe a startup.")
	}
	return nil
}
