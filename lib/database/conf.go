// File: conf.go
package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/roto17/zeus/lib/env" // Replace with your actual module path
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB // Singleton GORM database instance
	dbOnce sync.Once
)

// User represents a user in the database (adjust field types as needed)
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

// DBConnection returns a singleton instance of the GORM database connection
func DBConnection() (*gorm.DB, error) {
	var err error

	dbOnce.Do(func() {
		dbPort := env.GetEnv("dbport")
		dbName := env.GetEnv("dbname")
		dbUser := env.GetEnv("dbuser")
		dbPassword := env.GetEnv("dbpassword")
		dbHost := env.GetEnv("dbhost")
		dbsslmode := env.GetEnv("dbsslmode")

		connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbHost, dbUser, dbPassword, dbName, dbPort, dbsslmode)

		// Use GORM with PostgreSQL
		db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil {
			err = fmt.Errorf("unable to connect to database: %v", err)
		}

		// Automatically migrate the schema
		err = db.AutoMigrate(&User{})
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}
