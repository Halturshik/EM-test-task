package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (s *Store) DeleteSubscription(ctx context.Context, userID uuid.UUID, serviceName string, startDate time.Time) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	syncQuery := ` 
	SELECT id FROM subscriptions 
	WHERE user_id=$1 AND service_name=$2 AND start_date=$3
	`

	deletePricesQuery := `
		DELETE FROM subscription_prices 
		WHERE subscription_id=$1
	`
	deleteSubQuery := `
		DELETE FROM subscriptions 
		WHERE id=$1
	`

	var subID int
	if err := tx.QueryRowContext(ctx, syncQuery, userID, serviceName, startDate).Scan(&subID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSubNotFound
		}
		return err
	}

	if _, err := tx.ExecContext(ctx, deletePricesQuery, subID); err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, deleteSubQuery, subID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrSubNotFound
	}

	return tx.Commit()
}
