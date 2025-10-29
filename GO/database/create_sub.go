package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func (s *Store) CreateSubscription(ctx context.Context, sub *Subs) error {
	if sub.EndDate == nil {
		t := time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
		sub.EndDate = &t
	}

	var activeStart, activeEnd *time.Time
	checkConflictQuery := `
		SELECT start_date, end_date
        FROM subscriptions
        WHERE user_id=$1 AND service_name=$2
          AND (end_date IS NULL OR end_date >= CURRENT_DATE)
        LIMIT 1
	`

	err := s.DB.QueryRowContext(ctx, checkConflictQuery, sub.UserID, sub.ServiceName).Scan(&activeStart, &activeEnd)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if activeStart != nil {
		if sub.EndDate == nil || !sub.EndDate.Before(*activeStart) {
			return ErrSubIsExist
		}
	}

	query := `
		INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var subID int
	err = s.DB.QueryRowContext(ctx, query, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate).Scan(&subID)
	if err != nil {
		return err
	}

	priceQuery := `
		INSERT INTO subscription_prices (subscription_id, price, valid_from, valid_to)
		VALUES ($1, $2, $3, $4)
	`

	_, err = s.DB.ExecContext(ctx, priceQuery, subID, sub.Price, sub.StartDate, sub.EndDate)
	return err

}
