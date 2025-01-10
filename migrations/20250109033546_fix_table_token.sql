-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE tokens ADD COLUMN name VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE tokens DROP COLUMN name;
-- +goose StatementEnd