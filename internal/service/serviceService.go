package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceService struct {
	serviceRepo  *repository.ServiceRepository
	categoryRepo *repository.ServiceCategoryRepository
	specRepo     *repository.SpecializationRepository
}

func NewServiceService(db *pgxpool.Pool) *ServiceService {
	return &ServiceService{
		serviceRepo:  repository.NewServiceRepository(db),
		categoryRepo: repository.NewServiceCategoryRepository(db),
		specRepo:     repository.NewSpecializationRepository(db),
	}
}

func (s *ServiceService) CreateService(ctx context.Context, req *entity.ServiceCreateRequest) (*entity.Service, error) {
	// Валидация category_id если указан
	if req.ServiceCategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *req.ServiceCategoryID)
		if err != nil {
			return nil, errors.New("invalid service_category_id")
		}
	}

	// Валидация specialization_id если указан
	if req.SpecializationID != nil {
		_, err := s.specRepo.GetByID(ctx, *req.SpecializationID)
		if err != nil {
			return nil, errors.New("invalid specialization_id")
		}
	}

	// Валидация цены
	if req.Price != nil && *req.Price < 0 {
		return nil, errors.New("price must be positive")
	}

	service := &entity.Service{
		Name:              req.Name,
		Description:       req.Description,
		SpecificPhoto:     req.SpecificPhoto,
		Price:             req.Price,
		ServiceCategoryID: req.ServiceCategoryID,
		SpecializationID:  req.SpecializationID,
	}

	return s.serviceRepo.Create(ctx, service)
}

func (s *ServiceService) GetAllServices(ctx context.Context) ([]entity.Service, error) {
	return s.serviceRepo.GetAll(ctx)
}

func (s *ServiceService) GetServiceByID(ctx context.Context, id int) (*entity.Service, error) {
	return s.serviceRepo.GetByID(ctx, id)
}

func (s *ServiceService) GetServicesByCategory(ctx context.Context, categoryID int) ([]entity.Service, error) {
	// Проверяем существование категории
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	return s.serviceRepo.GetByCategory(ctx, categoryID)
}

func (s *ServiceService) GetServicesBySpecialization(ctx context.Context, specID int) ([]entity.Service, error) {
	// Проверяем существование специализации
	_, err := s.specRepo.GetByID(ctx, specID)
	if err != nil {
		return nil, errors.New("specialization not found")
	}

	return s.serviceRepo.GetBySpecialization(ctx, specID)
}

func (s *ServiceService) UpdateService(ctx context.Context, id int, req *entity.ServiceCreateRequest) (*entity.Service, error) {
	existing, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля
	existing.Name = req.Name
	existing.Description = req.Description
	existing.SpecificPhoto = req.SpecificPhoto
	existing.Price = req.Price
	existing.ServiceCategoryID = req.ServiceCategoryID
	existing.SpecializationID = req.SpecializationID

	return s.serviceRepo.Update(ctx, id, existing)
}

func (s *ServiceService) DeleteService(ctx context.Context, id int) error {
	return s.serviceRepo.Delete(ctx, id)
}
