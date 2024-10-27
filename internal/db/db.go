package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/config"
)

var (
	db   *sql.DB
	once sync.Once
	cfg  *config.Database
)

// Initialize the DB connection
func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to check the connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Connected to MySQL database successfully")
}

// Initialize the database configuration only once
func InitDBCredentials(dbConfig *config.Database) {
	cfg = dbConfig
}

// GetDB returns the singleton instance of the database
func GetDB() *sql.DB {
	once.Do(initDB) // Ensures initDB is only called once
	return db
}
