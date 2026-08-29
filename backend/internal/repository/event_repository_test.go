package repository

import (
	"testing"
	"time"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

func TestValidateEventRequest(t *testing.T) {
	t.Run("accepts valid range", func(t *testing.T) {
		request := model.EventAdminRequest{
			Name:      "Summer Open",
			VenueName: "Elite Club",
			StartDate: "2026-08-01",
			EndDate:   "2026-08-05",
		}

		if err := validateEventRequest(request); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("rejects end date before start date", func(t *testing.T) {
		request := model.EventAdminRequest{
			Name:      "Summer Open",
			VenueName: "Elite Club",
			StartDate: "2026-08-10",
			EndDate:   "2026-08-05",
		}

		if err := validateEventRequest(request); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("rejects invalid date format", func(t *testing.T) {
		request := model.EventAdminRequest{
			Name:      "Summer Open",
			VenueName: "Elite Club",
			StartDate: "not-a-date",
			EndDate:   "2026-08-05",
		}

		if err := validateEventRequest(request); err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("rejects deadline equal to start date", func(t *testing.T) {
		deadline := "2026-08-10T00:00:00Z"
		request := model.EventAdminRequest{
			Name:                 "Summer Open",
			VenueName:            "Elite Club",
			StartDate:            "2026-08-10",
			EndDate:              "2026-08-11",
			RegistrationDeadline: &deadline,
		}

		if err := validateEventRequest(request); err == nil {
			t.Fatal("expected validation error")
		}
	})
}

func TestParseDate(t *testing.T) {
	date, err := parseDate("2026-08-05")
	if err != nil {
		t.Fatalf("expected date to parse, got %v", err)
	}

	if !date.Equal(time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected parsed date: %v", date)
	}
}
