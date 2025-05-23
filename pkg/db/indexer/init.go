package indexer

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IndexerRepository struct {
	*gorm.DB
}

var _ db.IndexerAdapter = (*IndexerRepository)(nil)

type ConnConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func NewRepository(ctx context.Context, cfg *ConnConfig) *IndexerRepository {
	dbSource := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dbSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil
	}

	return &IndexerRepository{
		DB: db,
	}
}

func (r *IndexerRepository) Close() {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}
