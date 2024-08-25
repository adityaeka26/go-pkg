package postgres

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(username, password, host, port, dbname string, sslmode bool) (*Postgres, error) {
	ssl := "disable"
	if sslmode {
		ssl = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, username, password, dbname, port, ssl)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Postgres{
		db: db,
	}, err
}

func (p *Postgres) GetDb() *gorm.DB {
	return p.db
}

func (p *Postgres) Close(ctx context.Context) error {
	dbInstance, err := p.db.WithContext(ctx).DB()
	if err != nil {
		return err
	}
	return dbInstance.Close()
}
