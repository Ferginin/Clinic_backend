package utils

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\\-]+@[a-zA-Z0-9.\\-]+\\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\\+?[1-9]\\d{1,14}$`)
)

func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone format")
	}
	return nil
}

func ValidateTimeSlot(from, to string) error {
	fromTime, err := time.Parse("15:04", from)
	if err != nil {
		return errors.New("invalid time_from format (expected HH:MM)")
	}

	toTime, err := time.Parse("15:04", to)
	if err != nil {
		return errors.New("invalid time_to format (expected HH:MM)")
	}

	if !fromTime.Before(toTime) {
		return errors.New("time_from must be before time_to")
	}

	return nil
}

func ValidateDayOfWeek(day int) error {
	if day < 1 || day > 7 {
		return errors.New("day must be between 1 (Monday) and 7 (Sunday)")
	}
	return nil
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
