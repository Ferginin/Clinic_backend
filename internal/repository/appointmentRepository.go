package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepositoryInterface interface {
	Create(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error)
	GetByID(ctx context.Context, id int) (*entity.Appointment, error)
	GetByPatient(ctx context.Context, patientID int) ([]entity.Appointment, error)
	GetByDoctor(ctx context.Context, doctorID int) ([]entity.Appointment, error)
	GetByDoctorAndDateRange(ctx context.Context, doctorID int, from, to time.Time) ([]entity.Appointment, error)
	GetAll(ctx context.Context) ([]entity.Appointment, error)
	GetAllPaginated(ctx context.Context, limit, offset int) ([]entity.Appointment, int, error)
	Update(ctx context.Context, id int, appointment *entity.Appointment) (*entity.Appointment, error)
	Delete(ctx context.Context, id int) error
}

type AppointmentRepository struct {
	db *pgxpool.Pool
}

func NewAppointmentRepository(db *pgxpool.Pool) AppointmentRepositoryInterface {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(ctx context.Context, appointment *entity.Appointment) (*entity.Appointment, error) {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, scheduled_at, status, notes, results)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
	`

	var created entity.Appointment
	err := r.db.QueryRow(ctx, query,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.ScheduledAt,
		appointment.Status,
		appointment.Notes,
		appointment.Results,
	).Scan(
		&created.ID,
		&created.PatientID,
		&created.DoctorID,
		&created.ScheduledAt,
		&created.Status,
		&created.Notes,
		&created.Results,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return &created, nil
}

func (r *AppointmentRepository) scanAppointment(row pgx.Row) (*entity.Appointment, error) {
	var appointment entity.Appointment
	err := row.Scan(
		&appointment.ID,
		&appointment.PatientID,
		&appointment.DoctorID,
		&appointment.ScheduledAt,
		&appointment.Status,
		&appointment.Notes,
		&appointment.Results,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *AppointmentRepository) GetByID(ctx context.Context, id int) (*entity.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
		FROM appointments WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	appointment, err := r.scanAppointment(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	return appointment, nil
}

func (r *AppointmentRepository) GetByPatient(ctx context.Context, patientID int) ([]entity.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
		FROM appointments WHERE patient_id = $1 ORDER BY scheduled_at DESC
	`

	rows, err := r.db.Query(ctx, query, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments by patient: %w", err)
	}
	defer rows.Close()

	var appointments []entity.Appointment
	for rows.Next() {
		var appointment entity.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.ScheduledAt,
			&appointment.Status,
			&appointment.Notes,
			&appointment.Results,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByDoctor(ctx context.Context, doctorID int) ([]entity.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
		FROM appointments WHERE doctor_id = $1 ORDER BY scheduled_at DESC
	`

	rows, err := r.db.Query(ctx, query, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments by doctor: %w", err)
	}
	defer rows.Close()

	var appointments []entity.Appointment
	for rows.Next() {
		var appointment entity.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.ScheduledAt,
			&appointment.Status,
			&appointment.Notes,
			&appointment.Results,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetAll(ctx context.Context) ([]entity.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
		FROM appointments ORDER BY scheduled_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments: %w", err)
	}
	defer rows.Close()

	var appointments []entity.Appointment
	for rows.Next() {
		var appointment entity.Appointment
		if err := rows.Scan(
			&appointment.ID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.ScheduledAt,
			&appointment.Status,
			&appointment.Notes,
			&appointment.Results,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetAllPaginated(ctx context.Context, limit, offset int) ([]entity.Appointment, int, error) {
	countQuery := `SELECT COUNT(*) FROM appointments`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
		FROM appointments
		ORDER BY scheduled_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query appointments: %w", err)
	}
	defer rows.Close()

	var appointments []entity.Appointment
	for rows.Next() {
		var a entity.Appointment
		if err := rows.Scan(
			&a.ID, &a.PatientID, &a.DoctorID, &a.ScheduledAt,
			&a.Status, &a.Notes, &a.Results, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}
	return appointments, total, nil
}

func (r *AppointmentRepository) Update(ctx context.Context, id int, appointment *entity.Appointment) (*entity.Appointment, error) {
	query := `
		UPDATE appointments
		SET patient_id = $1, doctor_id = $2, scheduled_at = $3, status = $4, notes = $5, results = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
	`

	var updated entity.Appointment
	err := r.db.QueryRow(ctx, query,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.ScheduledAt,
		appointment.Status,
		appointment.Notes,
		appointment.Results,
		id,
	).Scan(
		&updated.ID,
		&updated.PatientID,
		&updated.DoctorID,
		&updated.ScheduledAt,
		&updated.Status,
		&updated.Notes,
		&updated.Results,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return &updated, nil
}

func (r *AppointmentRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM appointments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}
	return nil
}

func (r *AppointmentRepository) GetByDoctorAndDateRange(ctx context.Context, doctorID int, from, to time.Time) ([]entity.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, status, notes, results, created_at, updated_at
		FROM appointments
		WHERE doctor_id = $1 AND scheduled_at >= $2 AND scheduled_at < $3
		  AND status NOT IN ('canceled')
		ORDER BY scheduled_at
	`

	rows, err := r.db.Query(ctx, query, doctorID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments by date range: %w", err)
	}
	defer rows.Close()

	var appointments []entity.Appointment
	for rows.Next() {
		var a entity.Appointment
		if err := rows.Scan(
			&a.ID, &a.PatientID, &a.DoctorID, &a.ScheduledAt, &a.Status,
			&a.Notes, &a.Results, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}
