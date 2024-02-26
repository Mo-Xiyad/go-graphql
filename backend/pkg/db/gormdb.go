package db

import (
	"log"
	"os"
	"server/config"
	"server/pkg/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB(conf *config.Config) (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// db, err := gorm.Open(mysql.Open(conf.Database.URL), &gorm.Config{
	db, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to PlanetScale: %v", err)
	}

	log.Println("Database connection successfully established")
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Company{})

	return db, nil
}
