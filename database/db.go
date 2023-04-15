package database

import (
	"final-project-mygram/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// host     = os.Getenv("PGHOST")
	// user     = os.Getenv("PGUSER")
	// password = os.Getenv("PGPASSWORD")
	// dbPort   = os.Getenv("PGPORT")
	// dbName   = os.Getenv("PGDATABASE")
	host     = "localhost"
	user     = "postgres"
	password = "1"
	dbPort   = "5432"
	dbName   = "mygram"
	db  *gorm.DB
	err error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, dbPort)
	dsn := config
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database: ", err)
		return
	}

	fmt.Println("success connection to database")
	db.Debug().AutoMigrate(models.User{}, models.Photo{}, models.Comment{}, models.SocialMedia{})
}

func GetDB() *gorm.DB {
	return db
}
