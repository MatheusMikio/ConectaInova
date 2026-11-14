package services

import (
	"errors"

	"conectainova/internal/domain/dto/desafio"
	"conectainova/internal/domain/factory"
	"conectainova/internal/domain/models"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"conectainova/internal/validator"
	"gorm.io/gorm"
)

type IDesafioService interface {
	CreateDesafio(request *desafio.CreateDesafioRequest) (*desafio.DesafioResponse, error)
	GetDesafios(page, size uint) ([]*desafio.DesafioResponse, error)
	GetDesafiosByNecessidade(necessidadeID, page, size uint) ([]*desafio.DesafioResponse, error)
	UpdateDesafio(id uint, request *desafio.UpdateDesafioRequest) (*desafio.DesafioResponse, error)
}

type DesafioService struct {
	DesafioRepo     repositories.IDesafioRepository
	NecessidadeRepo repositories.INecessidadeRepository
}

func NewDesafioService(desafioRepo repositories.IDesafioRepository, necessidadeRepo repositories.INecessidadeRepository) IDesafioService {
	return &DesafioService{DesafioRepo: desafioRepo, NecessidadeRepo: necessidadeRepo}
}

func (s *DesafioService) CreateDesafio(request *desafio.CreateDesafioRequest) (*desafio.DesafioResponse, error) {
	entity := factory.NewDesafio(request)
	if err := validator.ValidateDesafio(entity); err != nil {
		return nil, err
	}
	if err := s.ensureNecessidade(entity.NecessidadeID); err != nil {
		return nil, err
	}
	if err := s.DesafioRepo.Create(entity); err != nil {
		return nil, err
	}
	return mapper.DesafioToResponse(entity), nil
}

func (s *DesafioService) GetDesafios(page, size uint) ([]*desafio.DesafioResponse, error) {
	entities, err := s.DesafioRepo.GetAll(page, size)
	if err != nil {
		return nil, err
	}
	return mapper.DesafiosToResponse(entities), nil
}

func (s *DesafioService) GetDesafiosByNecessidade(necessidadeID, page, size uint) ([]*desafio.DesafioResponse, error) {
	if _, err := s.NecessidadeRepo.GetByID(necessidadeID); err != nil {
		return nil, err
	}
	entities, err := s.DesafioRepo.GetByNecessidadeID(necessidadeID, page, size)
	if err != nil {
		return nil, err
	}
	return mapper.DesafiosToResponse(entities), nil
}

func (s *DesafioService) UpdateDesafio(id uint, request *desafio.UpdateDesafioRequest) (*desafio.DesafioResponse, error) {
	entity, err := s.DesafioRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	updated := factory.UpdateDesafio(entity, request)
	if err := validator.ValidateDesafio(updated); err != nil {
		return nil, err
	}
	if err := s.ensureNecessidade(updated.NecessidadeID); err != nil {
		return nil, err
	}
	if err := s.DesafioRepo.Update(updated); err != nil {
		return nil, err
	}
	return mapper.DesafioToResponse(updated), nil
}

func (s *DesafioService) ensureNecessidade(id uint) error {
	_, err := s.NecessidadeRepo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.NewErrorMessage("necessidade_id", "Necessidade não encontrada.")
	}
	return err
}
