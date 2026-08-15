package services

import (
	"habit-tracker/internal/database"
	"habit-tracker/internal/models"
	"time"
)

type AchievementService struct{}

func NewAchievementService() *AchievementService {
	return &AchievementService{}
}

func (s *AchievementService) CheckAndAwardAchievements(habitID uint) ([]models.UserAchievement, error) {
	var newAchievements []models.UserAchievement
	var habit models.Habit

	if err := database.DB.First(&habit, habitID).Error; err != nil {
		return nil, err
	}

	streakAchievements, err := s.checkStreakAchievements(&habit)
	if err != nil {
		return nil, err
	}
	newAchievements = append(newAchievements, streakAchievements...)

	totalAchievements, err := s.checkTotalAchievements(&habit)
	if err != nil {
		return nil, err
	}
	newAchievements = append(newAchievements, totalAchievements...)

	habitCountAchievements, err := s.checkHabitCountAchievements()
	if err != nil {
		return nil, err
	}
	newAchievements = append(newAchievements, habitCountAchievements...)

	return newAchievements, nil
}

func (s *AchievementService) checkStreakAchievements(habit *models.Habit) ([]models.UserAchievement, error) {
	var achievements []models.Achievement
	if err := database.DB.Where("type = ?", models.AchievementStreak).Find(&achievements).Error; err != nil {
		return nil, err
	}

	var earned []models.UserAchievement
	for _, a := range achievements {
		if habit.CurrentStreak >= a.Target {
			var existing models.UserAchievement
			result := database.DB.Where("achievement_id = ? AND habit_id = ?", a.ID, habit.ID).First(&existing)
			if result.Error != nil {
				userAchievement := models.UserAchievement{
					AchievementID: a.ID,
					HabitID:       &habit.ID,
					EarnedAt:      time.Now(),
				}
				if err := database.DB.Create(&userAchievement).Error; err != nil {
					return nil, err
				}
				userAchievement.Achievement = a
				earned = append(earned, userAchievement)
			}
		}
	}

	return earned, nil
}

func (s *AchievementService) checkTotalAchievements(habit *models.Habit) ([]models.UserAchievement, error) {
	var achievements []models.Achievement
	if err := database.DB.Where("type = ?", models.AchievementTotal).Find(&achievements).Error; err != nil {
		return nil, err
	}

	var earned []models.UserAchievement
	for _, a := range achievements {
		if habit.TotalCheckins >= a.Target {
			var existing models.UserAchievement
			result := database.DB.Where("achievement_id = ? AND habit_id = ?", a.ID, habit.ID).First(&existing)
			if result.Error != nil {
				userAchievement := models.UserAchievement{
					AchievementID: a.ID,
					HabitID:       &habit.ID,
					EarnedAt:      time.Now(),
				}
				if err := database.DB.Create(&userAchievement).Error; err != nil {
					return nil, err
				}
				userAchievement.Achievement = a
				earned = append(earned, userAchievement)
			}
		}
	}

	return earned, nil
}

func (s *AchievementService) checkHabitCountAchievements() ([]models.UserAchievement, error) {
	var count int64
	if err := database.DB.Model(&models.Habit{}).Where("is_archived = ?", false).Count(&count).Error; err != nil {
		return nil, err
	}

	var achievements []models.Achievement
	if err := database.DB.Where("type = ?", models.AchievementHabitCount).Find(&achievements).Error; err != nil {
		return nil, err
	}

	var earned []models.UserAchievement
	for _, a := range achievements {
		if int(count) >= a.Target {
			var existing models.UserAchievement
			result := database.DB.Where("achievement_id = ?", a.ID).First(&existing)
			if result.Error != nil {
				userAchievement := models.UserAchievement{
					AchievementID: a.ID,
					EarnedAt:      time.Now(),
				}
				if err := database.DB.Create(&userAchievement).Error; err != nil {
					return nil, err
				}
				userAchievement.Achievement = a
				earned = append(earned, userAchievement)
			}
		}
	}

	return earned, nil
}

func (s *AchievementService) GetAllAchievementsWithProgress() ([]map[string]interface{}, error) {
	var achievements []models.Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		return nil, err
	}

	var userAchievements []models.UserAchievement
	if err := database.DB.Preload("Achievement").Find(&userAchievements).Error; err != nil {
		return nil, err
	}

	earnedMap := make(map[uint]bool)
	for _, ua := range userAchievements {
		earnedMap[ua.AchievementID] = true
	}

	var habits []models.Habit
	database.DB.Where("is_archived = ?", false).Find(&habits)

	result := make([]map[string]interface{}, len(achievements))
	for i, a := range achievements {
		item := map[string]interface{}{
			"id":          a.ID,
			"type":        a.Type,
			"name":        a.Name,
			"description": a.Description,
			"target":      a.Target,
			"icon":        a.Icon,
			"earned":      earnedMap[a.ID],
			"progress":    0,
		}

		switch a.Type {
		case models.AchievementStreak:
			maxStreak := 0
			for _, h := range habits {
				if h.CurrentStreak > maxStreak {
					maxStreak = h.CurrentStreak
				}
			}
			item["progress"] = maxStreak
		case models.AchievementTotal:
			total := 0
			for _, h := range habits {
				total += h.TotalCheckins
			}
			item["progress"] = total
		case models.AchievementHabitCount:
			item["progress"] = len(habits)
		case models.AchievementSpecial:
			item["progress"] = 0
		}

		if earnedMap[a.ID] {
			item["progress"] = a.Target
		}

		result[i] = item
	}

	return result, nil
}
