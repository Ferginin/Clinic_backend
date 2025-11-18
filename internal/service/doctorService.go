package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DoctorService struct {
	doctorRepo *repository.DoctorRepository
	specRepo   *repository.SpecializationRepository
	schedRepo  *repository.ScheduleRepository
}

func NewDoctorService(db *pgxpool.Pool) *DoctorService {
	return &DoctorService{
		doctorRepo: repository.NewDoctorRepository(db),
		specRepo:   repository.NewSpecializationRepository(db),
		schedRepo:  repository.NewScheduleRepository(db),
	}
}

func (s *DoctorService) CreateDoctor(ctx context.Context, req *entity.DoctorCreateRequest) (*entity.Doctor, error) {
	// Валидация schedule_id если указан
	if req.ScheduleID != nil {
		_, err := s.schedRepo.GetByID(ctx, *req.ScheduleID)
		if err != nil {
			return nil, errors.New("invalid schedule_id")
		}
	}

	doctor := &entity.Doctor{
		Fullname:    req.Fullname,
		Description: req.Description,
		DoctorPhoto: req.DoctorPhoto,
		ScheduleID:  req.ScheduleID,
	}

	created, err := s.doctorRepo.Create(ctx, doctor)
	if err != nil {
		return nil, err
	}

	// Добавляем специализации
	for _, specID := range req.SpecializationIDs {
		if err := s.doctorRepo.AddSpecialization(ctx, created.ID, specID); err != nil {
			return nil, err
		}
	}

	// Загружаем специализации
	specializations, err := s.doctorRepo.GetSpecializations(ctx, created.ID)
	if err == nil {
		created.Specializations = specializations
	}

	// Загружаем расписание если есть
	if created.ScheduleID != nil {
		schedule, err := s.schedRepo.GetByID(ctx, *created.ScheduleID)
		if err == nil {
			created.Schedule = schedule
		}
	}

	return created, nil
}

func (s *DoctorService) GetAllDoctors(ctx context.Context) ([]entity.Doctor, error) {
	doctors, err := s.doctorRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Загружаем специализации для каждого врача
	for i := range doctors {
		specializations, _ := s.doctorRepo.GetSpecializations(ctx, doctors[i].ID)
		doctors[i].Specializations = specializations

		// Загружаем расписание
		if doctors[i].ScheduleID != nil {
			schedule, err := s.schedRepo.GetByID(ctx, *doctors[i].ScheduleID)
			if err == nil {
				doctors[i].Schedule = schedule
			}
		}
	}

	return doctors, nil
}

func (s *DoctorService) GetDoctorByID(ctx context.Context, id int) (*entity.Doctor, error) {
	doctor, err := s.doctorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Загружаем специализации
	specializations, _ := s.doctorRepo.GetSpecializations(ctx, doctor.ID)
	doctor.Specializations = specializations

	// Загружаем расписание
	if doctor.ScheduleID != nil {
		schedule, err := s.schedRepo.GetByID(ctx, *doctor.ScheduleID)
		if err == nil {
			doctor.Schedule = schedule
		}
	}

	return doctor, nil
}

func (s *DoctorService) GetDoctorsBySpecialization(ctx context.Context, specID int) ([]entity.Doctor, error) {
	doctors, err := s.doctorRepo.GetBySpecialization(ctx, specID)
	if err != nil {
		return nil, err
	}

	// Загружаем специализации для каждого врача
	for i := range doctors {
		specializations, _ := s.doctorRepo.GetSpecializations(ctx, doctors[i].ID)
		doctors[i].Specializations = specializations
	}

	return doctors, nil
}

func (s *DoctorService) UpdateDoctor(ctx context.Context, id int, req *entity.DoctorUpdateRequest) (*entity.Doctor, error) {
	existing, err := s.doctorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Обновляем только переданные поля
	if req.Fullname != nil {
		existing.Fullname = *req.Fullname
	}
	if req.Description != nil {
		existing.Description = req.Description
	}
	if req.DoctorPhoto != nil {
		existing.DoctorPhoto = req.DoctorPhoto
	}
	if req.ScheduleID != nil {
		existing.ScheduleID = req.ScheduleID
	}

	_, err = s.doctorRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	// Обновляем специализации если переданы
	if len(req.SpecializationIDs) > 0 {
		// Получаем текущие специализации
		currentSpecs, _ := s.doctorRepo.GetSpecializations(ctx, id)

		// Удаляем старые
		for _, spec := range currentSpecs {
			s.doctorRepo.RemoveSpecialization(ctx, id, spec.ID)
		}

		// Добавляем новые
		for _, specID := range req.SpecializationIDs {
			s.doctorRepo.AddSpecialization(ctx, id, specID)
		}
	}

	// Загружаем обновленные данные
	return s.GetDoctorByID(ctx, id)
}

func (s *DoctorService) DeleteDoctor(ctx context.Context, id int) error {
	return s.doctorRepo.Delete(ctx, id)
}

func (s *DoctorService) GetDoctorSchedule(ctx context.Context, doctorID int) (*entity.Schedule, error) {
	doctor, err := s.doctorRepo.GetByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	if doctor.ScheduleID == nil {
		return nil, errors.New("doctor has no schedule")
	}

	return s.schedRepo.GetByID(ctx, *doctor.ScheduleID)
}
