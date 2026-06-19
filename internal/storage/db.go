package storage

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//go:embed init.sql
var initSQL string

func CheckAndMigrate(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), initSQL)
	if err != nil {
		return fmt.Errorf("ошибка выполнения init.sql: %w", err)
	}

	var count int
	err = db.QueryRow(context.Background(), "select count(*) from users").Scan(&count)
	if err != nil {
		return fmt.Errorf("ошибка проверки количества пользователей в бд: %v", err)
	}

	if count == 0 {
		if err = InsertAdminUser(context.Background(), db); err != nil {
			return fmt.Errorf("ошибка добавления пользователя: %v", err)
		}
	}

	return nil
}

func InsertAdminUser(ctx context.Context, db *pgxpool.Pool) error {
	hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = db.Exec(
		ctx,
		`INSERT INTO users (username, email, password, role_id, confirmed)
		VALUES ($1, $2, $3, $4, $5)`,
		"admin",
		"admin@admin.ru",
		string(hash),
		1,
		true,
	)

	return err
}
