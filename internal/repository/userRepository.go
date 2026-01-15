package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, id int, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int) error
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = 'user'))
		RETURNING id, username, email, confirmed, blocked, role_id, created_at, updated_at
	`

	var createdUser entity.User
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.Confirmed,
		&createdUser.Blocked,
		&createdUser.RoleID,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get role name
	if createdUser.RoleID != nil {
		roleQuery := `SELECT name FROM roles WHERE id = $1`
		err = r.db.QueryRow(ctx, roleQuery, *createdUser.RoleID).Scan(&createdUser.RoleName)
		if err != nil {
			createdUser.RoleName = "user"
		}
	}

	return &createdUser, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT u.username, u.email, u.password, u.confirmed, u.blocked, 
		    	COALESCE(r.name, 'user') as role_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.email = $1
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		//&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Confirmed,
		&user.Blocked,
		//&user.RoleID,
		&user.RoleName,
		//&user.CreatedAt,
		//&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.confirmed, u.blocked, 
		       u.role_id, COALESCE(r.name, 'user') as role_name, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`

	var user entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Confirmed,
		&user.Blocked,
		&user.RoleID,
		&user.RoleName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.confirmed, u.blocked, 
		       u.role_id, COALESCE(r.name, 'user') as role_name, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Confirmed,
			&user.Blocked,
			&user.RoleID,
			&user.RoleName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id int, user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users 
		SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, username, email, confirmed, blocked, role_id, created_at, updated_at
	`

	var updatedUser entity.User
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, id).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.Confirmed,
		&updatedUser.Blocked,
		&updatedUser.RoleID,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
