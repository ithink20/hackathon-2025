package database

import (
	"log"
)

func Init() error {
	config := NewConfig()

	log.Println("Initializing database connection...")

	if err := Connect(config); err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func InitWithConfig(config *Config) error {
	log.Println("Initializing database connection with custom config...")

	if err := Connect(config); err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
