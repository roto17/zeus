// File: conf.go
package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/roto17/zeus/lib/env" // Replace with your actual module path
)

var (
	conn     *pgx.Conn // The singleton connection instance
	connOnce sync.Once // Ensures the connection is initialized only once
	connErr  error     // Captures any error during the connection process
)

// DBConnection returns a singleton PostgreSQL connection
func DBConnection() (*pgx.Conn, error) {
	connOnce.Do(func() {
		dbPort := env.GetEnv("dbport")
		dbName := env.GetEnv("dbname")
		dbUser := env.GetEnv("dbuser")
		dbPassword := env.GetEnv("dbpassword")
		dbHost := env.GetEnv("dbhost")
		dbsslmode := env.GetEnv("dbsslmode")
		dbDriver := env.GetEnv("dbdriver")

		connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
			dbDriver, dbUser, dbPassword, dbHost, dbPort, dbName, dbsslmode)

		// Create the connection
		conn, connErr = pgx.Connect(context.Background(), connString)
		if connErr != nil {
			connErr = fmt.Errorf("unable to connect to database: %v", connErr)
		}
	})

	return conn, connErr
}
