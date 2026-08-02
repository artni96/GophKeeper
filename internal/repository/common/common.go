package common

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func GetNextRecordNumber(tx *gorm.DB, userID uuid.UUID) (int64, error) {
	updateQuery := "UPDATE user_record_number SET current_number = current_number + 1 where user_id = ? RETURNING current_number;"
	var nextNumber int64 = 1

	result := tx.Raw(updateQuery, userID).Scan(&nextNumber)
	if result.RowsAffected == 0 {
		insertQuery := "INSERT INTO user_record_number(user_id, current_number) VALUES(?, ?);"
		err := tx.Exec(insertQuery, userID, nextNumber).Error
		if err != nil {
			return 0, fmt.Errorf("failed to initialize user records number counter: %w", result.Error)
		}
	}

	if result.Error != nil {
		return -1, fmt.Errorf("failed to get the next user record number: %w", result.Error)
	}
	return nextNumber, nil
}

func InitDBConnByGORM(db *sqlx.DB, logger *zap.Logger) (*gorm.DB, error) {
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
