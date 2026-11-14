package services

import (
	"conectainova/internal/domain/dto/startup"
	"conectainova/internal/domain/factory"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"conectainova/internal/validator"
)

type IStartupService interface {
	CreateStartup(request *startup.CreateStartupRequest) (*startup.StartupResponse, error)
	GetStartups(page, size uint) ([]*startup.StartupResponse, error)
	UpdateStartup(id uint, request *startup.UpdateStartupRequest) (*startup.StartupResponse, error)
}

type StartupService struct {
	StartupRepo repositories.IStartupRepository
}

func NewStartupService(startupRepo repositories.IStartupRepository) IStartupService {
	return &StartupService{StartupRepo: startupRepo}
}

func (s *StartupService) CreateStartup(request *startup.CreateStartupRequest) (*startup.StartupResponse, error) {
	entity := factory.NewStartup(request)
	if err := validator.ValidateStartup(entity); err != nil {
		return nil, err
	}
	if err := s.StartupRepo.Create(entity); err != nil {
		return nil, err
	}
	return mapper.StartupToResponse(entity), nil
}

func (s *StartupService) GetStartups(page, size uint) ([]*startup.StartupResponse, error) {
	entities, err := s.StartupRepo.GetAll(page, size)
	if err != nil {
		return nil, err
	}
	return mapper.StartupsToResponse(entities), nil
}

func (s *StartupService) UpdateStartup(id uint, request *startup.UpdateStartupRequest) (*startup.StartupResponse, error) {
	entity, err := s.StartupRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	updated := factory.UpdateStartup(entity, request)
	if err := validator.ValidateStartup(updated); err != nil {
		return nil, err
	}
	if err := s.StartupRepo.Update(updated); err != nil {
		return nil, err
	}
	return mapper.StartupToResponse(updated), nil
}
