package main

import (
	"log"

	"farm/internal/config"
	"farm/internal/database"
	"farm/internal/models"
	"farm/internal/routes"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate()

	app := fiber.New(fiber.Config{
		AppName: "Farm API",
	})

	routes.Setup(app, cfg)

	seedAdmin(cfg)
	database.Seed()

	log.Printf("Server starting on %s", cfg.AppPort)
	log.Fatal(app.Listen(cfg.AppPort))
}

func seedAdmin(cfg *config.Config) {
	var count int64
	database.DB.Model(&models.User{}).Where("role = ?", "Owner").Count(&count)
	if count == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Name:     "Admin",
			Email:    strPtr("admin@farm.com"),
			Username: strPtr("admin"),
			Password: string(hashed),
			Role:     "Owner",
			Status:   "Active",
		}
		database.DB.Create(&admin)
		log.Println("Default admin created: admin / admin123")
	}
}

func strPtr(s string) *string { return &s }
