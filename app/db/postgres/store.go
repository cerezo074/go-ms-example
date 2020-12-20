package postgres

import (
	"fmt"
	"user/core/entities"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgress", dataSourceName)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conecting to database: %w", err)
	}

	return &Store{
		UserRepository: &UserStore{DB: db},
	}, nil
}

type Store struct {
	entities.UserRepository
}
