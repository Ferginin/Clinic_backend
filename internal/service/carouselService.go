package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CarouselService struct {
	carouselRepo *repository.CarouselRepository
}

func NewCarouselService(db *pgxpool.Pool) *CarouselService {
	return &CarouselService{
		carouselRepo: repository.NewCarouselRepository(db),
	}
}

func (s *CarouselService) CreateSlide(ctx context.Context, carousel *entity.Carousel) (*entity.Carousel, error) {
	return s.carouselRepo.Create(ctx, carousel)
}

func (s *CarouselService) GetAllSlides(ctx context.Context) ([]entity.Carousel, error) {
	return s.carouselRepo.GetAll(ctx)
}

func (s *CarouselService) GetSlideByID(ctx context.Context, id int) (*entity.Carousel, error) {
	return s.carouselRepo.GetByID(ctx, id)
}

func (s *CarouselService) UpdateSlide(ctx context.Context, id int, carousel *entity.Carousel) (*entity.Carousel, error) {
	return s.carouselRepo.Update(ctx, id, carousel)
}

func (s *CarouselService) DeleteSlide(ctx context.Context, id int) error {
	return s.carouselRepo.Delete(ctx, id)
}
