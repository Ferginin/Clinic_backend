package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]entity.Role, error) {
	query := `SELECT id, name, created_at, updated_at FROM roles ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepository) GetByID(ctx context.Context, id int) (*entity.Role, error) {
	query := `SELECT id, name, created_at, updated_at FROM roles WHERE id = $1`

	var role entity.Role
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

func (r *RoleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT id, name, created_at, updated_at FROM roles WHERE name = $1`

	var role entity.Role
	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}
