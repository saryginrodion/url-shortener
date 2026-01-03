-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE urls (
    id VARCHAR NOT NULL PRIMARY KEY,
    url VARCHAR NOT NULL,
    name VARCHAR,
    author_id UUID REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE urls;
-- +goose StatementEnd
