package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryService struct {
	categoryRepo *repository.ServiceCategoryRepository
	specRepo     *repository.SpecializationRepository
}

func NewCategoryService(db *pgxpool.Pool) *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewServiceCategoryRepository(db),
		specRepo:     repository.NewSpecializationRepository(db),
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, category *entity.ServiceCategory) (*entity.ServiceCategory, error) {
	return s.categoryRepo.Create(ctx, category)
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]entity.ServiceCategory, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (*entity.ServiceCategory, error) {
	return s.categoryRepo.GetByID(ctx, id)
}

func (s *CategoryService) GetFavoriteCategories(ctx context.Context) ([]entity.ServiceCategory, error) {
	return s.categoryRepo.GetFavorites(ctx)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int, category *entity.ServiceCategory) (*entity.ServiceCategory, error) {
	return s.categoryRepo.Update(ctx, id, category)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.categoryRepo.Delete(ctx, id)
}

func (s *CategoryService) ToggleFavorite(ctx context.Context, id int) error {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.categoryRepo.SetFavorite(ctx, id, !category.Favorite)
}
