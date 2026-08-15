package handlers

import (
	"habit-tracker/internal/database"
	"habit-tracker/internal/models"
	"habit-tracker/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type HabitHandler struct {
	streakService      *services.StreakService
	achievementService *services.AchievementService
}

func NewHabitHandler() *HabitHandler {
	return &HabitHandler{
		streakService:      services.NewStreakService(),
		achievementService: services.NewAchievementService(),
	}
}

func (h *HabitHandler) GetHabits(c *gin.Context) {
	var habits []models.Habit
	archived := c.Query("archived") == "true"

	if err := database.DB.Where("is_archived = ?", archived).Order("sort_order asc, created_at desc").Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().Format("2006-01-02")
	result := make([]map[string]interface{}, len(habits))
	for i, habit := range habits {
		var checkin models.Checkin
		hasCheckin := database.DB.Where("habit_id = ? AND date = ?", habit.ID, today).First(&checkin).Error == nil

		result[i] = map[string]interface{}{
			"id":               habit.ID,
			"name":             habit.Name,
			"color":            habit.Color,
			"frequency_type":   habit.FrequencyType,
			"frequency_value":  habit.FrequencyValue,
			"specific_days":    habit.SpecificDays,
			"reminder_time":    habit.ReminderTime,
			"current_streak":   habit.CurrentStreak,
			"longest_streak":   habit.LongestStreak,
			"total_checkins":   habit.TotalCheckins,
			"sort_order":       habit.SortOrder,
			"is_archived":      habit.IsArchived,
			"created_at":       habit.CreatedAt,
			"today_checked_in": hasCheckin,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *HabitHandler) GetHabit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var habit models.Habit
	if err := database.DB.First(&habit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}
	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) CreateHabit(c *gin.Context) {
	var input struct {
		Name          string         `json:"name" binding:"required"`
		Color         string         `json:"color"`
		FrequencyType string         `json:"frequency_type" binding:"required"`
		FrequencyValue *int          `json:"frequency_value"`
		SpecificDays  string         `json:"specific_days"`
		ReminderTime  *string        `json:"reminder_time"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Color == "" {
		input.Color = "#4F46E5"
	}

	freqValue := 1
	if input.FrequencyValue != nil {
		freqValue = *input.FrequencyValue
	}

	var maxSortOrder int
	database.DB.Model(&models.Habit{}).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSortOrder)

	habit := models.Habit{
		Name:           input.Name,
		Color:          input.Color,
		FrequencyType:  models.FrequencyType(input.FrequencyType),
		FrequencyValue: freqValue,
		SpecificDays:   input.SpecificDays,
		ReminderTime:   input.ReminderTime,
		SortOrder:      maxSortOrder + 1,
	}

	if err := database.DB.Create(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.achievementService.CheckAndAwardAchievements(habit.ID)

	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) UpdateHabit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var habit models.Habit
	if err := database.DB.First(&habit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Habit not found"})
		return
	}

	var input struct {
		Name           *string `json:"name"`
		Color          *string `json:"color"`
		FrequencyType  *string `json:"frequency_type"`
		FrequencyValue *int    `json:"frequency_value"`
		SpecificDays   *string `json:"specific_days"`
		ReminderTime   **string `json:"reminder_time"`
		SortOrder      *int    `json:"sort_order"`
		IsArchived     *bool   `json:"is_archived"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	frequencyChanged := false
	if input.FrequencyType != nil && *input.FrequencyType != string(habit.FrequencyType) {
		frequencyChanged = true
	}
	if input.FrequencyValue != nil && *input.FrequencyValue != habit.FrequencyValue {
		frequencyChanged = true
	}
	if input.SpecificDays != nil && *input.SpecificDays != habit.SpecificDays {
		frequencyChanged = true
	}

	if input.Name != nil {
		habit.Name = *input.Name
	}
	if input.Color != nil {
		habit.Color = *input.Color
	}
	if input.FrequencyType != nil {
		habit.FrequencyType = models.FrequencyType(*input.FrequencyType)
	}
	if input.FrequencyValue != nil {
		habit.FrequencyValue = *input.FrequencyValue
	}
	if input.SpecificDays != nil {
		habit.SpecificDays = *input.SpecificDays
	}
	if input.ReminderTime != nil {
		habit.ReminderTime = *input.ReminderTime
	}
	if input.SortOrder != nil {
		habit.SortOrder = *input.SortOrder
	}
	if input.IsArchived != nil {
		habit.IsArchived = *input.IsArchived
	}

	if frequencyChanged {
		now := time.Now()
		habit.FrequencyModifiedAt = &now
		habit.CurrentStreak = 0
	}

	if err := database.DB.Save(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if frequencyChanged {
		h.streakService.UpdateHabitStreak(habit.ID)
	}

	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) DeleteHabit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Delete(&models.Habit{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Where("habit_id = ?", id).Delete(&models.Checkin{})
	c.JSON(http.StatusOK, gin.H{"message": "Habit deleted"})
}

func (h *HabitHandler) CheckIn(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var existing models.Checkin
	result := database.DB.Where("habit_id = ? AND date = ?", id, date).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"message":     "Already checked in",
			"streak":      existing,
			"encouragement": services.GetRandomEncouragement(),
		})
		return
	}

	isBackfill := date != time.Now().Format("2006-01-02")
	if isBackfill && !h.streakService.IsDateWithinBackfillWindow(date) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot backfill more than 7 days ago"})
		return
	}

	var input struct {
		Note string `json:"note"`
	}
	c.ShouldBindJSON(&input)

	checkin := models.Checkin{
		HabitID:   uint(id),
		Date:      date,
		Note:      input.Note,
		IsBackfill: isBackfill,
	}

	if err := database.DB.Create(&checkin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.streakService.UpdateHabitStreak(uint(id))
	newAchievements, _ := h.achievementService.CheckAndAwardAchievements(uint(id))

	var habit models.Habit
	database.DB.First(&habit, id)

	c.JSON(http.StatusOK, gin.H{
		"checkin":      checkin,
		"current_streak": habit.CurrentStreak,
		"encouragement": services.GetRandomEncouragement(),
		"new_achievements": newAchievements,
	})
}

func (h *HabitHandler) CancelCheckIn(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	if err := database.DB.Where("habit_id = ? AND date = ?", id, date).Delete(&models.Checkin{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.streakService.UpdateHabitStreak(uint(id))

	c.JSON(http.StatusOK, gin.H{"message": "Check-in cancelled"})
}

func (h *HabitHandler) GetCheckIns(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var checkins []models.Checkin
	if err := database.DB.Where("habit_id = ?", id).Order("date desc").Find(&checkins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, checkins)
}

func (h *HabitHandler) GetHabitStats(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	stats, err := h.streakService.GetHabitStats(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *HabitHandler) ReorderHabits(c *gin.Context) {
	var input struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, id := range input.IDs {
		database.DB.Model(&models.Habit{}).Where("id = ?", id).Update("sort_order", i+1)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Habits reordered"})
}

func (h *HabitHandler) GetSettings(c *gin.Context) {
	var settings []models.Setting
	if err := database.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, result)
}

func (h *HabitHandler) UpdateSettings(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range input {
		var setting models.Setting
		result := database.DB.Where("key = ?", key).First(&setting)
		if result.Error != nil {
			setting = models.Setting{Key: key, Value: value}
			database.DB.Create(&setting)
		} else {
			setting.Value = value
			database.DB.Save(&setting)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}
