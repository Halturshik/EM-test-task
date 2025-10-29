package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *Store) CalculateTotalSubscriptionCost(ctx context.Context, userID uuid.UUID, serviceName string, from, to time.Time) (int, string, error) {
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM subscriptions
			WHERE user_id = $1 AND service_name = $2)
	`

	if err := s.DB.QueryRowContext(ctx, checkQuery, userID, serviceName).Scan(&exists); err != nil {
		return 0, "", err
	}

	if !exists {
		return 0, "no_subscription", nil
	}

	query := `
	SELECT sp.price, sp.valid_from, sp.valid_to
	FROM subscriptions s
	JOIN subscription_prices sp ON sp.subscription_id = s.id
	WHERE s.user_id = $1
	  AND s.service_name = $2
	  AND s.start_date <= $3
	  AND s.end_date >= $4
	  AND sp.valid_from <= $3
	  AND (sp.valid_to IS NULL OR sp.valid_to >= $4)
	ORDER BY sp.valid_from
	`
	rows, err := s.DB.QueryContext(ctx, query, userID, serviceName, to, from)
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()

	total := 0
	hasOverlap := false

	for rows.Next() {
		var price int
		var validFrom, validTo time.Time

		if err := rows.Scan(&price, &validFrom, &validTo); err != nil {
			return 0, "", err
		}

		start := maxTime(validFrom, from)
		end := minTime(validTo, to)

		if !end.Before(start) {
			months := monthsBetween(start, end)
			total += price * months
			hasOverlap = true
		}
	}

	if err := rows.Err(); err != nil {
		return 0, "", err
	}

	if !hasOverlap {
		return 0, "no_overlap", nil
	}

	return total, "ok", nil
}

func monthsBetween(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())
	return yearDiff*12 + monthDiff + 1
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
