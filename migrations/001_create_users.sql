-- +goose Up
CREATE TABLE users (
    user_id VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    team_name VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE
);

-- +goose Down
DROP TABLE IF EXISTS users;