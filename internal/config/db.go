package config

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

//func (cfg *Config) InitDBConn(ctx context.Context) (*sqlx.DB, error) {
//	if cfg.DBDsn == "" {
//		return nil, errors.New("database dsn is not provided")
//	}
//
//	db, err := sqlx.Open("pgx", cfg.DBDsn)
//	if err != nil {
//		return nil, fmt.Errorf("failed to connect to database: %w", err)
//	}
//
//	localCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
//	defer cancel()
//
//	err = db.PingContext(localCtx)
//	if err != nil {
//		return nil, fmt.Errorf("failed to ping database: %w", err)
//	}
//
//	//if err := runMigrations(db); err != nil {
//	//	logger.Info("failed to run migrations",
//	//		zap.String("error message", err.Error()),
//	//	)
//	//} else {
//	//	logger.Info("Migrations completed successfully")
//	//}
//
//	return db, nil
//}

//func runMigrations(db *sqlx.DB) error {
//	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
//	if err != nil {
//		return fmt.Errorf("failed to initialize postgres driver: %w", err)
//	}
//
//	migrator, err := migrate.NewWithDatabaseInstance(
//		"file://migrations",
//		"postgres",
//		driver,
//	)
//	if err != nil {
//		return fmt.Errorf("failed to initialize migrator: %w", err)
//	}
//
//	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
//		return fmt.Errorf("failed to apply migrations: %w", err)
//	}
//
//	return nil
//}
