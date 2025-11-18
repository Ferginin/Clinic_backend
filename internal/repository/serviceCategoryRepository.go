package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceCategoryRepository struct {
	db *pgxpool.Pool
}

func NewServiceCategoryRepository(db *pgxpool.Pool) *ServiceCategoryRepository {
	return &ServiceCategoryRepository{db: db}
}

func (r *ServiceCategoryRepository) Create(ctx context.Context, category *entity.ServiceCategory) (*entity.ServiceCategory, error) {
	query := `
		INSERT INTO service_categories (name, description, category_photo, favorite, specialization_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, category_photo, favorite, specialization_id, created_at, updated_at
	`

	var created entity.ServiceCategory
	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.CategoryPhoto,
		category.Favorite,
		category.SpecializationID,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.CategoryPhoto,
		&created.Favorite,
		&created.SpecializationID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create service category: %w", err)
	}

	return &created, nil
}

func (r *ServiceCategoryRepository) GetAll(ctx context.Context) ([]entity.ServiceCategory, error) {
	query := `
		SELECT id, name, description, category_photo, favorite, specialization_id, created_at, updated_at
		FROM service_categories
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query service categories: %w", err)
	}
	defer rows.Close()

	var categories []entity.ServiceCategory
	for rows.Next() {
		var cat entity.ServiceCategory
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&cat.CategoryPhoto,
			&cat.Favorite,
			&cat.SpecializationID,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *ServiceCategoryRepository) GetByID(ctx context.Context, id int) (*entity.ServiceCategory, error) {
	query := `
		SELECT id, name, description, category_photo, favorite, specialization_id, created_at, updated_at
		FROM service_categories
		WHERE id = $1
	`

	var cat entity.ServiceCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cat.ID,
		&cat.Name,
		&cat.Description,
		&cat.CategoryPhoto,
		&cat.Favorite,
		&cat.SpecializationID,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("service category not found")
		}
		return nil, fmt.Errorf("failed to get service category: %w", err)
	}

	return &cat, nil
}

func (r *ServiceCategoryRepository) GetFavorites(ctx context.Context) ([]entity.ServiceCategory, error) {
	query := `
		SELECT id, name, description, category_photo, favorite, specialization_id, created_at, updated_at
		FROM service_categories
		WHERE favorite = true
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query favorite categories: %w", err)
	}
	defer rows.Close()

	var categories []entity.ServiceCategory
	for rows.Next() {
		var cat entity.ServiceCategory
		err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Description,
			&cat.CategoryPhoto,
			&cat.Favorite,
			&cat.SpecializationID,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service category: %w", err)
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *ServiceCategoryRepository) Update(ctx context.Context, id int, category *entity.ServiceCategory) (*entity.ServiceCategory, error) {
	query := `
		UPDATE service_categories
		SET name = $1, description = $2, category_photo = $3, favorite = $4, 
		    specialization_id = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, name, description, category_photo, favorite, specialization_id, created_at, updated_at
	`

	var updated entity.ServiceCategory
	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.CategoryPhoto,
		category.Favorite,
		category.SpecializationID,
		id,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.CategoryPhoto,
		&updated.Favorite,
		&updated.SpecializationID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update service category: %w", err)
	}

	return &updated, nil
}

func (r *ServiceCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM service_categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service category: %w", err)
	}
	return nil
}

func (r *ServiceCategoryRepository) SetFavorite(ctx context.Context, id int, favorite bool) error {
	query := `
		UPDATE service_categories 
		SET favorite = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, favorite, id)
	if err != nil {
		return fmt.Errorf("failed to set favorite: %w", err)
	}
	return nil
}
