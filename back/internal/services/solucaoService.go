package services

import (
	"errors"

	"conectainova/internal/domain/dto/solucao"
	"conectainova/internal/domain/factory"
	"conectainova/internal/domain/models"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"conectainova/internal/validator"
	"gorm.io/gorm"
)

type ISolucaoService interface {
	CreateSolucao(request *solucao.CreateSolucaoRequest) (*solucao.SolucaoResponse, error)
	GetSolucoes(page, size uint) ([]*solucao.SolucaoResponse, error)
	UpdateSolucao(id uint, request *solucao.UpdateSolucaoRequest) (*solucao.SolucaoResponse, error)
}

type SolucaoService struct {
	SolucaoRepo repositories.ISolucaoRepository
	StartupRepo repositories.IStartupRepository
}

func NewSolucaoService(solucaoRepo repositories.ISolucaoRepository, startupRepo repositories.IStartupRepository) ISolucaoService {
	return &SolucaoService{SolucaoRepo: solucaoRepo, StartupRepo: startupRepo}
}

func (s *SolucaoService) CreateSolucao(request *solucao.CreateSolucaoRequest) (*solucao.SolucaoResponse, error) {
	entity := factory.NewSolucao(request)
	if err := validator.ValidateSolucao(entity); err != nil {
		return nil, err
	}
	if err := s.ensureStartup(entity.StartupID); err != nil {
		return nil, err
	}
	if err := s.SolucaoRepo.Create(entity); err != nil {
		return nil, err
	}
	return mapper.SolucaoToResponse(entity), nil
}

func (s *SolucaoService) GetSolucoes(page, size uint) ([]*solucao.SolucaoResponse, error) {
	entities, err := s.SolucaoRepo.GetAll(page, size)
	if err != nil {
		return nil, err
	}
	return mapper.SolucoesToResponse(entities), nil
}

func (s *SolucaoService) UpdateSolucao(id uint, request *solucao.UpdateSolucaoRequest) (*solucao.SolucaoResponse, error) {
	entity, err := s.SolucaoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	updated := factory.UpdateSolucao(entity, request)
	if err := validator.ValidateSolucao(updated); err != nil {
		return nil, err
	}
	if err := s.ensureStartup(updated.StartupID); err != nil {
		return nil, err
	}
	if err := s.SolucaoRepo.Update(updated); err != nil {
		return nil, err
	}
	return mapper.SolucaoToResponse(updated), nil
}

func (s *SolucaoService) ensureStartup(id uint) error {
	_, err := s.StartupRepo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.NewErrorMessage("startup_id", "Startup não encontrada.")
	}
	return err
}
