package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"strings"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath) // Используем sqlite вместо sqlite3
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Storage{db: db}

	// Создаем таблицу, если её нет
	if err := storage.createTables(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return storage, nil
}

// createTables создаёт таблицу `subscriptions`, если её нет
func (s *Storage) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS subscriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL,
		card_number TEXT NOT NULL
	);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// AddSubscription saves subscription to db.
func (s *Storage) AddSubscription(ctx context.Context, email string, cardNumber string) (int64, error) {
	const op = "storage.sqlite.AddSubscription"

	stmt, err := s.db.Prepare("INSERT INTO subscriptions(email, card_number) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, email, cardNumber)
	if err != nil {
		// Проверяем, содержит ли ошибка нарушение уникального ограничения
		if strings.Contains(err.Error(), "constraint failed") {
			return 0, fmt.Errorf("%s: %w", op, errors.New("email already exists"))
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// DeleteSubscription removes a subscription from the DB by email and subscription ID.
func (s *Storage) DeleteSubscription(ctx context.Context, email string) error {
	const op = "storage.sqlite.DeleteSubscription"

	// Prepare the SQL statement to delete the subscription
	stmt, err := s.db.Prepare("DELETE FROM subscriptions WHERE email = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	// Execute the statement with the email and subscription ID
	_, err = stmt.ExecContext(ctx, email)
	if err != nil {
		// Handle errors, for example, if no matching subscription is found
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
