package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LicenseRepository struct {
	db *pgxpool.Pool
}

func NewLicenseRepository(db *pgxpool.Pool) *LicenseRepository {
	return &LicenseRepository{db: db}
}

func (r *LicenseRepository) Create(ctx context.Context, license *entity.License) (*entity.License, error) {
	query := `
		INSERT INTO licenses (photo, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, photo, name, description, created_at, updated_at
	`

	var created entity.License
	err := r.db.QueryRow(ctx, query, license.Photo, license.Name, license.Description).Scan(
		&created.ID, &created.Photo, &created.Name, &created.Description,
		&created.CreatedAt, &created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create license: %w", err)
	}

	return &created, nil
}

func (r *LicenseRepository) GetAll(ctx context.Context) ([]entity.License, error) {
	query := `
		SELECT id, photo, name, description, created_at, updated_at 
		FROM licenses 
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query licenses: %w", err)
	}
	defer rows.Close()

	var licenses []entity.License
	for rows.Next() {
		var license entity.License
		err := rows.Scan(
			&license.ID, &license.Photo, &license.Name, &license.Description,
			&license.CreatedAt, &license.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan license: %w", err)
		}
		licenses = append(licenses, license)
	}

	return licenses, nil
}

func (r *LicenseRepository) GetByID(ctx context.Context, id int) (*entity.License, error) {
	query := `
		SELECT id, photo, name, description, created_at, updated_at 
		FROM licenses 
		WHERE id = $1
	`

	var license entity.License
	err := r.db.QueryRow(ctx, query, id).Scan(
		&license.ID, &license.Photo, &license.Name, &license.Description,
		&license.CreatedAt, &license.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("license not found")
		}
		return nil, fmt.Errorf("failed to get license: %w", err)
	}

	return &license, nil
}

func (r *LicenseRepository) Update(ctx context.Context, id int, license *entity.License) (*entity.License, error) {
	query := `
		UPDATE licenses
		SET photo = $1, name = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, photo, name, description, created_at, updated_at
	`

	var updated entity.License
	err := r.db.QueryRow(ctx, query,
		license.Photo, license.Name, license.Description, id,
	).Scan(
		&updated.ID, &updated.Photo, &updated.Name, &updated.Description,
		&updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update license: %w", err)
	}

	return &updated, nil
}

func (r *LicenseRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM licenses WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete license: %w", err)
	}
	return nil
}
