-- +goose Up
-- +goose StatementBegin
-- Adding the `email_verify` column to the `users` table
ALTER TABLE users ADD COLUMN email_verify BOOL DEFAULT FALSE;

-- Adding the `number_phone` column to the `addresses` table
ALTER TABLE addresses ADD COLUMN number_phone VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Removing the `email_verify` column from the `users` table
ALTER TABLE users DROP COLUMN email_verify;

-- Removing the `number_phone` column from the `addresses` table
ALTER TABLE addresses DROP COLUMN number_phone;
-- +goose StatementEnd
