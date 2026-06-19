package service

import (
	"Clinic_backend/config"
	"fmt"
	"log/slog"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailServiceInterface interface {
	// Уведомления пациенту
	NotifyPatientAppointmentCreated(patientEmail, patientName, doctorName string, scheduledAt time.Time)
	NotifyPatientAppointmentConfirmed(patientEmail, patientName, doctorName string, scheduledAt time.Time)
	NotifyPatientAppointmentCanceled(patientEmail, patientName, doctorName string, scheduledAt time.Time)
	NotifyPatientAppointmentRescheduled(patientEmail, patientName, doctorName string, oldTime, newTime time.Time)
	NotifyPatientResultsReady(patientEmail, patientName, doctorName string, scheduledAt time.Time)

	// Уведомления врачу
	NotifyDoctorNewAppointment(doctorEmail, doctorName, patientName string, scheduledAt time.Time)
	NotifyDoctorAppointmentCanceled(doctorEmail, doctorName, patientName string, scheduledAt time.Time)
}

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailServiceInterface {
	return &EmailService{cfg: cfg}
}

// --- Уведомления пациенту ---

func (s *EmailService) NotifyPatientAppointmentCreated(patientEmail, patientName, doctorName string, scheduledAt time.Time) {
	subject := "Запись на приём подтверждена"
	body := s.buildHTML(
		patientName,
		"Запись на приём создана",
		fmt.Sprintf("Вы успешно записались на приём к врачу <strong>%s</strong>.", doctorName),
		formatDateTime(scheduledAt),
		"#4CAF50",
	)
	go s.send(patientEmail, subject, body)
}

func (s *EmailService) NotifyPatientAppointmentConfirmed(patientEmail, patientName, doctorName string, scheduledAt time.Time) {
	subject := "Приём подтверждён"
	body := s.buildHTML(
		patientName,
		"Ваш приём подтверждён",
		fmt.Sprintf("Врач <strong>%s</strong> подтвердил ваш приём.", doctorName),
		formatDateTime(scheduledAt),
		"#2196F3",
	)
	go s.send(patientEmail, subject, body)
}

func (s *EmailService) NotifyPatientAppointmentCanceled(patientEmail, patientName, doctorName string, scheduledAt time.Time) {
	subject := "Приём отменён"
	body := s.buildHTML(
		patientName,
		"Приём отменён",
		fmt.Sprintf("Ваш приём к врачу <strong>%s</strong> был отменён.", doctorName),
		formatDateTime(scheduledAt),
		"#f44336",
	)
	go s.send(patientEmail, subject, body)
}

func (s *EmailService) NotifyPatientAppointmentRescheduled(patientEmail, patientName, doctorName string, oldTime, newTime time.Time) {
	subject := "Время приёма изменено"
	body := s.buildHTML(
		patientName,
		"Время приёма изменено",
		fmt.Sprintf(
			"Время вашего приёма к врачу <strong>%s</strong> было изменено.<br><br>"+
				"<s>Старое время: %s</s><br>"+
				"<strong>Новое время: %s</strong>",
			doctorName, formatDateTime(oldTime), formatDateTime(newTime),
		),
		"",
		"#FF9800",
	)
	go s.send(patientEmail, subject, body)
}

func (s *EmailService) NotifyPatientResultsReady(patientEmail, patientName, doctorName string, scheduledAt time.Time) {
	subject := "Результаты приёма готовы"
	body := s.buildHTML(
		patientName,
		"Результаты приёма готовы",
		fmt.Sprintf("Результаты вашего приёма у врача <strong>%s</strong> доступны в личном кабинете.", doctorName),
		formatDateTime(scheduledAt),
		"#9C27B0",
	)
	go s.send(patientEmail, subject, body)
}

// --- Уведомления врачу ---

func (s *EmailService) NotifyDoctorNewAppointment(doctorEmail, doctorName, patientName string, scheduledAt time.Time) {
	subject := "Новая запись на приём"
	body := s.buildHTML(
		doctorName,
		"Новая запись на приём",
		fmt.Sprintf("Пациент <strong>%s</strong> записался к вам на приём.", patientName),
		formatDateTime(scheduledAt),
		"#4CAF50",
	)
	go s.send(doctorEmail, subject, body)
}

func (s *EmailService) NotifyDoctorAppointmentCanceled(doctorEmail, doctorName, patientName string, scheduledAt time.Time) {
	subject := "Запись на приём отменена"
	body := s.buildHTML(
		doctorName,
		"Запись отменена",
		fmt.Sprintf("Пациент <strong>%s</strong> отменил запись на приём.", patientName),
		formatDateTime(scheduledAt),
		"#f44336",
	)
	go s.send(doctorEmail, subject, body)
}

// --- Внутренние методы ---

func (s *EmailService) send(to, subject, htmlBody string) {
	if s.cfg.Env.SMTPHost == "" {
		slog.Warn("SMTP not configured, skipping email", "to", to, "subject", subject)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Env.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.cfg.Env.SMTPHost, s.cfg.Env.SMTPPort, s.cfg.Env.SMTPUsername, s.cfg.Env.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		slog.Error("Failed to send email", "to", to, "subject", subject, "error", err)
	} else {
		slog.Info("Email sent", "to", to, "subject", subject)
	}
}

func (s *EmailService) buildHTML(recipientName, title, message, dateTime, accentColor string) string {
	if accentColor == "" {
		accentColor = "#1976D2"
	}

	dateBlock := ""
	if dateTime != "" {
		dateBlock = fmt.Sprintf(`
			<div style="background: #f5f5f5; border-radius: 8px; padding: 12px 16px; margin: 16px 0; text-align: center;">
				<span style="color: #666; font-size: 13px;">Дата и время</span><br>
				<strong style="font-size: 18px; color: #333;">%s</strong>
			</div>`, dateTime)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0; padding:0; background:#f0f2f5; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="padding: 32px 16px;">
    <tr><td align="center">
      <table width="600" cellpadding="0" cellspacing="0" style="background:#fff; border-radius:12px; overflow:hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
        
        <!-- Header -->
        <tr><td style="background: %s; padding: 24px 32px;">
          <h1 style="margin:0; color:#fff; font-size:20px; font-weight:600;">%s</h1>
        </td></tr>

        <!-- Body -->
        <tr><td style="padding: 24px 32px;">
          <p style="margin:0 0 16px; color:#333; font-size:15px; line-height:1.6;">
            Здравствуйте, <strong>%s</strong>!
          </p>
          <p style="margin:0 0 8px; color:#333; font-size:15px; line-height:1.6;">
            %s
          </p>
          %s
        </td></tr>

        <!-- Footer -->
        <tr><td style="padding: 16px 32px 24px; border-top: 1px solid #eee;">
          <p style="margin:0; color:#999; font-size:12px; line-height:1.5;">
            С уважением, команда клиники<br>
            Это автоматическое уведомление, отвечать на него не нужно.
          </p>
        </td></tr>

      </table>
    </td></tr>
  </table>
</body>
</html>`, accentColor, title, recipientName, message, dateBlock)
}

func formatDateTime(t time.Time) string {
	months := []string{
		"", "января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}

	weekdays := []string{
		"воскресенье", "понедельник", "вторник", "среда",
		"четверг", "пятница", "суббота",
	}

	return fmt.Sprintf("%d %s %d г., %s, %s",
		t.Day(), months[t.Month()], t.Year(),
		weekdays[t.Weekday()],
		t.Format("15:04"),
	)
}
