package database

import (
	"fmt"
	"log"

	"github.com/roto17/zeus/lib/env"
	"gorm.io/driver/postgres" // or use MySQL driver if needed
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// dbDriver := env.GetEnv("dbdriver")
	dbHost := env.GetEnv("dbhost")
	dbPort := env.GetEnv("dbport")
	dbName := env.GetEnv("dbname")
	dbUser := env.GetEnv("dbuser")
	dbPassword := env.GetEnv("dbpassword")
	dbSSLMode := env.GetEnv("dbsslmode")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	log.Println("Database connection established!")
}
