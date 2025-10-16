-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    service_name VARCHAR(64) NOT NULL,
    price INTEGER NOT NULL,
    start_date CHAR(7) NOT NULL,
    end_date CHAR(7)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
