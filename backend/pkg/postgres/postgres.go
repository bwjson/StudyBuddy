package postgres

import (
	"fmt"
	"github.com/bwjson/StudyBuddy/configs"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(cfg configs.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s",
		cfg.PostgresDB.Host, cfg.PostgresDB.User, cfg.PostgresDB.Password, cfg.PostgresDB.Port, cfg.PostgresDB.DBName, cfg.PostgresDB.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
