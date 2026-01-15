package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpecializationRepositoryInterface interface {
	Create(ctx context.Context, spec *entity.Specialization) (*entity.Specialization, error)
	GetAll(ctx context.Context) ([]entity.Specialization, error)
	GetByID(ctx context.Context, id int) (*entity.Specialization, error)
	Update(ctx context.Context, id int, spec *entity.Specialization) (*entity.Specialization, error)
	Delete(ctx context.Context, id int) error
}

type SpecializationRepository struct {
	db *pgxpool.Pool
}

func NewSpecializationRepository(db *pgxpool.Pool) SpecializationRepositoryInterface {
	return &SpecializationRepository{db: db}
}

func (r *SpecializationRepository) Create(ctx context.Context, spec *entity.Specialization) (*entity.Specialization, error) {
	query := `
		INSERT INTO specializations (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`

	var created entity.Specialization
	err := r.db.QueryRow(ctx, query, spec.Name).Scan(
		&created.ID, &created.Name, &created.CreatedAt, &created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create specialization: %w", err)
	}

	return &created, nil
}

func (r *SpecializationRepository) GetAll(ctx context.Context) ([]entity.Specialization, error) {
	query := `SELECT id, name, created_at, updated_at FROM specializations ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query specializations: %w", err)
	}
	defer rows.Close()

	var specializations []entity.Specialization
	for rows.Next() {
		var spec entity.Specialization
		err := rows.Scan(&spec.ID, &spec.Name, &spec.CreatedAt, &spec.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan specialization: %w", err)
		}
		specializations = append(specializations, spec)
	}

	return specializations, nil
}

func (r *SpecializationRepository) GetByID(ctx context.Context, id int) (*entity.Specialization, error) {
	query := `SELECT id, name, created_at, updated_at FROM specializations WHERE id = $1`

	var spec entity.Specialization
	err := r.db.QueryRow(ctx, query, id).Scan(
		&spec.ID, &spec.Name, &spec.CreatedAt, &spec.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("specialization not found")
		}
		return nil, fmt.Errorf("failed to get specialization: %w", err)
	}

	return &spec, nil
}

func (r *SpecializationRepository) Update(ctx context.Context, id int, spec *entity.Specialization) (*entity.Specialization, error) {
	query := `
		UPDATE specializations
		SET name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, name, created_at, updated_at
	`

	var updated entity.Specialization
	err := r.db.QueryRow(ctx, query, spec.Name, id).Scan(
		&updated.ID, &updated.Name, &updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update specialization: %w", err)
	}

	return &updated, nil
}

func (r *SpecializationRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM specializations WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete specialization: %w", err)
	}
	return nil
}
