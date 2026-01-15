package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
)

type LicenseServiceInterface interface {
	CreateLicense(ctx context.Context, license *entity.License) (*entity.License, error)
	GetAllLicenses(ctx context.Context) ([]entity.License, error)
	GetLicenseByID(ctx context.Context, id int) (*entity.License, error)
	UpdateLicense(ctx context.Context, id int, license *entity.License) (*entity.License, error)
	DeleteLicense(ctx context.Context, id int) error
}

type LicenseService struct {
	licenseRepo repository.LicenseRepositoryInterface
}

func NewLicenseService(licenseRepo repository.LicenseRepositoryInterface) LicenseServiceInterface {
	return &LicenseService{
		licenseRepo: licenseRepo,
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
