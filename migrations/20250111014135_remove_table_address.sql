-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE addresses;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
