package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(ctx context.Context, service *entity.Service) (*entity.Service, error) {
	query := `
		INSERT INTO services (name, description, specific_photo, price, service_category_id, specialization_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, specific_photo, price, service_category_id, specialization_id, created_at, updated_at
	`

	var created entity.Service
	err := r.db.QueryRow(ctx, query,
		service.Name,
		service.Description,
		service.SpecificPhoto,
		service.Price,
		service.ServiceCategoryID,
		service.SpecializationID,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.SpecificPhoto,
		&created.Price,
		&created.ServiceCategoryID,
		&created.SpecializationID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return &created, nil
}

func (r *ServiceRepository) GetAll(ctx context.Context) ([]entity.Service, error) {
	query := `
		SELECT id, name, description, specific_photo, price, service_category_id, specialization_id, created_at, updated_at
		FROM services
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var service entity.Service
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.SpecificPhoto,
			&service.Price,
			&service.ServiceCategoryID,
			&service.SpecializationID,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) GetByID(ctx context.Context, id int) (*entity.Service, error) {
	query := `
		SELECT id, name, description, specific_photo, price, service_category_id, specialization_id, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	var service entity.Service
	err := r.db.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.SpecificPhoto,
		&service.Price,
		&service.ServiceCategoryID,
		&service.SpecializationID,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("service not found")
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return &service, nil
}

func (r *ServiceRepository) GetByCategory(ctx context.Context, categoryID int) ([]entity.Service, error) {
	query := `
		SELECT id, name, description, specific_photo, price, service_category_id, specialization_id, created_at, updated_at
		FROM services
		WHERE service_category_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query services by category: %w", err)
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var service entity.Service
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.SpecificPhoto,
			&service.Price,
			&service.ServiceCategoryID,
			&service.SpecializationID,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) GetBySpecialization(ctx context.Context, specializationID int) ([]entity.Service, error) {
	query := `
		SELECT id, name, description, specific_photo, price, service_category_id, specialization_id, created_at, updated_at
		FROM services
		WHERE specialization_id = $1
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, specializationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query services by specialization: %w", err)
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var service entity.Service
		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.SpecificPhoto,
			&service.Price,
			&service.ServiceCategoryID,
			&service.SpecializationID,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) Update(ctx context.Context, id int, service *entity.Service) (*entity.Service, error) {
	query := `
		UPDATE services
		SET name = $1, description = $2, specific_photo = $3, price = $4,
		    service_category_id = $5, specialization_id = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING id, name, description, specific_photo, price, service_category_id, specialization_id, created_at, updated_at
	`

	var updated entity.Service
	err := r.db.QueryRow(ctx, query,
		service.Name,
		service.Description,
		service.SpecificPhoto,
		service.Price,
		service.ServiceCategoryID,
		service.SpecializationID,
		id,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.SpecificPhoto,
		&updated.Price,
		&updated.ServiceCategoryID,
		&updated.SpecializationID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return &updated, nil
}

func (r *ServiceRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}
	return nil
}
