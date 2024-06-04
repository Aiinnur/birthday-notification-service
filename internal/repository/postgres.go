package repository

import (
	"birthday-notification-service/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) Repository {
	return Repository{client: client}
}

func (r Repository) AddUser(ctx context.Context, user entity.User) error {
	_, err := r.client.Exec(ctx, addUser, user.Email, user.Name, user.Birthday)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("user with email %s already exists", user.Email)
			}
		}
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (r Repository) Subscribe(ctx context.Context, subscription entity.Subscription) error {
	_, err := r.client.Exec(ctx, subscribe, subscription.SubscriberID, subscription.BirthdayUserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return fmt.Errorf("The employee %d already has such a subscriber %d", subscription.BirthdayUserID, subscription.SubscriberID)
			}
		}
		return fmt.Errorf("failed to add user: %w", err)
	}

	return nil
}

func (r Repository) Unsubscribe(ctx context.Context, subscription entity.Subscription) error {
	_, err := r.client.Exec(ctx, unsubscribe, subscription.SubscriberID, subscription.BirthdayUserID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (r Repository) GetSubscribersForTodayBirthdays(ctx context.Context) ([]entity.SubscriberInfo, error) {
	rows, err := r.client.Query(ctx, getSubscribersForTodayBirthdays)
	if err != nil {
		return nil, fmt.Errorf("error getting subscribers")
	}
	defer rows.Close()

	var subscribers []entity.SubscriberInfo
	for rows.Next() {
		var subscriber entity.SubscriberInfo
		if err := rows.Scan(&subscriber.SubscriberEmail, &subscriber.BirthdayName, &subscriber.BirthdayEmail); err != nil {
			return nil, fmt.Errorf("error scan subscribers")
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}
