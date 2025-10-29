package database

import "context"

func (s *Store) SyncSubscriptionPrices(ctx context.Context) error {
	query := `
	UPDATE subscriptions s
	SET price = sp.price
	FROM subscription_prices sp
	WHERE s.id = sp.subscription_id
	  AND sp.valid_from <= CURRENT_DATE
	  AND sp.valid_to >= CURRENT_DATE
	  AND s.price <> sp.price
	`
	_, err := s.DB.ExecContext(ctx, query)
	return err
}
