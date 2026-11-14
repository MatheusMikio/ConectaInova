package services

import (
	"errors"

	"conectainova/internal/domain/dto/necessidade"
	"conectainova/internal/domain/factory"
	"conectainova/internal/domain/models"
	"conectainova/internal/mapper"
	"conectainova/internal/repositories"
	"conectainova/internal/validator"
	"gorm.io/gorm"
)

type INecessidadeService interface {
	CreateNecessidade(request *necessidade.CreateNecessidadeRequest) (*necessidade.NecessidadeResponse, error)
	GetNecessidade(id uint) (*necessidade.NecessidadeResponse, error)
	GetNecessidades(page, size uint) ([]*necessidade.NecessidadeListResponse, error)
	UpdateNecessidade(id uint, request *necessidade.UpdateNecessidadeRequest) (*necessidade.NecessidadeResponse, error)
}

type NecessidadeService struct {
	NecessidadeRepo repositories.INecessidadeRepository
	LojistaRepo     repositories.ILojistaRepository
}

func NewNecessidadeService(necessidadeRepo repositories.INecessidadeRepository, lojistaRepo repositories.ILojistaRepository) INecessidadeService {
	return &NecessidadeService{NecessidadeRepo: necessidadeRepo, LojistaRepo: lojistaRepo}
}

func (s *NecessidadeService) CreateNecessidade(request *necessidade.CreateNecessidadeRequest) (*necessidade.NecessidadeResponse, error) {
	if err := validator.ValidateNecessidadeCreate(request); err != nil {
		return nil, err
	}
	entity := factory.NewNecessidade(request)
	if err := validator.ValidateNecessidade(entity); err != nil {
		return nil, err
	}
	if err := s.ensureLojista(entity.LojistaID); err != nil {
		return nil, err
	}
	if err := s.NecessidadeRepo.Create(entity); err != nil {
		return nil, err
	}
	return mapper.NecessidadeToResponse(entity), nil
}

func (s *NecessidadeService) GetNecessidade(id uint) (*necessidade.NecessidadeResponse, error) {
	entity, err := s.NecessidadeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return mapper.NecessidadeToResponse(entity), nil
}

func (s *NecessidadeService) GetNecessidades(page, size uint) ([]*necessidade.NecessidadeListResponse, error) {
	entities, err := s.NecessidadeRepo.GetAllWithLojista(page, size)
	if err != nil {
		return nil, err
	}
	return mapper.NecessidadesToListResponse(entities), nil
}

func (s *NecessidadeService) UpdateNecessidade(id uint, request *necessidade.UpdateNecessidadeRequest) (*necessidade.NecessidadeResponse, error) {
	entity, err := s.NecessidadeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if err := validator.ValidateNecessidadeUpdate(request); err != nil {
		return nil, err
	}
	updated := factory.UpdateNecessidade(entity, request)
	if err := validator.ValidateNecessidade(updated); err != nil {
		return nil, err
	}
	if err := s.ensureLojista(updated.LojistaID); err != nil {
		return nil, err
	}
	if err := s.NecessidadeRepo.Update(updated); err != nil {
		return nil, err
	}
	return mapper.NecessidadeToResponse(updated), nil
}

func (s *NecessidadeService) ensureLojista(id uint) error {
	_, err := s.LojistaRepo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.NewErrorMessage("lojista_id", "Lojista não encontrado.")
	}
	return err
}
