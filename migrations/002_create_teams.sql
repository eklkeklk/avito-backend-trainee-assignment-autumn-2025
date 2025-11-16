-- +goose Up
CREATE TABLE teams (
    team_name VARCHAR(255) PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS teams;