package validator

import (
	"conectainova/internal/domain/enum"
	"conectainova/internal/domain/models"
	"conectainova/internal/domain/schema"
)

func ValidateAvaliacao(entity *schema.Avaliacao) error {
	if entity.NecessidadeID == 0 {
		return models.NewErrorMessage("necessidade_id", "Selecione um problema.")
	}
	if entity.SolucaoID == 0 {
		return models.NewErrorMessage("solucao_id", "Selecione uma solução.")
	}
	if !entity.Status.IsValid() {
		return models.NewErrorMessage("status", "Status de avaliação inválido.")
	}
	if err := validateOptionalText(entity.Mensagem, "mensagem", 2000); err != nil {
		return err
	}
	if err := validateOptionalText(entity.Resposta, "resposta", 4000); err != nil {
		return err
	}
	if entity.Status == enum.AvaliacaoPilotoProposto || entity.Status == enum.AvaliacaoNaoAtendida || entity.Status == enum.AvaliacaoConcluida {
		return validateText(entity.Resposta, "resposta", 4000)
	}
	return nil
}
