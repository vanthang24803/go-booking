-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users RENAME COLUMN fist_name TO first_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users RENAME COLUMN first_name TO fist_name;
-- +goose StatementEnd
