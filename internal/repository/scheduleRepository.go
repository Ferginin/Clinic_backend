package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	query := `
		INSERT INTO schedules (day, time_from, time_to)
		VALUES ($1, $2, $3)
		RETURNING id, day, time_from, time_to, created_at, updated_at
	`

	var created entity.Schedule
	err := r.db.QueryRow(ctx, query, schedule.Day, schedule.TimeFrom, schedule.TimeTo).Scan(
		&created.ID, &created.Day, &created.TimeFrom, &created.TimeTo,
		&created.CreatedAt, &created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return &created, nil
}

func (r *ScheduleRepository) GetAll(ctx context.Context) ([]entity.Schedule, error) {
	query := `
		SELECT id, day, time_from, time_to, created_at, updated_at 
		FROM schedules 
		ORDER BY day, time_from
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query schedules: %w", err)
	}
	defer rows.Close()

	var schedules []entity.Schedule
	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(
			&schedule.ID, &schedule.Day, &schedule.TimeFrom, &schedule.TimeTo,
			&schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id int) (*entity.Schedule, error) {
	query := `
		SELECT id, day, time_from, time_to, created_at, updated_at 
		FROM schedules 
		WHERE id = $1
	`

	var schedule entity.Schedule
	err := r.db.QueryRow(ctx, query, id).Scan(
		&schedule.ID, &schedule.Day, &schedule.TimeFrom, &schedule.TimeTo,
		&schedule.CreatedAt, &schedule.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("schedule not found")
		}
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return &schedule, nil
}

func (r *ScheduleRepository) GetByDay(ctx context.Context, day int) ([]entity.Schedule, error) {
	query := `
		SELECT id, day, time_from, time_to, created_at, updated_at 
		FROM schedules 
		WHERE day = $1 
		ORDER BY time_from
	`

	rows, err := r.db.Query(ctx, query, day)
	if err != nil {
		return nil, fmt.Errorf("failed to query schedules by day: %w", err)
	}
	defer rows.Close()

	var schedules []entity.Schedule
	for rows.Next() {
		var schedule entity.Schedule
		err := rows.Scan(
			&schedule.ID, &schedule.Day, &schedule.TimeFrom, &schedule.TimeTo,
			&schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, id int, schedule *entity.Schedule) (*entity.Schedule, error) {
	query := `
		UPDATE schedules
		SET day = $1, time_from = $2, time_to = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, day, time_from, time_to, created_at, updated_at
	`

	var updated entity.Schedule
	err := r.db.QueryRow(ctx, query,
		schedule.Day, schedule.TimeFrom, schedule.TimeTo, id,
	).Scan(
		&updated.ID, &updated.Day, &updated.TimeFrom, &updated.TimeTo,
		&updated.CreatedAt, &updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return &updated, nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM schedules WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}
