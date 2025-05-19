package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"roster-management/pkg/postgres"
)

const dateFormat = "2006-01-02 15:04:05"

type Repository struct {
	db *sqlx.DB
}

func NewRepositoryFromConfig(config *viper.Viper) (*Repository, error) {
	db, err := postgres.Connectx(config)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}
