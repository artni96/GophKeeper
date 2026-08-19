package common

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/model/card"
	"github.com/artni96/GophKeeper/internal/server/model/common"
	"github.com/artni96/GophKeeper/internal/server/model/login"
	"github.com/artni96/GophKeeper/internal/server/model/text"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

// InitGORMDBConn initializes the database connection through GORM by the *sqxl.DB.
func InitGORMDBConn(db *sql.DB, logger *zap.Logger) (*gorm.DB, error) {
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

type CRepositoryI[
CreateT common.CreateEntityI,
UpdateT login.UpdateLogin | card.UpdateCard | text.UpdateText,
GetT login.Login | card.Card | text.Text,
] interface {
	Create(ctx context.Context, entity CreateT, notNullFields []string) (uint64, error)
	GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*GetT, error)
	Update(ctx context.Context, entity UpdateT, number uint64, userID uuid.UUID, fieldsToUpdate []string) error
	Delete(ctx context.Context, number uint64, userID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID) ([]GetT, error)
}

// CRepository implements the common repository to manage data through the database.
type CRepository[
CreateT common.CreateEntityI,
UpdateT login.UpdateLogin | card.UpdateCard | text.UpdateText,
GetT login.Login | card.Card | text.Text,
] struct {
	db        *gorm.DB
	tableName string
}

// NewCRepository initializes and return a new common repository instance.
func NewCRepository[
CreateT common.CreateEntityI,
UpdateT login.UpdateLogin | card.UpdateCard | text.UpdateText,
GetT login.Login | card.Card | text.Text,
](sqlDB *sql.DB, logger *zap.Logger, tableName string) (*CRepository[CreateT, UpdateT, GetT], error) {
	db, err := InitGORMDBConn(sqlDB, logger)
	if err != nil {
		return nil, err
	}
	return &CRepository[CreateT, UpdateT, GetT]{db: db, tableName: tableName}, nil
}

// Create creates a new Entity in the database.
func (r *CRepository[CreateT, UpdateT, GetT]) Create(ctx context.Context, entity CreateT, notNullFields []string) (uint64, error) {
	tx := r.db.WithContext(ctx).Begin()

	entityNumber, err := GetNextRecordNumber(tx, entity.GetUserID())
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to get the next user login number: %w", err)
	}
	entity.SetNumber(entityNumber)

	err = tx.Table(r.tableName).Select(notNullFields).Create(entity).Error
	if err != nil {
		tx.Rollback()
		var pgErr *pgconn.PgError
		fmt.Println(err.Error())
		if errors.As(err, &pgErr) {

			if pgErr.Code == "23505" {
				return 0, constants.ErrEntityAlreadyExists
			} else if pgErr.Code == "23502" {
				return 0, constants.ErrRequiredField
			}
		}
		return 0, fmt.Errorf("failed to create the entity: %w", err)
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create login: %w", err)
	}
	return entityNumber, nil
}

// GetByNumber selects and returns the Entity by its number from the database.
func (r *CRepository[CreateT, UpdateT, GetT]) GetByNumber(ctx context.Context, number uint64, userID uuid.UUID) (*GetT, error) {
	var entity GetT

	db := r.db.WithContext(ctx)
	if err := db.Table(r.tableName).Where("number = ? AND user_id = ?", number, userID).First(&entity).Error; err != nil {
		return nil, fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return &entity, nil
}

// Update updates the Entity in the database by its number (only for authors).
func (r *CRepository[CreateT, UpdateT, GetT]) Update(ctx context.Context, entity UpdateT, number uint64, userID uuid.UUID, fieldsToUpdate []string) error {
	db := r.db.WithContext(ctx)
	result := db.Table(r.tableName).Select(fieldsToUpdate).Where("number = ? AND user_id = ?", number, userID).Updates(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update login with the number %d: %w", number, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return nil
}

// Delete removes the Entity from the database by its number (only for authors).
func (r *CRepository[CreateT, UpdateT, GetT]) Delete(ctx context.Context, number uint64, userID uuid.UUID) error {
	db := r.db.WithContext(ctx)
	result := db.Table(r.tableName).Where("number = ? AND user_id = ?", number, userID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to delete login with the number %d: %w", number, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: number: %d", constants.ErrEntityNotFound, number)
	}
	return nil
}

// GetList returns the list of user's entities from the database (only for authors).
func (r *CRepository[CreateT, UpdateT, GetT]) GetList(ctx context.Context, userID uuid.UUID) ([]GetT, error) {
	var dbEntities []GetT
	result := r.db.WithContext(ctx).Table(r.tableName).Where("user_id = ?", userID).Find(&dbEntities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get list of logins: %w", result.Error)
	}
	return dbEntities, nil
}
