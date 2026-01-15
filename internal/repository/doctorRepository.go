package repository

import (
	"Clinic_backend/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DoctorRepositoryInterface interface {
	Create(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error)
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	GetByID(ctx context.Context, id int) (*entity.Doctor, error)
	GetBySpecialization(ctx context.Context, specializationID int) ([]entity.Doctor, error)
	Update(ctx context.Context, id int, doctor *entity.Doctor) (*entity.Doctor, error)
	Delete(ctx context.Context, id int) error
	AddSpecialization(ctx context.Context, doctorID, specializationID int) error
	RemoveSpecialization(ctx context.Context, doctorID, specializationID int) error
	GetSpecializations(ctx context.Context, doctorID int) ([]entity.Specialization, error)
}

type DoctorRepository struct {
	db *pgxpool.Pool
}

func NewDoctorRepository(db *pgxpool.Pool) DoctorRepositoryInterface {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) Create(ctx context.Context, doctor *entity.Doctor) (*entity.Doctor, error) {
	query := `
		INSERT INTO doctors (fullname, description, doctor_photo, schedule_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, fullname, description, doctor_photo, schedule_id, created_at, updated_at
	`

	var created entity.Doctor
	err := r.db.QueryRow(ctx, query,
		doctor.Fullname,
		doctor.Description,
		doctor.DoctorPhoto,
		doctor.ScheduleID,
	).Scan(
		&created.ID,
		&created.Fullname,
		&created.Description,
		&created.DoctorPhoto,
		&created.ScheduleID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create doctor: %w", err)
	}

	return &created, nil
}

func (r *DoctorRepository) GetAll(ctx context.Context) ([]entity.Doctor, error) {
	query := `
		SELECT id, fullname, description, doctor_photo, schedule_id, created_at, updated_at
		FROM doctors
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query doctors: %w", err)
	}
	defer rows.Close()

	var doctors []entity.Doctor
	for rows.Next() {
		var doctor entity.Doctor
		err := rows.Scan(
			&doctor.ID,
			&doctor.Fullname,
			&doctor.Description,
			&doctor.DoctorPhoto,
			&doctor.ScheduleID,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan doctor: %w", err)
		}
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}

func (r *DoctorRepository) GetByID(ctx context.Context, id int) (*entity.Doctor, error) {
	query := `
		SELECT id, fullname, description, doctor_photo, schedule_id, created_at, updated_at
		FROM doctors
		WHERE id = $1
	`

	var doctor entity.Doctor
	err := r.db.QueryRow(ctx, query, id).Scan(
		&doctor.ID,
		&doctor.Fullname,
		&doctor.Description,
		&doctor.DoctorPhoto,
		&doctor.ScheduleID,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("doctor not found")
		}
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	return &doctor, nil
}

func (r *DoctorRepository) GetBySpecialization(ctx context.Context, specializationID int) ([]entity.Doctor, error) {
	query := `
		SELECT d.id, d.fullname, d.description, d.doctor_photo, d.schedule_id, d.created_at, d.updated_at
		FROM doctors d
		INNER JOIN doctor_specializations ds ON d.id = ds.doctor_id
		WHERE ds.specialization_id = $1
		ORDER BY d.id
	`

	rows, err := r.db.Query(ctx, query, specializationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query doctors by specialization: %w", err)
	}
	defer rows.Close()

	var doctors []entity.Doctor
	for rows.Next() {
		var doctor entity.Doctor
		err := rows.Scan(
			&doctor.ID,
			&doctor.Fullname,
			&doctor.Description,
			&doctor.DoctorPhoto,
			&doctor.ScheduleID,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan doctor: %w", err)
		}
		doctors = append(doctors, doctor)
	}

	return doctors, nil
}

func (r *DoctorRepository) Update(ctx context.Context, id int, doctor *entity.Doctor) (*entity.Doctor, error) {
	query := `
		UPDATE doctors
		SET fullname = $1, description = $2, doctor_photo = $3, schedule_id = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING id, fullname, description, doctor_photo, schedule_id, created_at, updated_at
	`

	var updated entity.Doctor
	err := r.db.QueryRow(ctx, query,
		doctor.Fullname,
		doctor.Description,
		doctor.DoctorPhoto,
		doctor.ScheduleID,
		id,
	).Scan(
		&updated.ID,
		&updated.Fullname,
		&updated.Description,
		&updated.DoctorPhoto,
		&updated.ScheduleID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update doctor: %w", err)
	}

	return &updated, nil
}

func (r *DoctorRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM doctors WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete doctor: %w", err)
	}
	return nil
}

func (r *DoctorRepository) AddSpecialization(ctx context.Context, doctorID, specializationID int) error {
	query := `
		INSERT INTO doctor_specializations (doctor_id, specialization_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, doctorID, specializationID)
	if err != nil {
		return fmt.Errorf("failed to add specialization: %w", err)
	}
	return nil
}

func (r *DoctorRepository) RemoveSpecialization(ctx context.Context, doctorID, specializationID int) error {
	query := `
		DELETE FROM doctor_specializations
		WHERE doctor_id = $1 AND specialization_id = $2
	`
	_, err := r.db.Exec(ctx, query, doctorID, specializationID)
	if err != nil {
		return fmt.Errorf("failed to remove specialization: %w", err)
	}
	return nil
}

func (r *DoctorRepository) GetSpecializations(ctx context.Context, doctorID int) ([]entity.Specialization, error) {
	query := `
		SELECT s.id, s.name, s.created_at, s.updated_at
		FROM specializations s
		INNER JOIN doctor_specializations ds ON s.id = ds.specialization_id
		WHERE ds.doctor_id = $1
		ORDER BY s.id
	`

	rows, err := r.db.Query(ctx, query, doctorID)
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
