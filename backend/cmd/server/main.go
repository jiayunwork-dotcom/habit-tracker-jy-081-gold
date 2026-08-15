package main

import (
	"habit-tracker/internal/database"
	"habit-tracker/internal/handlers"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database.Connect()
	database.Migrate()
	database.SeedAchievements()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	habitHandler := handlers.NewHabitHandler()
	statsHandler := handlers.NewStatsHandler()

	api := r.Group("/api")
	{
		habits := api.Group("/habits")
		{
			habits.GET("", habitHandler.GetHabits)
			habits.POST("", habitHandler.CreateHabit)
			habits.POST("/reorder", habitHandler.ReorderHabits)
			habits.GET("/:id", habitHandler.GetHabit)
			habits.PUT("/:id", habitHandler.UpdateHabit)
			habits.DELETE("/:id", habitHandler.DeleteHabit)
			habits.POST("/:id/checkin", habitHandler.CheckIn)
			habits.DELETE("/:id/checkin", habitHandler.CancelCheckIn)
			habits.GET("/:id/checkins", habitHandler.GetCheckIns)
			habits.GET("/:id/stats", habitHandler.GetHabitStats)
		}

		stats := api.Group("/stats")
		{
			stats.GET("/heatmap", statsHandler.GetHeatmap)
			stats.GET("/overview", statsHandler.GetOverview)
			stats.GET("/achievements", statsHandler.GetAchievements)
		}

		settings := api.Group("/settings")
		{
			settings.GET("", habitHandler.GetSettings)
			settings.PUT("", habitHandler.UpdateSettings)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
