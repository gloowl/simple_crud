-- +goose Up
-- +goose StatementBegin
CREATE TABLE usage_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS usage_types;
-- +goose StatementEnd
