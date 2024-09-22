package database

import (
	"fmt"

	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/logs"
	"gorm.io/driver/postgres" // or use MySQL driver if needed
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	// dbDriver := config.GetEnv("dbdriver")
	dbHost := config.GetEnv("dbhost")
	dbPort := config.GetEnv("dbport")
	dbName := config.GetEnv("dbname")
	dbUser := config.GetEnv("dbuser")
	dbPassword := config.GetEnv("dbpassword")
	dbSSLMode := config.GetEnv("dbsslmode")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Failed to connect to the database:%s", err))
	} else {
		// log.Println("Database connection established!")
		logs.AddLog("Info", "roto", "Database connection established!")
	}

}
