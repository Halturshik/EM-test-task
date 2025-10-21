package database

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

func (s *Store) GetSubscriptions(ctx context.Context, userID uuid.UUID, serviceName string) ([]Subs, error) {
	query := `
        SELECT s.id, s.user_id, s.service_name, sp.price, s.start_date, s.end_date
        FROM subscriptions s
        JOIN subscription_prices sp
          ON s.id = sp.subscription_id
        WHERE s.user_id = $1
          AND sp.valid_from <= CURRENT_DATE
          AND (sp.valid_to IS NULL OR sp.valid_to >= CURRENT_DATE)
    `
	args := []any{userID}

	if strings.TrimSpace(serviceName) != "" {
		query += " AND s.service_name = $2"
		args = append(args, serviceName)
	}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Subs
	for rows.Next() {
		var s Subs
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.ServiceName, &s.Price, &s.StartDate, &s.EndDate,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
