package database

import (
	"context"
)

func (s *Store) CreateSubscription(ctx context.Context, sub *Subs) error {
	var exists bool
	queryCheck := `
        SELECT EXISTS(
            SELECT 1 FROM subscriptions
            WHERE user_id=$1 AND service_name=$2)
    `

	err := s.DB.QueryRowContext(ctx, queryCheck, sub.UserID, sub.ServiceName).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrSubIsExist
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
