package database

import (
	"fmt"
	"habit-tracker/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.Habit{},
		&models.Checkin{},
		&models.Achievement{},
		&models.UserAchievement{},
		&models.Setting{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}

func SeedAchievements() {
	achievements := []models.Achievement{
		{Type: models.AchievementStreak, Name: "初心者", Description: "连续打卡7天", Target: 7, Icon: "🌟"},
		{Type: models.AchievementStreak, Name: "坚持者", Description: "连续打卡21天", Target: 21, Icon: "⭐"},
		{Type: models.AchievementStreak, Name: "月度达人", Description: "连续打卡30天", Target: 30, Icon: "🏆"},
		{Type: models.AchievementStreak, Name: "双月战士", Description: "连续打卡60天", Target: 60, Icon: "🥇"},
		{Type: models.AchievementStreak, Name: "百日英雄", Description: "连续打卡100天", Target: 100, Icon: "👑"},
		{Type: models.AchievementStreak, Name: "年度传奇", Description: "连续打卡365天", Target: 365, Icon: "💎"},

		{Type: models.AchievementTotal, Name: "起步者", Description: "累计打卡50次", Target: 50, Icon: "📝"},
		{Type: models.AchievementTotal, Name: "百次俱乐部", Description: "累计打卡100次", Target: 100, Icon: "💯"},
		{Type: models.AchievementTotal, Name: "五百大师", Description: "累计打卡500次", Target: 500, Icon: "🎯"},
		{Type: models.AchievementTotal, Name: "千次传奇", Description: "累计打卡1000次", Target: 1000, Icon: "🎖️"},

		{Type: models.AchievementHabitCount, Name: "多面手", Description: "同时维护3个活跃习惯", Target: 3, Icon: "🎭"},
		{Type: models.AchievementHabitCount, Name: "习惯达人", Description: "同时维护5个活跃习惯", Target: 5, Icon: "🌈"},
		{Type: models.AchievementHabitCount, Name: "生活大师", Description: "同时维护10个活跃习惯", Target: 10, Icon: "🌅"},

		{Type: models.AchievementSpecial, Name: "完美一周", Description: "一周内所有习惯全部按要求完成", Target: 1, Icon: "✨"},
		{Type: models.AchievementSpecial, Name: "完美一月", Description: "一个月内所有习惯全部按要求完成", Target: 1, Icon: "🌸"},
	}

	for _, a := range achievements {
		var count int64
		DB.Model(&models.Achievement{}).Where("name = ?", a.Name).Count(&count)
		if count == 0 {
			DB.Create(&a)
		}
	}
	log.Println("Achievements seeded successfully")
}
