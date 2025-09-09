package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var db *sql.DB

// NewConnection creates a new database connection (используем вашу реализацию)
func NewConnection(cfg Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Тест соединения
	if err = connection.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(25)

	log.Println("Successfully connected to database")
	return connection, nil
}

// Connect устанавливает глобальное соединение для CLI
//func Connect(cfg Config) error {
//	var err error
//	db, err = NewConnection(cfg)
//	return err
//}

func Connect(cfg Config) error {
	var err error
	if db, err = NewConnection(cfg); err != nil {
		return err
	}
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// DevConfig - ваша конфигурация для разработки
func DevConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "pwd4adm",
		DBName:   "simple_crud_db",
		SSLMode:  "disable",
	}
}
