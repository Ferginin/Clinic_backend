package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LicenseService struct {
	licenseRepo *repository.LicenseRepository
}

func NewLicenseService(db *pgxpool.Pool) *LicenseService {
	return &LicenseService{
		licenseRepo: repository.NewLicenseRepository(db),
	}
}

func (s *LicenseService) CreateLicense(ctx context.Context, license *entity.License) (*entity.License, error) {
	return s.licenseRepo.Create(ctx, license)
}

func (s *LicenseService) GetAllLicenses(ctx context.Context) ([]entity.License, error) {
	return s.licenseRepo.GetAll(ctx)
}

func (s *LicenseService) GetLicenseByID(ctx context.Context, id int) (*entity.License, error) {
	return s.licenseRepo.GetByID(ctx, id)
}

func (s *LicenseService) UpdateLicense(ctx context.Context, id int, license *entity.License) (*entity.License, error) {
	return s.licenseRepo.Update(ctx, id, license)
}

func (s *LicenseService) DeleteLicense(ctx context.Context, id int) error {
	return s.licenseRepo.Delete(ctx, id)
}
