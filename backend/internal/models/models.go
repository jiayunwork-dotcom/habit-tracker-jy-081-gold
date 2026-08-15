package models

import (
	"time"

	"gorm.io/gorm"
)

type FrequencyType string

const (
	FrequencyDaily      FrequencyType = "daily"
	FrequencyWeeklyN    FrequencyType = "weekly_n"
	FrequencySpecific   FrequencyType = "specific_days"
	FrequencyMonthlyN   FrequencyType = "monthly_n"
)

type Habit struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"not null" json:"name"`
	Color               string         `gorm:"not null;default:'#4F46E5'" json:"color"`
	FrequencyType       FrequencyType  `gorm:"not null" json:"frequency_type"`
	FrequencyValue      int            `gorm:"not null;default:1" json:"frequency_value"`
	SpecificDays        string         `json:"specific_days"`
	ReminderTime        *string        `json:"reminder_time"`
	CurrentStreak       int            `gorm:"not null;default:0" json:"current_streak"`
	LongestStreak       int            `gorm:"not null;default:0" json:"longest_streak"`
	TotalCheckins       int            `gorm:"not null;default:0" json:"total_checkins"`
	SortOrder           int            `gorm:"not null;default:0" json:"sort_order"`
	IsArchived          bool           `gorm:"not null;default:false" json:"is_archived"`
	FrequencyModifiedAt *time.Time     `json:"frequency_modified_at"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

type Checkin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	HabitID   uint           `gorm:"not null;index" json:"habit_id"`
	Date      string         `gorm:"not null;index" json:"date"`
	Note      string         `json:"note"`
	IsBackfill bool          `gorm:"not null;default:false" json:"is_backfill"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type AchievementType string

const (
	AchievementStreak      AchievementType = "streak"
	AchievementTotal       AchievementType = "total"
	AchievementHabitCount  AchievementType = "habit_count"
	AchievementSpecial     AchievementType = "special"
)

type Achievement struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Type        AchievementType `gorm:"not null" json:"type"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"not null" json:"description"`
	Target      int            `gorm:"not null" json:"target"`
	Icon        string         `gorm:"not null" json:"icon"`
	CreatedAt   time.Time      `json:"created_at"`
}

type UserAchievement struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AchievementID uint      `gorm:"not null;index" json:"achievement_id"`
	HabitID       *uint     `json:"habit_id"`
	EarnedAt      time.Time `gorm:"not null" json:"earned_at"`
	CreatedAt     time.Time  `json:"created_at"`

	Achievement Achievement `gorm:"foreignKey:AchievementID" json:"achievement"`
}

type Setting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
