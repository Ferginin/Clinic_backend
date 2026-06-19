package service

import (
	"Clinic_backend/config"
	"Clinic_backend/internal/entity"
	"Clinic_backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"
)

type AppointmentServiceInterface interface {
	CreateAppointment(ctx context.Context, patientID int, req *entity.AppointmentCreateRequest) (*entity.Appointment, error)
	GetAppointmentByID(ctx context.Context, id int) (*entity.Appointment, error)
	GetAppointmentsByPatient(ctx context.Context, patientID int) ([]entity.Appointment, error)
	GetAppointmentsByDoctor(ctx context.Context, doctorID int) ([]entity.Appointment, error)
	GetAllAppointments(ctx context.Context) ([]entity.Appointment, error)
	GetAllAppointmentsPaginated(ctx context.Context, limit, offset int) ([]entity.Appointment, int, error)
	UpdateAppointment(ctx context.Context, appointmentID int, patientID int, role string, req *entity.AppointmentUpdateRequest) (*entity.Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID int, patientID int, role string) error
	AddAppointmentResult(ctx context.Context, appointmentID int, userID int, resultReq *entity.AppointmentResultRequest) (*entity.Appointment, error)
	GetAvailableSlots(ctx context.Context, doctorID int, date string) (*entity.AvailableSlotsResponse, error)
	AdminCreateAppointment(ctx context.Context, adminID int, req *entity.AppointmentAdminCreateRequest) (*entity.Appointment, error)
}

type AppointmentService struct {
	appointmentRepo repository.AppointmentRepositoryInterface
	userRepo        repository.UserRepositoryInterface
	doctorRepo      repository.DoctorRepositoryInterface
	scheduleRepo    repository.ScheduleRepositoryInterface
	emailService    EmailServiceInterface
	cfg             *config.Config
}

func NewAppointmentService(
	ar repository.AppointmentRepositoryInterface,
	ur repository.UserRepositoryInterface,
	dr repository.DoctorRepositoryInterface,
	sr repository.ScheduleRepositoryInterface,
	es EmailServiceInterface,
	cfg *config.Config,
) AppointmentServiceInterface {
	return &AppointmentService{
		appointmentRepo: ar,
		userRepo:        ur,
		doctorRepo:      dr,
		scheduleRepo:    sr,
		emailService:    es,
		cfg:             cfg,
	}
}

func (s *AppointmentService) GetAvailableSlots(ctx context.Context, doctorID int, date string) (*entity.AvailableSlotsResponse, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	if parsedDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("date must be today or in the future")
	}

	if _, err := s.doctorRepo.GetByID(ctx, doctorID); err != nil {
		return nil, errors.New("doctor not found")
	}

	isoDay := int(parsedDate.Weekday())
	if isoDay == 0 {
		isoDay = 7
	}

	schedule, err := s.scheduleRepo.GetByDoctorAndDay(ctx, doctorID, isoDay)
	if err != nil {
		return &entity.AvailableSlotsResponse{
			DoctorID: doctorID,
			Date:     date,
			Slots:    []entity.TimeSlot{},
		}, nil
	}

	timeFrom, err := parseTime(schedule.TimeFrom)
	if err != nil {
		return nil, fmt.Errorf("invalid time_from in schedule: %s", schedule.TimeFrom)
	}

	timeTo, err := parseTime(schedule.TimeTo)
	if err != nil {
		return nil, fmt.Errorf("invalid time_to in schedule: %s", schedule.TimeTo)
	}

	dayStart := parsedDate
	dayEnd := parsedDate.Add(24 * time.Hour)
	existingAppointments, err := s.appointmentRepo.GetByDoctorAndDateRange(ctx, doctorID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	bookedSlots := make(map[string]bool)
	for _, a := range existingAppointments {
		bookedSlots[a.ScheduledAt.Format("15:04")] = true
	}

	duration := time.Duration(entity.AppointmentDurationMinutes) * time.Minute
	var slots []entity.TimeSlot
	now := time.Now()

	for t := timeFrom; t.Add(duration).Before(timeTo) || t.Add(duration).Equal(timeTo); t = t.Add(duration) {
		slotStart := t.Format("15:04")
		slotEnd := t.Add(duration).Format("15:04")

		available := !bookedSlots[slotStart]

		if parsedDate.Year() == now.Year() && parsedDate.YearDay() == now.YearDay() {
			slotDateTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
			if slotDateTime.Before(now) {
				available = false
			}
		}

		slots = append(slots, entity.TimeSlot{
			StartTime: slotStart,
			EndTime:   slotEnd,
			Available: available,
		})
	}

	return &entity.AvailableSlotsResponse{
		DoctorID: doctorID,
		Date:     date,
		Slots:    slots,
	}, nil
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, patientID int, req *entity.AppointmentCreateRequest) (*entity.Appointment, error) {
	patient, err := s.userRepo.GetByID(ctx, patientID)
	if err != nil {
		return nil, errors.New("patient not found")
	}

	doctor, err := s.doctorRepo.GetByID(ctx, req.DoctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return nil, errors.New("invalid scheduled_at format, use RFC3339")
	}

	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("scheduled_at must be in the future")
	}

	isoDay := int(scheduledAt.Weekday())
	if isoDay == 0 {
		isoDay = 7
	}

	schedule, err := s.scheduleRepo.GetByDoctorAndDay(ctx, req.DoctorID, isoDay)
	if err != nil {
		return nil, errors.New("doctor does not work on this day")
	}

	timeFrom, _ := parseTime(schedule.TimeFrom)
	timeTo, _ := parseTime(schedule.TimeTo)
	appointmentTime, _ := time.Parse("15:04", scheduledAt.Format("15:04"))
	appointmentEnd := appointmentTime.Add(time.Duration(entity.AppointmentDurationMinutes) * time.Minute)

	if appointmentTime.Before(timeFrom) || appointmentEnd.After(timeTo) {
		return nil, fmt.Errorf("appointment time must be between %s and %s", schedule.TimeFrom, schedule.TimeTo)
	}

	dayStart := time.Date(scheduledAt.Year(), scheduledAt.Month(), scheduledAt.Day(), 0, 0, 0, 0, scheduledAt.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	existing, err := s.appointmentRepo.GetByDoctorAndDateRange(ctx, req.DoctorID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}

	for _, a := range existing {
		if a.ScheduledAt.Format("15:04") == scheduledAt.Format("15:04") {
			return nil, errors.New("this time slot is already booked")
		}
	}

	appointment := &entity.Appointment{
		PatientID:   patientID,
		DoctorID:    req.DoctorID,
		ScheduledAt: scheduledAt,
		Status:      "pending",
		Notes:       req.Notes,
	}

	created, err := s.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Уведомление пациенту
	s.emailService.NotifyPatientAppointmentCreated(patient.Email, patient.Username, doctor.Fullname, scheduledAt)

	// Уведомление врачу
	if doctor.UserID != nil {
		doctorUser, err := s.userRepo.GetByID(ctx, *doctor.UserID)
		if err == nil {
			s.emailService.NotifyDoctorNewAppointment(doctorUser.Email, doctor.Fullname, patient.Username, scheduledAt)
		}
	}

	return created, nil
}

func (s *AppointmentService) GetAppointmentByID(ctx context.Context, id int) (*entity.Appointment, error) {
	return s.appointmentRepo.GetByID(ctx, id)
}

func (s *AppointmentService) GetAppointmentsByPatient(ctx context.Context, patientID int) ([]entity.Appointment, error) {
	if _, err := s.userRepo.GetByID(ctx, patientID); err != nil {
		return nil, errors.New("patient not found")
	}
	appointments, err := s.appointmentRepo.GetByPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}
	s.enrichWithDoctors(ctx, appointments)
	return appointments, nil
}

func (s *AppointmentService) GetAppointmentsByDoctor(ctx context.Context, doctorID int) ([]entity.Appointment, error) {
	if _, err := s.doctorRepo.GetByID(ctx, doctorID); err != nil {
		return nil, errors.New("doctor not found")
	}
	appointments, err := s.appointmentRepo.GetByDoctor(ctx, doctorID)
	if err != nil {
		return nil, err
	}
	s.enrichWithDoctors(ctx, appointments)
	return appointments, nil
}

func (s *AppointmentService) GetAllAppointments(ctx context.Context) ([]entity.Appointment, error) {
	appointments, err := s.appointmentRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	s.enrichWithDoctors(ctx, appointments)
	return appointments, nil
}

func (s *AppointmentService) GetAllAppointmentsPaginated(ctx context.Context, limit, offset int) ([]entity.Appointment, int, error) {
	appointments, total, err := s.appointmentRepo.GetAllPaginated(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	s.enrichWithDoctors(ctx, appointments)
	return appointments, total, nil
}

func (s *AppointmentService) enrichWithDoctors(ctx context.Context, appointments []entity.Appointment) {
	cache := make(map[int]*entity.Doctor)
	for i := range appointments {
		id := appointments[i].DoctorID
		if _, ok := cache[id]; !ok {
			doctor, err := s.doctorRepo.GetByID(ctx, id)
			if err == nil {
				specs, _ := s.doctorRepo.GetSpecializations(ctx, id)
				doctor.Specializations = specs
				cache[id] = doctor
			}
		}
		appointments[i].Doctor = cache[id]
	}
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, appointmentID int, patientID int, role string, req *entity.AppointmentUpdateRequest) (*entity.Appointment, error) {
	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}

	if role != "admin" && appointment.PatientID != patientID {
		return nil, errors.New("not authorized")
	}

	oldStatus := appointment.Status
	oldScheduledAt := appointment.ScheduledAt

	if req.Status != nil {
		if !entity.IsValidTransition(appointment.Status, *req.Status) {
			return nil, fmt.Errorf("cannot change status from '%s' to '%s'", appointment.Status, *req.Status)
		}
		appointment.Status = *req.Status
	}

	if req.ScheduledAt != nil {
		scheduledAt, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return nil, errors.New("invalid scheduled_at format, use RFC3339")
		}
		if scheduledAt.Before(time.Now()) {
			return nil, errors.New("scheduled_at must be in the future")
		}

		dayStart := time.Date(scheduledAt.Year(), scheduledAt.Month(), scheduledAt.Day(), 0, 0, 0, 0, scheduledAt.Location())
		dayEnd := dayStart.Add(24 * time.Hour)
		existing, _ := s.appointmentRepo.GetByDoctorAndDateRange(ctx, appointment.DoctorID, dayStart, dayEnd)
		for _, a := range existing {
			if a.ID != appointmentID && a.ScheduledAt.Format("15:04") == scheduledAt.Format("15:04") {
				return nil, errors.New("this time slot is already booked")
			}
		}

		appointment.ScheduledAt = scheduledAt
	}

	if req.Notes != nil {
		appointment.Notes = req.Notes
	}

	updated, err := s.appointmentRepo.Update(ctx, appointmentID, appointment)
	if err != nil {
		return nil, err
	}

	go s.sendUpdateNotifications(ctx, updated, oldStatus, oldScheduledAt)

	return updated, nil
}

func (s *AppointmentService) CancelAppointment(ctx context.Context, appointmentID int, patientID int, role string) error {
	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return err
	}

	if role != "admin" && appointment.PatientID != patientID {
		return errors.New("not authorized")
	}

	if !entity.IsValidTransition(appointment.Status, "canceled") {
		return fmt.Errorf("cannot cancel appointment with status '%s'", appointment.Status)
	}

	appointment.Status = "canceled"
	_, err = s.appointmentRepo.Update(ctx, appointmentID, appointment)
	if err != nil {
		return err
	}

	go s.sendCancelNotifications(ctx, appointment)

	return nil
}

func (s *AppointmentService) AddAppointmentResult(ctx context.Context, appointmentID int, userID int, resultReq *entity.AppointmentResultRequest) (*entity.Appointment, error) {
	doctor, err := s.doctorRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("doctor profile not found for this user")
	}

	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}

	if appointment.DoctorID != doctor.ID {
		return nil, errors.New("not authorized: this appointment belongs to another doctor")
	}

	appointment.Results = &resultReq.Results
	appointment.Status = resultReq.Status

	updated, err := s.appointmentRepo.Update(ctx, appointmentID, appointment)
	if err != nil {
		return nil, err
	}

	patient, err := s.userRepo.GetByID(ctx, appointment.PatientID)
	if err == nil {
		s.emailService.NotifyPatientResultsReady(patient.Email, patient.Username, doctor.Fullname, updated.ScheduledAt)
	}

	return updated, nil
}

func (s *AppointmentService) sendUpdateNotifications(ctx context.Context, appointment *entity.Appointment, oldStatus string, oldScheduledAt time.Time) {
	patient, err := s.userRepo.GetByID(ctx, appointment.PatientID)
	if err != nil {
		return
	}

	doctor, err := s.doctorRepo.GetByID(ctx, appointment.DoctorID)
	if err != nil {
		return
	}

	// Подтверждение
	if oldStatus == "pending" && appointment.Status == "confirmed" {
		s.emailService.NotifyPatientAppointmentConfirmed(patient.Email, patient.Username, doctor.Fullname, appointment.ScheduledAt)
	}

	// Перенос
	if !oldScheduledAt.Equal(appointment.ScheduledAt) {
		s.emailService.NotifyPatientAppointmentRescheduled(patient.Email, patient.Username, doctor.Fullname, oldScheduledAt, appointment.ScheduledAt)
	}
}

func (s *AppointmentService) sendCancelNotifications(ctx context.Context, appointment *entity.Appointment) {
	patient, err := s.userRepo.GetByID(ctx, appointment.PatientID)
	if err != nil {
		return
	}

	doctor, err := s.doctorRepo.GetByID(ctx, appointment.DoctorID)
	if err != nil {
		return
	}

	s.emailService.NotifyPatientAppointmentCanceled(patient.Email, patient.Username, doctor.Fullname, appointment.ScheduledAt)

	if doctor.UserID != nil {
		doctorUser, err := s.userRepo.GetByID(ctx, *doctor.UserID)
		if err == nil {
			s.emailService.NotifyDoctorAppointmentCanceled(doctorUser.Email, doctor.Fullname, patient.Username, appointment.ScheduledAt)
		}
	}
}

func (s *AppointmentService) AdminCreateAppointment(ctx context.Context, adminID int, req *entity.AppointmentAdminCreateRequest) (*entity.Appointment, error) {
	// Resolve patient ID and build notes prefix
	var patientID int
	var notesPrefix string

	if req.PatientID != nil {
		if _, err := s.userRepo.GetByID(ctx, *req.PatientID); err != nil {
			return nil, errors.New("patient not found")
		}
		patientID = *req.PatientID
	} else if req.GuestName != nil {
		patientID = adminID
		phone := ""
		if req.GuestPhone != nil {
			phone = ", " + *req.GuestPhone
		}
		notesPrefix = "[Callback: " + *req.GuestName + phone + "] "
	} else {
		return nil, errors.New("either patient_id or guest_name must be provided")
	}

	// Build combined notes
	var combinedNotes *string
	if notesPrefix != "" || req.Notes != nil {
		combined := notesPrefix
		if req.Notes != nil {
			combined += *req.Notes
		}
		combinedNotes = &combined
	}

	// Validate doctor
	doctor, err := s.doctorRepo.GetByID(ctx, req.DoctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}

	// Parse and validate scheduled time
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return nil, errors.New("invalid scheduled_at format, use RFC3339")
	}
	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("scheduled_at must be in the future")
	}

	// doctor's schedule
	isoDay := int(scheduledAt.Weekday())
	if isoDay == 0 {
		isoDay = 7
	}
	schedule, err := s.scheduleRepo.GetByDoctorAndDay(ctx, req.DoctorID, isoDay)
	if err != nil {
		return nil, errors.New("doctor does not work on this day")
	}

	timeFrom, _ := parseTime(schedule.TimeFrom)
	timeTo, _ := parseTime(schedule.TimeTo)
	appointmentTime, _ := time.Parse("15:04", scheduledAt.Format("15:04"))
	appointmentEnd := appointmentTime.Add(time.Duration(entity.AppointmentDurationMinutes) * time.Minute)

	if appointmentTime.Before(timeFrom) || appointmentEnd.After(timeTo) {
		return nil, fmt.Errorf("appointment time must be between %s and %s", schedule.TimeFrom, schedule.TimeTo)
	}

	// slot availability
	dayStart := time.Date(scheduledAt.Year(), scheduledAt.Month(), scheduledAt.Day(), 0, 0, 0, 0, scheduledAt.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	existing, err := s.appointmentRepo.GetByDoctorAndDateRange(ctx, req.DoctorID, dayStart, dayEnd)
	if err != nil {
		return nil, err
	}
	for _, a := range existing {
		if a.ScheduledAt.Format("15:04") == scheduledAt.Format("15:04") {
			return nil, errors.New("this time slot is already booked")
		}
	}

	appointment := &entity.Appointment{
		PatientID:   patientID,
		DoctorID:    req.DoctorID,
		ScheduledAt: scheduledAt,
		Status:      "pending",
		Notes:       combinedNotes,
	}

	created, err := s.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Send email only in existing-patient mode
	if req.PatientID != nil {
		patient, err := s.userRepo.GetByID(ctx, patientID)
		if err == nil {
			s.emailService.NotifyPatientAppointmentCreated(patient.Email, patient.Username, doctor.Fullname, scheduledAt)
		}
	}

	return created, nil
}

func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		t, err = time.Parse("15:04:05", s)
	}
	return t, err
}
