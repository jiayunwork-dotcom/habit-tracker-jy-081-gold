package handlers

import (
	"habit-tracker/internal/database"
	"habit-tracker/internal/models"
	"habit-tracker/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	streakService      *services.StreakService
	achievementService *services.AchievementService
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{
		streakService:      services.NewStreakService(),
		achievementService: services.NewAchievementService(),
	}
}

func (h *StatsHandler) GetHeatmap(c *gin.Context) {
	data, err := h.streakService.GetHeatmapData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *StatsHandler) GetOverview(c *gin.Context) {
	dueHabits, err := h.streakService.GetTodayDueHabits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")
	checkedCount := 0
	for _, h := range dueHabits {
		var checkin models.Checkin
		hasCheckin := database.DB.Where("habit_id = ? AND date = ?", h.ID, today).First(&checkin).Error == nil
		if hasCheckin {
			checkedCount++
		}
	}

	var totalHabits int64
	database.DB.Model(&models.Habit{}).Where("is_archived = ?", false).Count(&totalHabits)

	c.JSON(http.StatusOK, gin.H{
		"due_today":       len(dueHabits),
		"completed_today": checkedCount,
		"total_habits":    totalHabits,
		"completion_rate": 0,
	})
}

func (h *StatsHandler) GetAchievements(c *gin.Context) {
	achievements, err := h.achievementService.GetAllAchievementsWithProgress()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, achievements)
}
