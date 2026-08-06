package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/artni96/GophKeeper/internal/server/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TestDependencies struct {
	cfg *config.Config
	DB  *sqlx.DB
}

func NewTestDependencies(ctx context.Context) (*TestDependencies, error) {
	td := &TestDependencies{}
	if err := td.initConfig(); err != nil {
		return nil, err
	}
	if err := td.initDBConn(ctx); err != nil {
		return nil, err
	}
	if err := td.applyMigrations(); err != nil {
		return nil, err
	}

	return td, nil

}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func (td *TestDependencies) initConfig() error {
	cfg := config.Config{}
	configFile := config.ConfigFile{}

	root, err := findProjectRoot()
	if err != nil {
		return err
	}

	configPath := filepath.Join(root, "tests/test_config.json")

	file, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(file, &configFile)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	dbDSN, err := configFile.AssembleDBDsn()
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	cfg.DBDsn = dbDSN
	td.cfg = &cfg
	return nil
}

func (td *TestDependencies) initDBConn(ctx context.Context) error {
	if td.cfg.DBDsn == "" {
		return fmt.Errorf("database dsn is not provided")
	}
	db, err := sqlx.Open("pgx", td.cfg.DBDsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	localCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err = db.PingContext(localCtx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	td.DB = db
	return nil
}

func (c *TestDependencies) applyMigrations() error {
	driver, err := postgres.WithInstance(c.DB.DB, &postgres.Config{})
	if err != nil {
		log.Println(fmt.Errorf("failed to create database driver: %w", err))
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://../../../migrations", "postgres", driver)
	if err != nil {
		log.Println(fmt.Errorf("failed to initialize test migrator: %w", err))
		return err

	}
	if err = migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println(fmt.Errorf("failed to clean up test database: %w", err))
		return err
	}
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Println(fmt.Errorf("failed to run migrations: %w", err))
		return err
	}
	return nil
}
