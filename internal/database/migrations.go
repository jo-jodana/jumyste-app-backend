package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"jumyste-app-backend/config"
	"jumyste-app-backend/pkg/logger"
	"log/slog"
	"path/filepath"
	"strings"
)

func RunMigrations() {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Failed to connect to database for migrations",
			slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Error("Failed to create migration driver",
			slog.String("error", err.Error()))
		return
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		logger.Log.Error("Failed to resolve migrations path",
			slog.String("error", err.Error()))
		return
	}

	migrationsPath = strings.ReplaceAll(migrationsPath, "\\", "/")
	migrationsPath = "file://" + migrationsPath

	logger.Log.Info("Using migrations path", slog.String("path", migrationsPath))

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		logger.Log.Error("Failed to create migration instance",
			slog.String("error", err.Error()))
		return
	}

	logger.Log.Info("Starting migrations...")

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Info("No new migrations to apply")
		} else {
			logger.Log.Error("Migration failed", slog.String("error", err.Error()))
		}
		return
	}

	logger.Log.Info("Migrations applied successfully")
}
