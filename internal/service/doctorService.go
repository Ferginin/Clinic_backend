package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"errors"
)

type DoctorServiceInterface interface {
	CreateDoctor(ctx context.Context, req *entity.DoctorCreateRequest) (*entity.Doctor, error)
	GetAllDoctors(ctx context.Context) ([]entity.Doctor, error)
	GetDoctorByID(ctx context.Context, id int) (*entity.Doctor, error)
	GetDoctorByUserID(ctx context.Context, userID int) (*entity.Doctor, error)
	GetDoctorsBySpecialization(ctx context.Context, specID int) ([]entity.Doctor, error)
	UpdateDoctor(ctx context.Context, id int, req *entity.DoctorUpdateRequest) (*entity.Doctor, error)
	DeleteDoctor(ctx context.Context, id int) error
	GetDoctorSchedule(ctx context.Context, doctorID int) ([]entity.Schedule, error)
	GetMySchedule(ctx context.Context, userID int) ([]entity.Schedule, error)
}

type DoctorService struct {
	doctorRepo   repository.DoctorRepositoryInterface
	specRepo     repository.SpecializationRepositoryInterface
	scheduleRepo repository.ScheduleRepositoryInterface
}

func NewDoctorService(doctorRepo repository.DoctorRepositoryInterface, specRepo repository.SpecializationRepositoryInterface, scheduleRepo repository.ScheduleRepositoryInterface) DoctorServiceInterface {
	return &DoctorService{
		doctorRepo:   doctorRepo,
		specRepo:     specRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (s *DoctorService) loadDoctorRelations(ctx context.Context, doctor *entity.Doctor) {
	specializations, _ := s.doctorRepo.GetSpecializations(ctx, doctor.ID)
	doctor.Specializations = specializations

	schedules, _ := s.scheduleRepo.GetByDoctorID(ctx, doctor.ID)
	doctor.Schedules = schedules
}

func (s *DoctorService) CreateDoctor(ctx context.Context, req *entity.DoctorCreateRequest) (*entity.Doctor, error) {
	if req.ScheduleID != nil {
		_, err := s.scheduleRepo.GetByID(ctx, *req.ScheduleID)
		if err != nil {
			return nil, errors.New("invalid schedule_id")
		}
	}

	doctor := &entity.Doctor{
		Fullname:    req.Fullname,
		Description: req.Description,
		DoctorPhoto: req.DoctorPhoto,
		UserID:      req.UserID,
		ScheduleID:  req.ScheduleID,
	}

	created, err := s.doctorRepo.Create(ctx, doctor)
	if err != nil {
		return nil, err
	}

	for _, specID := range req.SpecializationIDs {
		if err := s.doctorRepo.AddSpecialization(ctx, created.ID, specID); err != nil {
			return nil, err
		}
	}

	s.loadDoctorRelations(ctx, created)
	return created, nil
}

func (s *DoctorService) GetAllDoctors(ctx context.Context) ([]entity.Doctor, error) {
	doctors, err := s.doctorRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for i := range doctors {
		s.loadDoctorRelations(ctx, &doctors[i])
	}

	return doctors, nil
}

func (s *DoctorService) GetDoctorByID(ctx context.Context, id int) (*entity.Doctor, error) {
	doctor, err := s.doctorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.loadDoctorRelations(ctx, doctor)
	return doctor, nil
}

func (s *DoctorService) GetDoctorByUserID(ctx context.Context, userID int) (*entity.Doctor, error) {
	doctor, err := s.doctorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.loadDoctorRelations(ctx, doctor)
	return doctor, nil
}

func (s *DoctorService) GetDoctorsBySpecialization(ctx context.Context, specID int) ([]entity.Doctor, error) {
	doctors, err := s.doctorRepo.GetBySpecialization(ctx, specID)
	if err != nil {
		return nil, err
	}

	for i := range doctors {
		s.loadDoctorRelations(ctx, &doctors[i])
	}

	return doctors, nil
}

func (s *DoctorService) UpdateDoctor(ctx context.Context, id int, req *entity.DoctorUpdateRequest) (*entity.Doctor, error) {
	existing, err := s.doctorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

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
	if req.UserID != nil {
		existing.UserID = req.UserID
	}

	_, err = s.doctorRepo.Update(ctx, id, existing)
	if err != nil {
		return nil, err
	}

	if len(req.SpecializationIDs) > 0 {
		currentSpecs, _ := s.doctorRepo.GetSpecializations(ctx, id)
		for _, spec := range currentSpecs {
			_ = s.doctorRepo.RemoveSpecialization(ctx, id, spec.ID)
		}
		for _, specID := range req.SpecializationIDs {
			_ = s.doctorRepo.AddSpecialization(ctx, id, specID)
		}
	}

	return s.GetDoctorByID(ctx, id)
}

func (s *DoctorService) DeleteDoctor(ctx context.Context, id int) error {
	return s.doctorRepo.Delete(ctx, id)
}

func (s *DoctorService) GetDoctorSchedule(ctx context.Context, doctorID int) ([]entity.Schedule, error) {
	if _, err := s.doctorRepo.GetByID(ctx, doctorID); err != nil {
		return nil, err
	}

	schedules, err := s.scheduleRepo.GetByDoctorID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, errors.New("doctor has no schedule")
	}

	return schedules, nil
}

func (s *DoctorService) GetMySchedule(ctx context.Context, userID int) ([]entity.Schedule, error) {
	doctor, err := s.doctorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("doctor profile not found")
	}

	schedules, err := s.scheduleRepo.GetByDoctorID(ctx, doctor.ID)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, errors.New("no schedule assigned")
	}

	return schedules, nil
}
