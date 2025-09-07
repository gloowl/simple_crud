-- +goose Up
-- +goose StatementBegin
CREATE TABLE herbs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    latin_name VARCHAR(255),
    description TEXT,
    is_poisonous BOOLEAN DEFAULT FALSE,
    image_path VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS herbs;
-- +goose StatementEnd
