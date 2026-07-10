package database

import (
	"fmt"
	"log"
	"time"

	"farm/internal/config"
	"farm/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func Migrate() {
	DB.AutoMigrate(
		&models.User{},
		&models.Species{},
		&models.Breed{},
		&models.Animal{},
		&models.AnimalWeightHistory{},
		&models.Vaccine{},
		&models.AnimalVaccination{},
		&models.InventoryCategory{},
		&models.InventoryItem{},
		&models.InventoryTransaction{},
		&models.AccountHead{},
		&models.AccountTransaction{},
		&models.AnimalPregnancy{},
	)
}
