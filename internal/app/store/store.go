package store

import (
	"database/sql"
	"hack2023/internal/app/config"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db *sql.DB
}

func New(config config.Config) (*Store, error) {
	db, err := sql.Open("mysql", config.URL)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) NewNullString(str string) sql.NullString {
	if len(str) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: str,
		Valid:  true,
	}
}
