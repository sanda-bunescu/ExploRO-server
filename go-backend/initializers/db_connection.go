package initializers

import (
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Database *gorm.DB

func ConnectToDB() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	Database = db
}

func MigrateDB() {

	err := Database.AutoMigrate(
		&models.Users{},
		&models.City{},
		&models.UserCities{},
		&models.TouristicAttraction{},
		&models.Group{},
		&models.UserGroup{},
		&models.TripPlan{},
		&models.Itinerary{},
		&models.StopPoint{},
		&models.Expense{},
		&models.Debt{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	print("migration is successful")
}
