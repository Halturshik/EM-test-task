package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (s *Store) UpdateSubscription(ctx context.Context, userID uuid.UUID, serviceName string, newPrice *int, newEndDate *time.Time) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var current Subs
	query := `
		SELECT id, user_id, service_name, price, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1 AND service_name = $2
	`
	err = tx.QueryRowContext(ctx, query, userID, serviceName).Scan(
		&current.ID, &current.UserID, &current.ServiceName, &current.Price, &current.StartDate, &current.EndDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrSubNotFound
	}
	if err != nil {
		return err
	}

	today := time.Now()
	firstNextMonth := time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, today.Location())

	if newPrice != nil && *newPrice != current.Price {
		if *newPrice > current.Price {
			updateSubQuery := `UPDATE subscriptions SET price=$1 WHERE id=$2`
			if _, err := tx.ExecContext(ctx, updateSubQuery, *newPrice, current.ID); err != nil {
				return err
			}

			priceQuery := `
				INSERT INTO subscription_prices (subscription_id, price, valid_from, valid_to)
				VALUES ($1, $2, $3, $4)
			`
			if _, err := tx.ExecContext(ctx, priceQuery, current.ID, *newPrice, today, current.EndDate); err != nil {
				return err
			}

		} else {
			priceQuery := `
				INSERT INTO subscription_prices (subscription_id, price, valid_from, valid_to)
				VALUES ($1, $2, $3, $4)
			`
			if _, err := tx.ExecContext(ctx, priceQuery, current.ID, *newPrice, firstNextMonth, current.EndDate); err != nil {
				return err
			}
		}
	}

	if newEndDate != nil {
		if current.EndDate == nil || newEndDate.After(*current.EndDate) {
			updateEndQuery := `UPDATE subscriptions SET end_date=$1 WHERE id=$2`
			if _, err := tx.ExecContext(ctx, updateEndQuery, newEndDate, current.ID); err != nil {
				return err
			}

			priceHistUpdate := `
				UPDATE subscription_prices
				SET valid_to=$1
				WHERE subscription_id=$2 AND valid_to IS NULL
			`
			if _, err := tx.ExecContext(ctx, priceHistUpdate, newEndDate, current.ID); err != nil {
				return err
			}

		} else if newEndDate.Before(*current.EndDate) && newEndDate.After(today) {
			updateEndQuery := `UPDATE subscriptions SET end_date=$1 WHERE id=$2`
			if _, err := tx.ExecContext(ctx, updateEndQuery, newEndDate, current.ID); err != nil {
				return err
			}

			priceHistUpdate := `
				UPDATE subscription_prices
				SET valid_to=$1
				WHERE subscription_id=$2 AND (valid_to IS NULL OR valid_to > $1)
			`
			if _, err := tx.ExecContext(ctx, priceHistUpdate, newEndDate, current.ID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
