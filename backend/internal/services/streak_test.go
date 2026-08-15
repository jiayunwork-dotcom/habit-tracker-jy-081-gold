package services

import (
	"testing"
	"time"

	"habit-tracker/internal/models"
)

func TestIsDateWithinBackfillWindow_IncludesSevenDaysAgo(t *testing.T) {
	d := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	s := NewStreakService()
	if !s.IsDateWithinBackfillWindow(d) {
		t.Fatal("a date 7 days ago should be within the backfill window")
	}
}

func TestIsDateWithinBackfillWindow_RejectsBadDate(t *testing.T) {
	s := NewStreakService()
	if s.IsDateWithinBackfillWindow("not-a-date") {
		t.Fatal("invalid date string should be rejected")
	}
}

func TestShouldCheckinOnDate_SundayMatchesDaySeven(t *testing.T) {
	s := NewStreakService()
	sunday := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC) // Sunday
	habit := &models.Habit{FrequencyType: models.FrequencySpecific, SpecificDays: "7"}
	ok, err := s.ShouldCheckinOnDate(habit, sunday)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Sunday should match specific day 7")
	}
}

func TestShouldCheckinToday_UnknownFrequencyErrors(t *testing.T) {
	s := NewStreakService()
	_, err := s.ShouldCheckinToday(&models.Habit{FrequencyType: "nope"})
	if err == nil {
		t.Fatal("unknown frequency type should return an error")
	}
}

func TestGetDaysInMonth_FebruaryLeap(t *testing.T) {
	if got := getDaysInMonth(2024, 2); got != 29 {
		t.Fatalf("Feb 2024 days=%d, want 29", got)
	}
}
