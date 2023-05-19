-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE galleries (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  user_id INT REFERENCES users (id) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE galleries;
-- +goose StatementEnd
