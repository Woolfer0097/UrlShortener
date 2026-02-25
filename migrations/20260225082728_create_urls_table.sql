-- +goose Up
CREATE TABLE urls (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    url_code VARCHAR(10) NOT NULL UNIQUE,
    original_url VARCHAR(2048) NOT NULL
);

CREATE UNIQUE INDEX idx_url_code ON urls(url_code);
CREATE INDEX idx_original_url ON urls(original_url);

-- +goose Down
DROP TABLE urls;