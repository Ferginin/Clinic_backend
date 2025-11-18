package service

import (
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"Clinic_backend/internal/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepository
}

func NewScheduleService(db *pgxpool.Pool) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: repository.NewScheduleRepository(db),
	}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, schedule *entity.Schedule) (*entity.Schedule, error) {
	// Валидация дня недели
	if err := utils.ValidateDayOfWeek(schedule.Day); err != nil {
		return nil, err
	}

	// Валидация временного слота
	if err := utils.ValidateTimeSlot(schedule.TimeFrom, schedule.TimeTo); err != nil {
		return nil, err
	}

	return s.scheduleRepo.Create(ctx, schedule)
}

func (s *ScheduleService) GetAllSchedules(ctx context.Context) ([]entity.Schedule, error) {
	return s.scheduleRepo.GetAll(ctx)
}

func (s *ScheduleService) GetScheduleByID(ctx context.Context, id int) (*entity.Schedule, error) {
	return s.scheduleRepo.GetByID(ctx, id)
}

func (s *ScheduleService) GetScheduleByDay(ctx context.Context, day int) ([]entity.Schedule, error) {
	if err := utils.ValidateDayOfWeek(day); err != nil {
		return nil, err
	}

	return s.scheduleRepo.GetByDay(ctx, day)
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id int, schedule *entity.Schedule) (*entity.Schedule, error) {
	// Валидация дня недели
	if err := utils.ValidateDayOfWeek(schedule.Day); err != nil {
		return nil, err
	}

	// Валидация временного слота
	if err := utils.ValidateTimeSlot(schedule.TimeFrom, schedule.TimeTo); err != nil {
		return nil, err
	}

	return s.scheduleRepo.Update(ctx, id, schedule)
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id int) error {
	return s.scheduleRepo.Delete(ctx, id)
}
