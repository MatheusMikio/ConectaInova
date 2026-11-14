package mapper

import (
	"conectainova/internal/domain/dto/common"
	"conectainova/internal/domain/dto/startup"
	"conectainova/internal/domain/schema"
)

func StartupToResponse(entity *schema.Startup) *startup.StartupResponse {
	return &startup.StartupResponse{
		ResourceMetadata: common.ResourceMetadata{ID: entity.ID, CreatedAt: entity.CreatedAt, UpdatedAt: entity.UpdatedAt},
		Nome:             entity.Nome,
		Descricao:        entity.Descricao,
	}
}

func StartupsToResponse(entities []*schema.Startup) []*startup.StartupResponse {
	response := make([]*startup.StartupResponse, 0, len(entities))
	for _, entity := range entities {
		response = append(response, StartupToResponse(entity))
	}
	return response
}
