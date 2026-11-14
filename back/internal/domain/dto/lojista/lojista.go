package lojista

import "conectainova/internal/domain/dto/common"

type CreateLojistaRequest struct {
	Nome     string `json:"nome"`
	Segmento string `json:"segmento"`
}

type UpdateLojistaRequest struct {
	Nome     *string `json:"nome"`
	Segmento *string `json:"segmento"`
}

type LojistaResponse struct {
	common.ResourceMetadata
	Nome     string `json:"nome"`
	Segmento string `json:"segmento"`
}
