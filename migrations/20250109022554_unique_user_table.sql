-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE users DROP CONSTRAINT users_username_unique;
-- +goose StatementEnd
