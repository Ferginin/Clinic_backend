package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SpecializationService struct {
	specRepo *repository.SpecializationRepository
}

func NewSpecializationService(db *pgxpool.Pool) *SpecializationService {
	return &SpecializationService{
		specRepo: repository.NewSpecializationRepository(db),
	}
}

func (s *SpecializationService) CreateSpecialization(ctx context.Context, spec *entity.Specialization) (*entity.Specialization, error) {
	return s.specRepo.Create(ctx, spec)
}

func (s *SpecializationService) GetAllSpecializations(ctx context.Context) ([]entity.Specialization, error) {
	return s.specRepo.GetAll(ctx)
}

func (s *SpecializationService) GetSpecializationByID(ctx context.Context, id int) (*entity.Specialization, error) {
	return s.specRepo.GetByID(ctx, id)
}

func (s *SpecializationService) UpdateSpecialization(ctx context.Context, id int, spec *entity.Specialization) (*entity.Specialization, error) {
	return s.specRepo.Update(ctx, id, spec)
}

func (s *SpecializationService) DeleteSpecialization(ctx context.Context, id int) error {
	return s.specRepo.Delete(ctx, id)
}
