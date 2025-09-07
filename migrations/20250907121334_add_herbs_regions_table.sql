-- +goose Up
-- +goose StatementBegin
CREATE TABLE herbs_regions (
    herb_id INT NOT NULL REFERENCES herbs(id) ON DELETE CASCADE,
    region_id INT NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    PRIMARY KEY (herb_id, region_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS herbs_regions;
-- +goose StatementEnd
