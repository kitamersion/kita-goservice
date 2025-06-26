package database

import (
	"fmt"

	"github.com/kitamersion/go-goservice/internal/config"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Event{},
	)
}
