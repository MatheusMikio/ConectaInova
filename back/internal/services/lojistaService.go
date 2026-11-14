package services

import (
	"conectainova/internal/domain/dto/lojista"
	"conectainova/internal/domain/factory"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"conectainova/internal/validator"
)

type ILojistaService interface {
	CreateLojista(request *lojista.CreateLojistaRequest) (*lojista.LojistaResponse, error)
	GetLojista(id uint) (*lojista.LojistaResponse, error)
	UpdateLojista(id uint, request *lojista.UpdateLojistaRequest) (*lojista.LojistaResponse, error)
}

type LojistaService struct {
	LojistaRepo repositories.ILojistaRepository
}

func NewLojistaService(lojistaRepo repositories.ILojistaRepository) ILojistaService {
	return &LojistaService{LojistaRepo: lojistaRepo}
}

func (s *LojistaService) CreateLojista(request *lojista.CreateLojistaRequest) (*lojista.LojistaResponse, error) {
	entity := factory.NewLojista(request)
	if err := validator.ValidateLojista(entity); err != nil {
		return nil, err
	}
	if err := s.LojistaRepo.Create(entity); err != nil {
		return nil, err
	}
	return mapper.LojistaToResponse(entity), nil
}

func (s *LojistaService) GetLojista(id uint) (*lojista.LojistaResponse, error) {
	entity, err := s.LojistaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return mapper.LojistaToResponse(entity), nil
}

func (s *LojistaService) UpdateLojista(id uint, request *lojista.UpdateLojistaRequest) (*lojista.LojistaResponse, error) {
	entity, err := s.LojistaRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	updated := factory.UpdateLojista(entity, request)
	if err := validator.ValidateLojista(updated); err != nil {
		return nil, err
	}
	if err := s.LojistaRepo.Update(updated); err != nil {
		return nil, err
	}
	return mapper.LojistaToResponse(updated), nil
}
