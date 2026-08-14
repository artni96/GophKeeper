package common

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GetNextRecordNumber creates or updates user record number counter.
func GetNextRecordNumber(tx *gorm.DB, userID uuid.UUID) (uint64, error) {
	updateQuery := "UPDATE user_record_number SET current_number = current_number + 1 where user_id = ? RETURNING current_number;"
	var nextNumber uint64 = 1

	result := tx.Raw(updateQuery, userID).Scan(&nextNumber)
	if result.RowsAffected == 0 {
		insertQuery := "INSERT INTO user_record_number(user_id, current_number) VALUES(?, ?);"
		err := tx.Exec(insertQuery, userID, nextNumber).Error
		if err != nil {
			return 0, fmt.Errorf("failed to initialize user records number counter: %w", result.Error)
		}
	}

	if result.Error != nil {
		return 0, fmt.Errorf("failed to get the next user record number: %w", result.Error)
	}
	return nextNumber, nil
}

// InitDBConnByGORM initializes the database connection through GORM by the *sqxl.DB.
func InitDBConnByGORM(db *sql.DB, logger *zap.Logger) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		logger.Info("Failed to initialize the database connection using GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize the database connection using GORM: %w", err)
	}
	return conn, nil
}
