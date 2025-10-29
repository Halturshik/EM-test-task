package database

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

func (s *Store) GetSubscriptions(ctx context.Context, userID uuid.UUID, serviceName string) ([]Subs, error) {
	query := `
        SELECT s.id, s.user_id, s.service_name, s.price, s.start_date, s.end_date
        FROM subscriptions s
        WHERE s.user_id = $1
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
