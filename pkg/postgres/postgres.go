package postgres

import (
	"birthday-notification-service/internal/repository"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func NewPostgresClient(url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Проверка подключения
	if err = conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("could not acquire connection from PostgreSQL pool: %w", err)
	}

	if err := createTablesIfNotExist(conn); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func ClosePostgresClient(conn *pgxpool.Pool) {
	conn.Close()
}

func createTablesIfNotExist(conn *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if _, err := conn.Exec(ctx, repository.CreateTableUsers); err != nil {
		fmt.Println("Error")
		return fmt.Errorf("could not create users table: %w", err)
	}

	if _, err := conn.Exec(ctx, repository.CreateTableSubscriptions); err != nil {
		return fmt.Errorf("could not create subscriptions table: %w", err)
	}

	return nil
}
