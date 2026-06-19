package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CallbackRepositoryInterface interface {
	Create(ctx context.Context, req *entity.CallbackRequest) (*entity.CallbackRequest, error)
	GetAll(ctx context.Context) ([]entity.CallbackRequest, error)
	GetAllPaginated(ctx context.Context, limit, offset int) ([]entity.CallbackRequest, int, error)
	GetByID(ctx context.Context, id int) (*entity.CallbackRequest, error)
	Update(ctx context.Context, id int, req *entity.CallbackRequest) (*entity.CallbackRequest, error)
	Delete(ctx context.Context, id int) error
}

type CallbackRepository struct {
	db *pgxpool.Pool
}

func NewCallbackRepository(db *pgxpool.Pool) CallbackRepositoryInterface {
	return &CallbackRepository{db: db}
}

func (r *CallbackRepository) Create(ctx context.Context, req *entity.CallbackRequest) (*entity.CallbackRequest, error) {
	query := `
		INSERT INTO callback_requests (name, phone, message, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, phone, message, status, created_at, updated_at
	`

	var created entity.CallbackRequest
	err := r.db.QueryRow(ctx, query,
		req.Name,
		req.Phone,
		req.Message,
		req.Status,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Phone,
		&created.Message,
		&created.Status,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create callback request: %w", err)
	}

	return &created, nil
}

func (r *CallbackRepository) GetAll(ctx context.Context) ([]entity.CallbackRequest, error) {
	query := `
		SELECT id, name, phone, message, status, created_at, updated_at
		FROM callback_requests
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query callback requests: %w", err)
	}
	defer rows.Close()

	var requests []entity.CallbackRequest
	for rows.Next() {
		var req entity.CallbackRequest
		if err := rows.Scan(
			&req.ID,
			&req.Name,
			&req.Phone,
			&req.Message,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan callback request: %w", err)
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *CallbackRepository) GetAllPaginated(ctx context.Context, limit, offset int) ([]entity.CallbackRequest, int, error) {
	countQuery := `SELECT COUNT(*) FROM callback_requests`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count callback requests: %w", err)
	}

	query := `
		SELECT id, name, phone, message, status, created_at, updated_at
		FROM callback_requests
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query callback requests: %w", err)
	}
	defer rows.Close()

	var requests []entity.CallbackRequest
	for rows.Next() {
		var req entity.CallbackRequest
		if err := rows.Scan(
			&req.ID, &req.Name, &req.Phone, &req.Message,
			&req.Status, &req.CreatedAt, &req.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan callback request: %w", err)
		}
		requests = append(requests, req)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return requests, total, nil
}

func (r *CallbackRepository) GetByID(ctx context.Context, id int) (*entity.CallbackRequest, error) {
	query := `
		SELECT id, name, phone, message, status, created_at, updated_at
		FROM callback_requests WHERE id = $1
	`

	var req entity.CallbackRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&req.ID,
		&req.Name,
		&req.Phone,
		&req.Message,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("callback request not found")
		}
		return nil, fmt.Errorf("failed to get callback request: %w", err)
	}

	return &req, nil
}

func (r *CallbackRepository) Update(ctx context.Context, id int, req *entity.CallbackRequest) (*entity.CallbackRequest, error) {
	query := `
		UPDATE callback_requests
		SET name = $1, phone = $2, message = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, name, phone, message, status, created_at, updated_at
	`

	var updated entity.CallbackRequest
	err := r.db.QueryRow(ctx, query,
		req.Name,
		req.Phone,
		req.Message,
		req.Status,
		id,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Phone,
		&updated.Message,
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update callback request: %w", err)
	}

	return &updated, nil
}

func (r *CallbackRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM callback_requests WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete callback request: %w", err)
	}
	return nil
}
