package services

import (
	"fmt"
	"habit-tracker/internal/database"
	"habit-tracker/internal/models"
	"strconv"
	"strings"
	"time"
)

type StreakService struct{}

func NewStreakService() *StreakService {
	return &StreakService{}
}

func (s *StreakService) CalculateStreak(habit *models.Habit) (int, error) {
	var checkins []models.Checkin
	query := database.DB.Where("habit_id = ?", habit.ID)
	
	if habit.FrequencyModifiedAt != nil {
		query = query.Where("created_at >= ?", habit.FrequencyModifiedAt)
	}
	
	if err := query.Order("date desc").Find(&checkins).Error; err != nil {
		return 0, err
	}

	if len(checkins) == 0 {
		return 0, nil
	}

	checkinMap := make(map[string]bool)
	for _, c := range checkins {
		checkinMap[c.Date] = true
	}

	switch habit.FrequencyType {
	case models.FrequencyDaily:
		return s.calculateDailyStreak(checkinMap)
	case models.FrequencyWeeklyN:
		return s.calculateWeeklyNStreak(checkinMap, habit.FrequencyValue)
	case models.FrequencySpecific:
		return s.calculateSpecificDaysStreak(checkinMap, habit.SpecificDays)
	case models.FrequencyMonthlyN:
		return s.calculateMonthlyNStreak(checkinMap, habit.FrequencyValue)
	default:
		return 0, fmt.Errorf("unknown frequency type: %s", habit.FrequencyType)
	}
}

func (s *StreakService) calculateDailyStreak(checkins map[string]bool) (int, error) {
	today := time.Now().Format("2006-01-02")
	streak := 0
	currentDate := today

	for {
		if checkins[currentDate] {
			streak++
		} else {
			break
		}

		date, _ := time.Parse("2006-01-02", currentDate)
		currentDate = date.AddDate(0, 0, -1).Format("2006-01-02")
	}

	return streak, nil
}

func (s *StreakService) calculateWeeklyNStreak(checkins map[string]bool, n int) (int, error) {
	streakDays := 0
	consecutiveFullWeeks := 0
	currentWeekStart := getWeekStart(time.Now())
	today := time.Now()

	for {
		weekEnd := currentWeekStart.AddDate(0, 0, 6)
		isCurrentWeek := !today.Before(currentWeekStart) && !today.After(weekEnd)
		
		weekCheckins := 0
		for i := 0; i < 7; i++ {
			date := currentWeekStart.AddDate(0, 0, i)
			dateStr := date.Format("2006-01-02")
			if checkins[dateStr] {
				weekCheckins++
			}
		}

		if isCurrentWeek {
			if weekCheckins > 0 {
				streakDays += weekCheckins
			} else if consecutiveFullWeeks == 0 {
				break
			}
		} else {
			if weekCheckins >= n {
				consecutiveFullWeeks++
				streakDays += weekCheckins
			} else {
				break
			}
		}

		currentWeekStart = currentWeekStart.AddDate(0, 0, -7)
	}

	if streakDays == 0 {
		for _, hasCheckin := range checkins {
			if hasCheckin {
				return 1, nil
			}
		}
	}

	return streakDays, nil
}

func (s *StreakService) calculateSpecificDaysStreak(checkins map[string]bool, specificDays string) (int, error) {
	days := strings.Split(specificDays, ",")
	dayMap := make(map[int]bool)
	for _, d := range days {
		dayNum, _ := strconv.Atoi(d)
		dayMap[dayNum] = true
	}

	streak := 0
	currentDate := time.Now()

	for {
		weekday := int(currentDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		if dayMap[weekday] {
			if checkins[currentDate.Format("2006-01-02")] {
				streak++
			} else {
				break
			}
		}

		currentDate = currentDate.AddDate(0, 0, -1)
	}

	return streak, nil
}

func (s *StreakService) calculateMonthlyNStreak(checkins map[string]bool, n int) (int, error) {
	streakDays := 0
	consecutiveFullMonths := 0
	currentMonth := time.Now()
	today := time.Now()

	for {
		isCurrentMonth := currentMonth.Year() == today.Year() && currentMonth.Month() == today.Month()
		
		monthCheckins := 0
		daysInMonth := getDaysInMonth(currentMonth.Year(), int(currentMonth.Month()))
		
		for i := 1; i <= daysInMonth; i++ {
			date := time.Date(currentMonth.Year(), currentMonth.Month(), i, 0, 0, 0, 0, time.Local)
			if checkins[date.Format("2006-01-02")] {
				monthCheckins++
			}
		}

		if isCurrentMonth {
			if monthCheckins > 0 {
				streakDays += monthCheckins
			} else if consecutiveFullMonths == 0 {
				break
			}
		} else {
			if monthCheckins >= n {
				consecutiveFullMonths++
				streakDays += monthCheckins
			} else {
				break
			}
		}

		currentMonth = currentMonth.AddDate(0, -1, 0)
	}

	if streakDays == 0 {
		for _, hasCheckin := range checkins {
			if hasCheckin {
				return 1, nil
			}
		}
	}

	return streakDays, nil
}

func (s *StreakService) UpdateHabitStreak(habitID uint) error {
	var habit models.Habit
	if err := database.DB.First(&habit, habitID).Error; err != nil {
		return err
	}

	streak, err := s.CalculateStreak(&habit)
	if err != nil {
		return err
	}

	habit.CurrentStreak = streak
	if streak > habit.LongestStreak {
		habit.LongestStreak = streak
	}

	var total int64
	database.DB.Model(&models.Checkin{}).Where("habit_id = ?", habitID).Count(&total)
	habit.TotalCheckins = int(total)

	return database.DB.Save(&habit).Error
}

func (s *StreakService) GetSortedCheckinDates(habitID uint) ([]string, error) {
	var checkins []models.Checkin
	if err := database.DB.Where("habit_id = ?", habitID).Order("date asc").Find(&checkins).Error; err != nil {
		return nil, err
	}

	dates := make([]string, len(checkins))
	for i, c := range checkins {
		dates[i] = c.Date
	}
	return dates, nil
}

func (s *StreakService) IsDateWithinBackfillWindow(dateStr string) bool {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return true
	}

	today := time.Now()
	sevenDaysAgo := today.AddDate(0, 0, -6)

	return !date.Before(sevenDaysAgo) && !date.After(today)
}

func (s *StreakService) ShouldCheckinToday(habit *models.Habit) (bool, error) {
	today := time.Now()
	weekday := int(today.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	switch habit.FrequencyType {
	case models.FrequencyDaily:
		return true, nil
	case models.FrequencySpecific:
		days := strings.Split(habit.SpecificDays, ",")
		for _, d := range days {
			dayNum, _ := strconv.Atoi(d)
			if dayNum == weekday {
				return true, nil
			}
		}
		return false, nil
	case models.FrequencyWeeklyN:
		return true, nil
	case models.FrequencyMonthlyN:
		return true, nil
	default:
		return false, nil
	}
}

func (s *StreakService) GetTodayDueHabits() ([]models.Habit, error) {
	var habits []models.Habit
	if err := database.DB.Where("is_archived = ?", false).Find(&habits).Error; err != nil {
		return nil, err
	}

	var dueHabits []models.Habit
	for _, h := range habits {
		should, _ := s.ShouldCheckinToday(&h)
		if should {
			dueHabits = append(dueHabits, h)
		}
	}

	return dueHabits, nil
}

func (s *StreakService) GetHeatmapData() (map[string]int, error) {
	var habits []models.Habit
	database.DB.Where("is_archived = ?", false).Find(&habits)

	heatmap := make(map[string]int)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	for d := oneYearAgo; !d.After(time.Now()); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		totalDue := 0
		totalDone := 0

		for _, h := range habits {
			var checkin models.Checkin
			hasCheckin := database.DB.Where("habit_id = ? AND date = ?", h.ID, dateStr).First(&checkin).Error == nil

			shouldCheckin, _ := s.ShouldCheckinOnDate(&h, d)
			if shouldCheckin {
				totalDue++
				if hasCheckin {
					totalDone++
				}
			}
		}

		if totalDue > 0 {
			heatmap[dateStr] = (totalDone * 100) / totalDue
		}
	}

	return heatmap, nil
}

func (s *StreakService) ShouldCheckinOnDate(habit *models.Habit, date time.Time) (bool, error) {
	weekday := int(date.Weekday())

	switch habit.FrequencyType {
	case models.FrequencyDaily:
		return true, nil
	case models.FrequencySpecific:
		days := strings.Split(habit.SpecificDays, ",")
		for _, d := range days {
			dayNum, _ := strconv.Atoi(d)
			if dayNum == weekday {
				return true, nil
			}
		}
		return false, nil
	default:
		return true, nil
	}
}

func (s *StreakService) GetHabitStats(habitID uint) (map[string]interface{}, error) {
	var habit models.Habit
	if err := database.DB.First(&habit, habitID).Error; err != nil {
		return nil, err
	}

	dates, err := s.GetSortedCheckinDates(habitID)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]interface{})
	stats["current_streak"] = habit.CurrentStreak
	stats["longest_streak"] = habit.LongestStreak
	stats["total_checkins"] = habit.TotalCheckins

	weekStart := getWeekStart(time.Now())
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)

	weekCount := 0
	monthCount := 0
	for _, d := range dates {
		date, _ := time.Parse("2006-01-02", d)
		if !date.Before(weekStart) {
			weekCount++
		}
		if !date.Before(monthStart) {
			monthCount++
		}
	}

	var weekDue, monthDue, totalDue int
	switch habit.FrequencyType {
	case models.FrequencyDaily:
		weekDue = 7
		monthDue = getDaysInMonth(time.Now().Year(), int(time.Now().Month()))
		totalDue = int(time.Now().Sub(habit.CreatedAt).Hours() / 24)
	case models.FrequencyWeeklyN:
		weekDue = habit.FrequencyValue
		monthDue = habit.FrequencyValue * 4
		totalDue = int(time.Now().Sub(habit.CreatedAt).Hours()/24/7) * habit.FrequencyValue
	case models.FrequencySpecific:
		days := strings.Split(habit.SpecificDays, ",")
		weekDue = len(days)
		monthDue = len(days) * 4
		totalDue = int(time.Now().Sub(habit.CreatedAt).Hours()/24/7) * len(days)
	case models.FrequencyMonthlyN:
		weekDue = (habit.FrequencyValue + 3) / 4
		monthDue = habit.FrequencyValue
		totalDue = int(time.Now().Sub(habit.CreatedAt).Hours()/24/30) * habit.FrequencyValue
	}

	stats["week_rate"] = 0
	if weekDue > 0 {
		stats["week_rate"] = (weekCount * 100) / weekDue
	}
	stats["month_rate"] = 0
	if monthDue > 0 {
		stats["month_rate"] = (monthCount * 100) / monthDue
	}
	stats["overall_rate"] = 0
	if totalDue > 0 {
		stats["overall_rate"] = (habit.TotalCheckins * 100) / totalDue
	}

	return stats, nil
}

func getWeekStart(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
}

func getDaysInMonth(year, month int) int {
	return time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.Local).Day()
}

func GetRandomEncouragement() string {
	encouragements := []string{
		"太棒了！继续保持！💪",
		"你正在变得更好！🌟",
		"每一步都很重要！✨",
		"坚持就是胜利！🏆",
		"今天的你比昨天更强！💫",
		"好习惯正在养成！🌱",
		"你做到了！为你骄傲！🎉",
		"继续加油，你可以的！💯",
	}
	return encouragements[time.Now().UnixNano()%int64(len(encouragements))]
}
