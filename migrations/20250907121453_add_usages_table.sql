-- +goose Up
-- +goose StatementBegin
CREATE TABLE usages (
    id SERIAL PRIMARY KEY,
    herb_id INT NOT NULL REFERENCES herbs(id) ON DELETE CASCADE,
    usage_type_id INT NOT NULL REFERENCES usage_types(id) ON DELETE RESTRICT,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS usages;
-- +goose StatementEnd
