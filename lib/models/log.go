package models

import (
	// or any other driver (e.g., sqlite, mysql)
	"time"
)

type Log struct {
	ID        uint      `gorm:"primaryKey"`         // Automatically adds an ID as the primary key
	Username  string    `gorm:"size:100;not null"`  // Username (up to 100 characters)
	LogType   string    `gorm:"size:50;not null"`   // Log type (up to 50 characters)
	Message   string    `gorm:"type:text;not null"` // Log message
	CreatedAt time.Time `gorm:"autoCreateTime"`     // Automatically stores the creation timestamp

}
