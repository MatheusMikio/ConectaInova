package startup

import "conectainova/internal/domain/dto/common"

type CreateStartupRequest struct {
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}

type UpdateStartupRequest struct {
	Nome      *string `json:"nome"`
	Descricao *string `json:"descricao"`
}

type StartupResponse struct {
	common.ResourceMetadata
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
}
