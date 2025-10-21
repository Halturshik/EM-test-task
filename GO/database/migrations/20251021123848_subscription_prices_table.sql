-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE subscription_prices (
    id SERIAL PRIMARY KEY,
    subscription_id INT NOT NULL REFERENCES subscriptions(id),
    price INT NOT NULL,
    valid_from DATE NOT NULL,
    valid_to DATE NULL
);

CREATE INDEX idx_subscription_prices_dates
ON subscription_prices(subscription_id, valid_from, valid_to);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS subscription_prices;
