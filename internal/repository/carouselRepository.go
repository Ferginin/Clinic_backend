package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarouselRepositoryInterface interface {
	Create(ctx context.Context, course *entity.Carousel) (*entity.Carousel, error)
	GetAll(ctx context.Context) ([]entity.Carousel, error)
	GetByID(ctx context.Context, id int) (*entity.Carousel, error)
	Update(ctx context.Context, id int, carousel *entity.Carousel) (*entity.Carousel, error)
	Delete(ctx context.Context, id int) error
}

type CarouselRepository struct {
	db *pgxpool.Pool
}

func NewCarouselRepository(db *pgxpool.Pool) CarouselRepositoryInterface {
	return &CarouselRepository{db: db}
}

func (r *CarouselRepository) Create(ctx context.Context, carousel *entity.Carousel) (*entity.Carousel, error) {
	query := `
		INSERT INTO main_carusel (image, header, description)
		VALUES ($1, $2, $3)
		RETURNING id, image, header, description, created_at, updated_at
	`

	var created entity.Carousel
	err := r.db.QueryRow(ctx, query, carousel.Image, carousel.Header, carousel.Description).Scan(
		&created.ID, &created.Image, &created.Header, &created.Description,
		&created.CreatedAt, &created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create carousel: %w", err)
	}

	return &created, nil
}

func (r *CarouselRepository) GetAll(ctx context.Context) ([]entity.Carousel, error) {
	query := `
		SELECT id, image, header, description, created_at, updated_at 
		FROM main_carusel 
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query carousel: %w", err)
	}
	defer rows.Close()

	var carousels []entity.Carousel
	for rows.Next() {
		var carousel entity.Carousel
		err := rows.Scan(
			&carousel.ID, &carousel.Image, &carousel.Header, &carousel.Description,
			&carousel.CreatedAt, &carousel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan carousel: %w", err)
		}
		carousels = append(carousels, carousel)
	}

	return carousels, nil
}

func (r *CarouselRepository) GetByID(ctx context.Context, id int) (*entity.Carousel, error) {
	query := `
		SELECT id, image, header, description, created_at, updated_at 
		FROM main_carusel 
		WHERE id = $1
	`

	var carousel entity.Carousel
	err := r.db.QueryRow(ctx, query, id).Scan(
		&carousel.ID, &carousel.Image, &carousel.Header, &carousel.Description,
		&carousel.CreatedAt, &carousel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("carousel not found")
		}
		return nil, fmt.Errorf("failed to get carousel: %w", err)
	}

	return &carousel, nil
}

func (r *CarouselRepository) Update(ctx context.Context, id int, carousel *entity.Carousel) (*entity.Carousel, error) {
	query := `
		UPDATE main_carusel
		SET image = $1, header = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, image, header, description, created_at, updated_at
	`

	var updated entity.Carousel
	err := r.db.QueryRow(ctx, query,
		carousel.Image, carousel.Header, carousel.Description, id,
	).Scan(
		&updated.ID, &updated.Image, &updated.Header, &updated.Description,
		&updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update carousel: %w", err)
	}

	return &updated, nil
}

func (r *CarouselRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM main_carusel WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete carousel: %w", err)
	}
	return nil
}
